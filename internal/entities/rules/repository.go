package rules

import (
	"context"
	"database/sql"
	"slices"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"golang.org/x/sync/errgroup"
)

type RulesRepository interface {
	Get(ruleId sql.NullInt32) (rules []DbRule, conditions []DbRuleCondition, ruleActions []DbRuleAction, args []DbRuleConditionOrActionArgument, mappings []DbRuleActionArgumentMapping, err error)
	Create(rule DbRule, conditions []DbRuleCondition, actions []DbRuleAction, arguments []DbRuleConditionOrActionArgument, mappings []DbRuleActionArgumentMapping) (int64, error)
}

var _ RulesRepository = (*rulesRepository)(nil)

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
	Name       string
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
	// DeviceId     sql.NullString
}

type DbRuleConditionOrActionArgument struct {
	Id            int32
	RuleId        int32
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
	RuleId     int32
	ArgumentId int32
	Key        string
	Value      string
}

func actionInsertTx(
	act DbRuleAction,
	ctx context.Context,
	tx *sql.Tx,
) (sql.Result, error) {
	return db.Insert(
		tx,
		ctx,
		`INSERT INTO rule_actions(rule_id, function_type) VALUES(?,?)`,
		act.RuleId, act.FunctionType,
	)
}

func conditionInsertTx(
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

func mappingInsertTx(
	mapping DbRuleActionArgumentMapping,
	ctx context.Context,
	tx *sql.Tx,
) (sql.Result, error) {
	return db.Insert(
		tx,
		ctx,
		`INSERT INTO rule_action_argument_mappings(rule_id, argument_id, key, value) VALUES(?,?,?,?)`,
		mapping.RuleId, mapping.ArgumentId, mapping.Key, mapping.Value,
	)
}

func argumentInsertTx(
	arg DbRuleConditionOrActionArgument,
	ctx context.Context,
	tx *sql.Tx,
) (sql.Result, error) {
	return db.Insert(
		tx,
		ctx,
		`INSERT INTO rule_condition_or_action_arguments(rule_id, condition_id, action_id, argument_name, is_list, value, device_id, device_class_id) VALUES(?,?,?,?,?,?,?,?)`,
		arg.RuleId, arg.ConditionId, arg.ActionId, arg.ArgumentName, arg.IsList, arg.Value, arg.DeviceId, arg.DeviceClassId,
	)
}

func ruleInsertTx(
	rule DbRule,
	ctx context.Context,
	tx *sql.Tx,
) (sql.Result, error) {
	return db.Insert(
		tx,
		ctx,
		`INSERT INTO rules(name, is_disabled, throttle) VALUES(?,?,?)`,
		rule.Name, rule.IsDisabled, rule.Throttle,
	)
}

func rulesSelectTx(ctx context.Context, tx *sql.Tx, ruleId sql.NullInt32) ([]DbRule, error) {
	return db.Select(
		tx,
		ctx,
		`SELECT
			id,
			name,
			is_disabled,
			throttle
		FROM 
			rules`,
		func(rows *sql.Rows, m *DbRule) error {
			return rows.Scan(&m.Id, &m.Name, &m.IsDisabled, &m.Throttle)
		},
		db.Where{
			"id": ruleId,
		},
	)
}

func conditionsSelectTx(ctx context.Context, tx *sql.Tx, ruleId sql.NullInt32) ([]DbRuleCondition, error) {
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
		db.Where{
			"rule_id": ruleId,
		},
	)
}

func ruleActionsSelectTx(ctx context.Context, tx *sql.Tx, ruleId sql.NullInt32) ([]DbRuleAction, error) {
	return db.Select(
		tx,
		ctx,
		`SELECT
			id,
			rule_id,
			function_type
		FROM
			rule_actions`,
		func(rows *sql.Rows, m *DbRuleAction) error {
			return rows.Scan(&m.Id, &m.RuleId, &m.FunctionType)
		},
		db.Where{
			"rule_id": ruleId,
		},
	)
}

func argsSelectTx(ctx context.Context, tx *sql.Tx, ruleId sql.NullInt32) ([]DbRuleConditionOrActionArgument, error) {
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
		db.Where{
			"rule_id": ruleId,
		},
	)
}

