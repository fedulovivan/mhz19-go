package rules

import (
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/samber/lo"
)

type RulesService interface {
	Get() ([]engine.Rule, error)
	Create(rule engine.Rule) error
}

type rulesService struct {
	repository RulesRepository
}

func NewService(r RulesRepository) RulesService {
	return rulesService{
		repository: r,
	}
}

func (s rulesService) Create(rule engine.Rule) error {
	dbRule, dbConditions, dbArguments, err := ToDb(rule, utils.NewSeq())
	if err != nil {
		return err
	}
	return s.repository.Create(
		dbRule,
		dbConditions,
		dbArguments,
	)
}

func (s rulesService) Get() ([]engine.Rule, error) {
	rules,
		conditions,
		ruleActions,
		ruleConditionOrActionArguments,
		ruleActionArgumentMappings,
		err := s.repository.Get()
	if err != nil {
		return nil, err
	}
	return Build(
		rules,
		conditions,
		ruleActions,
		ruleConditionOrActionArguments,
		ruleActionArgumentMappings,
	), nil
}

func Build(
	allRules []DbRule,
	allConditions []DbRuleCondition,
	allRuleActions []DbRuleAction,
	allArgs []DbRuleConditionOrActionArgument,
	allMappings []DbRuleActionArgumentMapping,
) (result []engine.Rule) {
	for _, r := range allRules {
		rootCond, rootCondFound := lo.Find(allConditions, func(c DbRuleCondition) bool {
			return c.RuleId == r.Id && !c.ParentConditionId.Valid
		})
		cond := engine.Condition{}
		if rootCondFound {
			cond = BuildCondition(rootCond.Id, allConditions, allArgs)
		}
		rule := engine.Rule{
			Id:        r.Id,
			Comments:  r.Comments,
			Disabled:  r.IsDisabled.Int32 == 1,
			Condition: cond,
			Actions:   BuildActions(r.Id, allRuleActions, allArgs),
		}
		result = append(result, rule)
	}
	return
}

func BuildCondition(
	rootConditionId int32,
	conditions []DbRuleCondition,
	allArgs []DbRuleConditionOrActionArgument,
) (cond engine.Condition) {
	if len(conditions) == 0 {
		return
	}
	root, rootFound := lo.Find(conditions, func(c DbRuleCondition) bool {
		return c.Id == rootConditionId
	})
	if !rootFound {
		return
	}
	isFn := root.FunctionType.Valid
	if isFn {
		// build function node
		args := lo.Filter(allArgs, func(arg DbRuleConditionOrActionArgument, i int) bool {
			return arg.ConditionId.Valid && arg.ConditionId.Int32 == root.Id
		})
		cond = engine.Condition{
			Id:   int(root.Id),
			Fn:   engine.CondFn(root.FunctionType.Int32),
			Args: BuildArguments(args),
		}
	} else {
		// recursively build list nodes
		list := []engine.Condition{}
		children := lo.Filter(conditions, func(c DbRuleCondition, i int) bool {
			return c.ParentConditionId.Valid && c.ParentConditionId.Int32 == rootConditionId
		})
		for _, child := range children {
			list = append(list, BuildCondition(child.Id, conditions, allArgs))
		}
		cond = engine.Condition{
			Id:   int(root.Id),
			List: list,
			Or:   root.LogicOr.Int32 == 1,
		}
	}
	return
}

func BuildActions(
	ruleId int32,
	allRuleActions []DbRuleAction,
	allArgs []DbRuleConditionOrActionArgument,
) (result []engine.Action) {
	actions := lo.Filter(allRuleActions, func(a DbRuleAction, i int) bool {
		return a.RuleId == ruleId
	})
	for _, a := range actions {
		args := lo.Filter(allArgs, func(arg DbRuleConditionOrActionArgument, i int) bool {
			return arg.ActionId.Valid && arg.ActionId.Int32 == a.Id
		})
		result = append(result, engine.Action{
			Id:      int(a.Id),
			Fn:      engine.ActionFn(a.FunctionType.Int32),
			Args:    BuildArguments(args),
			Mapping: engine.Mapping{},
		})
	}
	return
}

func BuildArguments(args []DbRuleConditionOrActionArgument) (result engine.Args) {
	result = engine.Args{}
	lists := make(map[string][]any)
	for _, a := range args {
		islist := a.IsList.Valid && a.IsList.Int32 == 1
		var value any
		if a.Value.Valid {
			value = a.Value.String
		} else if a.DeviceId.Valid {
			value = engine.DeviceId(a.DeviceId.String)
		} else if a.DeviceClassId.Valid {
			value = engine.DeviceClass(a.DeviceClassId.Int32)
		} else {
			panic("unexpected conditions")
		}
		if islist {
			lists[a.ArgumentName] = append(lists[a.ArgumentName], value)
		} else {
			result[a.ArgumentName] = value
		}
	}
	for k, v := range lists {
		result[k] = v
	}
	return
}

func ToDb(inrule engine.Rule, seq utils.Seq) (
	DbRule,
	[]DbRuleCondition,
	[]DbRuleConditionOrActionArgument,
	error,
) {
	outrule := ToDbRule(inrule, seq)
	args := make([]DbRuleConditionOrActionArgument, 0)
	outconds := ToDbConditions(outrule.Id, nil, inrule.Condition, seq, &args)
	return outrule, outconds, args, nil
}

func ToDbRule(rule engine.Rule, seq utils.Seq) DbRule {
	return DbRule{
		Id:         int32(seq.Next()),
		Comments:   rule.Comments,
		IsDisabled: db.NewNullInt32FromBool(rule.Disabled),
		Throttle:   db.NewNullInt32(int32(rule.Throttle.Seconds())),
	}
}

func ToDbConditions(
	ruleId int32,
	parent *DbRuleCondition,
	condition engine.Condition,
	seq utils.Seq,
	args *[]DbRuleConditionOrActionArgument,
) (res []DbRuleCondition) {
	withList := len(condition.List) > 0
	withFn := condition.Fn > 0
	if !withList && !withFn {
		return
	}
	if withList && withFn {
		panic("unexpected conditions")
	}
	if withList {
		cond := DbRuleCondition{
			Id:      int32(seq.Next()),
			RuleId:  ruleId,
			LogicOr: db.NewNullInt32FromBool(condition.Or),
		}
		if parent != nil {
			cond.ParentConditionId = db.NewNullInt32(parent.Id)
		}
		res = append(res, cond)
		for _, childIn := range condition.List {
			res = append(res, ToDbConditions(ruleId, &cond, childIn, seq, args)...)
		}
	} else if withFn {
		node := DbRuleCondition{
			Id:           int32(seq.Next()),
			RuleId:       ruleId,
			FunctionType: db.NewNullInt32(int32(condition.Fn)),
		}
		if parent != nil {
			node.ParentConditionId = db.NewNullInt32(parent.Id)
		}
		res = append(res, node)
		for key, value := range condition.Args {
			*args = append(*args, ToDbArgument(
				node.Id,
				key,
				value,
				seq,
			))
		}
	}
	return
}

// TODO handle all fields
func ToDbArgument(condId int32, key string, value any, seq utils.Seq) DbRuleConditionOrActionArgument {
	return DbRuleConditionOrActionArgument{
		Id:           int32(seq.Next()),
		ConditionId:  db.NewNullInt32(condId),
		ArgumentName: key,
		Value:        db.NewNullString(fmt.Sprintf("%v", value)),
	}
}
