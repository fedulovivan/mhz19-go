package rules

import (
	"context"
	"database/sql"
	"slices"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"golang.org/x/sync/errgroup"
)

type RulesRepository interface {
	Get() (rules []DbRule, conditions []DbRuleCondition, ruleActions []DbRuleAction, args []DbRuleConditionOrActionArgument, mappings []DbRuleActionArgumentMapping, err error)
	Create(rule DbRule, conditions []DbRuleCondition, arguments []DbRuleConditionOrActionArgument) error
}

type rulesRepository struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) RulesRepository {
	return rulesRepository{
		database: database,
	}
}

type DbRule struct {
	Id         int32
	Comments   string
	IsDisabled sql.NullInt32
	Throttle   sql.NullInt32
}

type DbRuleCondition struct {
	Id                int32
	RuleId            int32
	FunctionType      sql.NullInt32
	LogicOr           sql.NullInt32
	ParentConditionId sql.NullInt32
}

type DbRuleAction struct {
	Id           int32
	RuleId       int32
	FunctionType sql.NullInt32
	DeviceId     sql.NullString
}

type DbRuleConditionOrActionArgument struct {
	Id            int32
	ConditionId   sql.NullInt32
	ActionId      sql.NullInt32
	ArgumentName  string
	IsList        sql.NullInt32
	Value         sql.NullString
	DeviceId      sql.NullString
	DeviceClassId sql.NullInt32
}

type DbRuleActionArgumentMapping struct {
	Id         int32
	ArgumentId int32
	Key        string
	Value      string
}

func conditionsSelect(ctx context.Context, tx *sql.Tx) ([]DbRuleCondition, error) {
	return db.Select(
		tx,
		ctx,
		`SELECT
			id,
			rule_id,
			function_type,
			logic_or,
			parent_condition_id
		FROM
			rule_conditions`,
		func(rows *sql.Rows, m *DbRuleCondition) error {
			return rows.Scan(&m.Id, &m.RuleId, &m.FunctionType, &m.LogicOr, &m.ParentConditionId)
		},
	)
}

func conditionInsert(
	cond DbRuleCondition,
	ctx context.Context,
	tx *sql.Tx,
) (sql.Result, error) {
	return db.Insert(
		tx,
		ctx,
		`INSERT INTO rule_conditions(rule_id, function_type, logic_or, parent_condition_id) VALUES(?,?,?,?)`,
		cond.RuleId, cond.FunctionType, cond.LogicOr, cond.ParentConditionId,
	)
}

func argumentInsert(
	arg DbRuleConditionOrActionArgument,
	ctx context.Context,
	tx *sql.Tx,
) (sql.Result, error) {
	return db.Insert(
		tx,
		ctx,
		`INSERT INTO rule_condition_or_action_arguments(condition_id, action_id, argument_name, is_list, value, device_id, device_class_id) VALUES(?,?,?,?,?,?,?)`,
		arg.ConditionId, arg.ActionId, arg.ArgumentName, arg.IsList, arg.Value, arg.DeviceId, arg.DeviceClassId,
	)
}

func ruleInsert(
	rule DbRule,
	ctx context.Context,
	tx *sql.Tx,
) (sql.Result, error) {
	return db.Insert(
		tx,
		ctx,
		`INSERT INTO rules(comments, is_disabled, throttle) VALUES(?,?,?)`,
		rule.Comments, rule.IsDisabled, rule.Throttle,
	)
}

func rulesSelect(ctx context.Context, tx *sql.Tx) ([]DbRule, error) {
	return db.Select(
		tx,
		ctx,
		`SELECT
			id,
			comments,
			is_disabled,
			throttle
		FROM 
			rules`,
		func(rows *sql.Rows, m *DbRule) error {
			return rows.Scan(&m.Id, &m.Comments, &m.IsDisabled, &m.Throttle)
		},
	)
}

func ruleActionsSelect(ctx context.Context, tx *sql.Tx) ([]DbRuleAction, error) {
	return db.Select(
		tx,
		ctx,
		`SELECT
			id,
			rule_id,
			function_type,
			device_id
		FROM
			rule_actions`,
		func(rows *sql.Rows, m *DbRuleAction) error {
			return rows.Scan(&m.Id, &m.RuleId, &m.FunctionType, &m.DeviceId)
		},
	)
}

