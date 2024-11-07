package rules

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"golang.org/x/sync/errgroup"
)

type RulesRepository interface {
	Delete(ruleId int32) error
	Get(ruleId sql.NullInt32) (rules []DbRule, conditions []DbRuleCondition, ruleActions []DbRuleAction, args []DbRuleConditionOrActionArgument, mappings []DbRuleActionArgumentMapping, err error)
	Create(rule DbRule, conditions []DbRuleCondition, actions []DbRuleAction, arguments []DbRuleConditionOrActionArgument, mappings []DbRuleActionArgumentMapping) (int64, error)
}

var _ RulesRepository = (*repo)(nil)

type repo struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) repo {
	return repo{
		database: database,
	}
}

type DbRule struct {
	Id         int32
	Name       string
	IsDisabled sql.NullInt32
	ThrottleMs sql.NullInt32
}

type DbRuleCondition struct {
	Id                int32
	RuleId            int32
	FunctionType      sql.NullInt32
	LogicOr           sql.NullInt32
	Not               sql.NullInt32
	ParentConditionId sql.NullInt32
	OtherDeviceId     sql.NullString
	// FunctionType      *int32
	// LogicOr           *int32
	// Not               *int32
	// ParentConditionId *int32
	// OtherDeviceId     *string
}

type DbRuleAction struct {
	Id           int32
	RuleId       int32
	FunctionType sql.NullInt32
}

type DbRuleConditionOrActionArgument struct {
	Id            int32
	RuleId        int32
	ConditionId   sql.NullInt32
	ActionId      sql.NullInt32
	ArgumentName  string
	IsList        sql.NullInt32
	Value         sql.NullString
	ValueDataType sql.NullString
	DeviceId      sql.NullString
	DeviceClassId sql.NullInt32
	ChannelTypeId sql.NullInt32
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
) (sql.Result, error) {
	return db.Exec(
		ctx,
		`INSERT INTO rule_actions(rule_id, function_type) VALUES(?,?)`,
		act.RuleId, act.FunctionType,
	)
}

func conditionInsertTx(
	cond DbRuleCondition,
	ctx context.Context,
) (sql.Result, error) {
	return db.Exec(
		ctx,
		`INSERT INTO rule_conditions(rule_id, function_type, logic_or, parent_condition_id, other_device_id, function_inverted) VALUES(?,?,?,?,?,?)`,
		cond.RuleId, cond.FunctionType, cond.LogicOr, cond.ParentConditionId, cond.OtherDeviceId, cond.Not,
	)
}

func mappingInsertTx(
	mapping DbRuleActionArgumentMapping,
	ctx context.Context,
) (sql.Result, error) {
	return db.Exec(
		ctx,
		`INSERT INTO rule_action_argument_mappings(rule_id, argument_id, key, value) VALUES(?,?,?,?)`,
		mapping.RuleId, mapping.ArgumentId, mapping.Key, mapping.Value,
	)
}

func argumentInsertTx(
	arg DbRuleConditionOrActionArgument,
	ctx context.Context,
) (sql.Result, error) {
	return db.Exec(
		ctx,
		`INSERT INTO rule_condition_or_action_arguments(
			rule_id, 
			condition_id, 
			action_id, 
			argument_name, 
			is_list, 
			value, 
			value_data_type,
			device_id, 
			device_class_id,
			channel_type_id
		) VALUES(?,?,?,?,?,?,?,?,?,?)`,
		arg.RuleId,
		arg.ConditionId,
		arg.ActionId,
		arg.ArgumentName,
		arg.IsList,
		arg.Value,
		arg.ValueDataType,
		arg.DeviceId,
		arg.DeviceClassId,
		arg.ChannelTypeId,
	)
}

func ruleInsertTx(
	rule DbRule,
	ctx context.Context,
) (sql.Result, error) {
	return db.Exec(
		ctx,
		`INSERT INTO rules(name, is_disabled, throttle_ms) VALUES(?,?,?)`,
		rule.Name, rule.IsDisabled, rule.ThrottleMs,
	)
}

func rulesSelectTx(ctx context.Context, ruleId sql.NullInt32) ([]DbRule, error) {
	return db.Select(
		ctx,
		`SELECT
			id,
			name,
			is_disabled,
			throttle_ms
		FROM 
			rules`,
		func(rows *sql.Rows, m *DbRule) error {
			return rows.Scan(&m.Id, &m.Name, &m.IsDisabled, &m.ThrottleMs)
		},
		db.Where{
			"id": ruleId,
		},
	)
}

func conditionsSelectTx(ctx context.Context, ruleId sql.NullInt32) ([]DbRuleCondition, error) {
	return db.Select(
		ctx,
		`SELECT
			id,
			rule_id,
			function_type,
			logic_or,
			parent_condition_id,
			other_device_id,
			function_inverted
		FROM
			rule_conditions`,
		func(rows *sql.Rows, m *DbRuleCondition) error {
			return rows.Scan(&m.Id, &m.RuleId, &m.FunctionType, &m.LogicOr, &m.ParentConditionId, &m.OtherDeviceId, &m.Not)
		},
		db.Where{
			"rule_id": ruleId,
		},
	)
}