func mappingsSelectTx(ctx context.Context, tx *sql.Tx, ruleId sql.NullInt32) ([]DbRuleActionArgumentMapping, error) {
	return db.Select(
		tx,
		ctx,
		`SELECT id, argument_id, key, value FROM rule_action_argument_mappings`,
		func(rows *sql.Rows, m *DbRuleActionArgumentMapping) error {
			return rows.Scan(
				&m.Id, &m.ArgumentId, &m.Key, &m.Value,
			)
		},
		db.Where{
			"rule_id": ruleId,
		},
	)
}

func CountTx(ctx context.Context, tx *sql.Tx) (int32, error) {
	return db.Count(
		tx,
		ctx,
		`SELECT COUNT(*) FROM rules`,
	)
}

func (r rulesRepository) Create(
	rule DbRule,
	conditions []DbRuleCondition,
	actions []DbRuleAction,
	arguments []DbRuleConditionOrActionArgument,
	mappings []DbRuleActionArgumentMapping,
) (ruleId int64, err error) {
	var realCondIdsMap = make(map[int32]int32, len(conditions))
	var realActionIdsMap = make(map[int32]int32, len(actions))
	var realArgIdsMap = make(map[int32]int32, len(arguments))
	ctx := context.Background()
	tx, err := r.database.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return
	}
	// rule
	result, err := ruleInsertTx(rule, ctx, tx)
	if err != nil {
		return
	}
	ruleId, err = result.LastInsertId()
	if err != nil {
		return
	}
	slices.SortFunc(conditions, func(a, b DbRuleCondition) int {
		return int(a.ParentConditionId.Int32 - b.ParentConditionId.Int32)
	})
	// conditions
	for _, cond := range conditions {
		cond.RuleId = int32(ruleId)
		if cond.ParentConditionId.Valid {
			realParentId := realCondIdsMap[cond.ParentConditionId.Int32]
			cond.ParentConditionId = db.NewNullInt32(realParentId)
		}
		result, err = conditionInsertTx(cond, ctx, tx)
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
	// actions
	for _, act := range actions {
		act.RuleId = int32(ruleId)
		result, err = actionInsertTx(act, ctx, tx)
		if err != nil {
			return
		}
		var actId int64
		actId, err = result.LastInsertId()
		if err != nil {
			return
		}
		realActionIdsMap[act.Id] = int32(actId)
	}
	// arguments
	for _, arg := range arguments {
		arg.RuleId = int32(ruleId)
		if arg.ConditionId.Valid {
			realCondId := realCondIdsMap[arg.ConditionId.Int32]
			arg.ConditionId = db.NewNullInt32(realCondId)
		}
		if arg.ActionId.Valid {
			realActId := realActionIdsMap[arg.ActionId.Int32]
			arg.ActionId = db.NewNullInt32(realActId)
		}
		result, err = argumentInsertTx(arg, ctx, tx)
		if err != nil {
			return
		}
		var argId int64
		argId, err = result.LastInsertId()
		if err != nil {
			return
		}
		realArgIdsMap[arg.Id] = int32(argId)
	}
	// mappings
	for _, mapping := range mappings {
		mapping.RuleId = int32(ruleId)
		mapping.ArgumentId = realArgIdsMap[mapping.ArgumentId]
		_, err = mappingInsertTx(mapping, ctx, tx)
		if err != nil {
			return
		}
	}
	if err == nil {
		err = tx.Commit()
	}
	return
}

func (r rulesRepository) Get(ruleId sql.NullInt32) (
	rules []DbRule,
	conditions []DbRuleCondition,
	ruleActions []DbRuleAction,
	args []DbRuleConditionOrActionArgument,
	mappings []DbRuleActionArgumentMapping,
	err error,
) {
	g, ctx := errgroup.WithContext(context.Background())
	tx, err := r.database.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return
	}
	g.Go(func() (e error) { rules, e = rulesSelectTx(ctx, tx, ruleId); return })
	g.Go(func() (e error) { conditions, e = conditionsSelectTx(ctx, tx, ruleId); return })
	g.Go(func() (e error) { ruleActions, e = ruleActionsSelectTx(ctx, tx, ruleId); return })
	g.Go(func() (e error) { args, e = argsSelectTx(ctx, tx, ruleId); return })
	g.Go(func() (e error) { mappings, e = mappingsSelectTx(ctx, tx, ruleId); return })
	err = g.Wait()
	if err == nil {
		err = db.Commit(tx)
	}
	return
}