func argsSelect(ctx context.Context, tx *sql.Tx) ([]DbRuleConditionOrActionArgument, error) {
	return db.Select(
		tx,
		ctx,
		`SELECT
			id,
			condition_id,
			action_id,
			argument_name,
			is_list,
			value,
			device_id,
			device_class_id
		FROM
			rule_condition_or_action_arguments`,
		func(rows *sql.Rows, m *DbRuleConditionOrActionArgument) error {
			return rows.Scan(
				&m.Id, &m.ConditionId, &m.ActionId, &m.ArgumentName, &m.IsList, &m.Value, &m.DeviceId, &m.DeviceClassId,
			)
		},
	)
}

func mappingsSelect(ctx context.Context, tx *sql.Tx) ([]DbRuleActionArgumentMapping, error) {
	return db.Select(
		tx,
		ctx,
		`SELECT id, argument_id, key, value FROM rule_action_argument_mappings`,
		func(rows *sql.Rows, m *DbRuleActionArgumentMapping) error {
			return rows.Scan(
				&m.Id, &m.ArgumentId, &m.Key, &m.Value,
			)
		},
	)
}

func (repo rulesRepository) Create(
	rule DbRule,
	conditions []DbRuleCondition,
	arguments []DbRuleConditionOrActionArgument,
) (err error) {
	defer utils.TimeTrack(logTag, time.Now(), "RuleCreate")
	ctx := context.Background()
	tx, err := repo.database.Begin()
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	// rule
	result, err := ruleInsert(rule, ctx, tx)
	if err != nil {
		return
	}
	ruleId, err := result.LastInsertId()
	if err != nil {
		return
	}
	slices.SortFunc(conditions, func(a, b DbRuleCondition) int {
		return int(a.ParentConditionId.Int32 - b.ParentConditionId.Int32)
	})
	var realCondIdsMap = make(map[int32]int32, len(conditions))
	// conditions
	for _, cond := range conditions {
		cond.RuleId = int32(ruleId)
		if cond.ParentConditionId.Valid {
			realParentId := realCondIdsMap[cond.ParentConditionId.Int32]
			cond.ParentConditionId = db.NewNullInt32(realParentId)
		}
		result, err = conditionInsert(cond, ctx, tx)
		if err != nil {
			return
		}
		var condId int64
		condId, err = result.LastInsertId()
		if err != nil {
			return
		}
		realCondIdsMap[cond.Id] = int32(condId)
	}
	// arguments
	for _, arg := range arguments {
		realCondId := realCondIdsMap[arg.ConditionId.Int32]
		arg.ConditionId = db.NewNullInt32(realCondId)
		_, err = argumentInsert(arg, ctx, tx)
		if err != nil {
			return
		}
	}
	if err == nil {
		err = tx.Commit()
	}
	return
}

func (repo rulesRepository) Get() (
	rules []DbRule,
	conditions []DbRuleCondition,
	ruleActions []DbRuleAction,
	args []DbRuleConditionOrActionArgument,
	mappings []DbRuleActionArgumentMapping,
	err error,
) {
	defer utils.TimeTrack(logTag, time.Now(), "RuleFetchAll")
	g, ctx := errgroup.WithContext(context.Background())
	tx, err := repo.database.Begin()
	if err != nil {
		return
	}
	g.Go(func() (e error) {
		rules, e = rulesSelect(ctx, tx)
		return
	})
	g.Go(func() (e error) {
		conditions, e = conditionsSelect(ctx, tx)
		return
	})
	g.Go(func() (e error) {
		ruleActions, e = ruleActionsSelect(ctx, tx)
		return
	})
	g.Go(func() (e error) {
		args, e = argsSelect(ctx, tx)
		return
	})
	g.Go(func() (e error) {
		mappings, e = mappingsSelect(ctx, tx)
		return
	})
	err = g.Wait()
	if err == nil {
		err = db.Commit(tx)
	}
	return
}