func ruleActionsSelectTx(ctx context.Context, ruleId sql.NullInt32) ([]DbRuleAction, error) {
	return db.Select(
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

func argsSelectTx(ctx context.Context, ruleId sql.NullInt32) ([]DbRuleConditionOrActionArgument, error) {
	return db.Select(
		ctx,
		`SELECT
			id,
			condition_id,
			action_id,
			argument_name,
			is_list,
			value,
			value_data_type,
			device_id,
			device_class_id,
			channel_type_id
		FROM
			rule_condition_or_action_arguments`,
		func(rows *sql.Rows, m *DbRuleConditionOrActionArgument) error {
			return rows.Scan(
				&m.Id,
				&m.ConditionId,
				&m.ActionId,
				&m.ArgumentName,
				&m.IsList,
				&m.Value,
				&m.ValueDataType,
				&m.DeviceId,
				&m.DeviceClassId,
				&m.ChannelTypeId,
			)
		},
		db.Where{
			"rule_id": ruleId,
		},
	)
}

func mappingsSelectTx(ctx context.Context, ruleId sql.NullInt32) ([]DbRuleActionArgumentMapping, error) {
	return db.Select(
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

func CountTx(ctx context.Context) (int32, error) {
	return db.Count(
		ctx,
		`SELECT COUNT(*) FROM rules`,
	)
}

func CountActionsTx(ctx db.CtxEnhanced) (int32, error) {
	return db.Count(
		ctx,
		`SELECT COUNT(*) FROM rule_actions`,
	)
}

func CountCondsTx(ctx db.CtxEnhanced) (int32, error) {
	return db.Count(
		ctx,
		`SELECT COUNT(*) FROM rule_conditions`,
	)
}

func CountArgsTx(ctx db.CtxEnhanced) (int32, error) {
	return db.Count(
		ctx,
		`SELECT COUNT(*) FROM rule_condition_or_action_arguments`,
	)
}

func CountMappingsTx(ctx db.CtxEnhanced) (int32, error) {
	return db.Count(
		ctx,
		`SELECT COUNT(*) FROM rule_action_argument_mappings`,
	)
}

func (r repo) Create(
	rule DbRule,
	conditions []DbRuleCondition,
	actions []DbRuleAction,
	arguments []DbRuleConditionOrActionArgument,
	mappings []DbRuleActionArgumentMapping,
) (ruleId int64, err error) {
	err = db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {

		var realCondIdsMap = make(map[int32]int32, len(conditions))
		var realActionIdsMap = make(map[int32]int32, len(actions))
		var realArgIdsMap = make(map[int32]int32, len(arguments))

		// rule
		result, err := ruleInsertTx(rule, ctx)
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
			result, err = conditionInsertTx(cond, ctx)
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
			result, err = actionInsertTx(act, ctx)
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
			result, err = argumentInsertTx(arg, ctx)
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
			_, err = mappingInsertTx(mapping, ctx)
			if err != nil {
				return
			}
		}

		return

	})

	// if err == nil {
	// 	err = tx.Commit()
	// }

	return
}

func (r repo) Get(ruleId sql.NullInt32) (
	rules []DbRule,
	conditions []DbRuleCondition,
	ruleActions []DbRuleAction,
	args []DbRuleConditionOrActionArgument,
	mappings []DbRuleActionArgumentMapping,
	err error,
) {
	err = db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() (e error) { rules, e = rulesSelectTx(ctx, ruleId); return })
		g.Go(func() (e error) { conditions, e = conditionsSelectTx(ctx, ruleId); return })
		g.Go(func() (e error) { ruleActions, e = ruleActionsSelectTx(ctx, ruleId); return })
		g.Go(func() (e error) { args, e = argsSelectTx(ctx, ruleId); return })
		g.Go(func() (e error) { mappings, e = mappingsSelectTx(ctx, ruleId); return })
		return g.Wait()
	})
	return

}

func ruleDeleteTx(
	ruleId int32,
	ctx context.Context,
) (sql.Result, error) {
	return db.Exec(
		ctx,
		`DELETE FROM rules WHERE id = ?`,
		ruleId,
	)
}

func (r repo) Delete(ruleId int32) error {
	return db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
		res, err := ruleDeleteTx(ruleId, ctx)
		if err != nil {
			return
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return
		}
		if rowsAffected == 0 {
			err = fmt.Errorf("no one rule was deleted")
		}
		return
	})
}

// baseCtx := context.WithValue(
// 	context.Background(),
// 	db.Ctxkey_tag{},
// 	db.BaseTag.WithTid("Tx"),
// )
// g, ctx := errgroup.WithContext(baseCtx)
// tx, err := r.database.Begin()
// defer db.Rollback(tx)
// if err != nil {
// 	return
// }
// g.Go(func() (e error) { rules, e = rulesSelectTx(ctx, tx, ruleId); return })
// g.Go(func() (e error) { conditions, e = conditionsSelectTx(ctx, tx, ruleId); return })
// g.Go(func() (e error) { ruleActions, e = ruleActionsSelectTx(ctx, tx, ruleId); return })
// g.Go(func() (e error) { args, e = argsSelectTx(ctx, tx, ruleId); return })
// g.Go(func() (e error) { mappings, e = mappingsSelectTx(ctx, tx, ruleId); return })
// err = g.Wait()
// if err == nil {
// 	err = db.Commit(tx)
// }

// ctx := context.Background()
// tx, err := r.database.Begin()
// ctx = context.WithValue(ctx, db.Ctxkey_tx{}, tx)
// ctx = context.WithValue(ctx, db.Ctxkey_tag{}, db.BaseTag.WithTid("Tx"))
// defer db.Rollback(ctx)
// if err != nil {
// 	return
// }
