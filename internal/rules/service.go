package rules

import (
	"database/sql"
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/samber/lo"
)

type RulesService interface {
	GetOne(ruleId int32) (engine.Rule, error)
	Get() ([]engine.Rule, error)
	Create(rule engine.Rule) (int64, error)
}

type rulesService struct {
	repository RulesRepository
}

func NewService(r RulesRepository) RulesService {
	return rulesService{
		repository: r,
	}
}

func (s rulesService) Create(rule engine.Rule) (int64, error) {
	dbRule, dbConditions, dbActions, dbArguments, dbMappings := ToDb(
		rule,
		utils.NewSeq(),
	)
	return s.repository.Create(
		dbRule,
		dbConditions,
		dbActions,
		dbArguments,
		dbMappings,
	)
}

func (s rulesService) GetOne(ruleId int32) (res engine.Rule, err error) {
	rules,
		conditions,
		ruleActions,
		ruleConditionOrActionArguments,
		ruleActionArgumentMappings,
		err := s.repository.Get(db.NewNullInt32(ruleId))
	if err != nil {
		return
	}
	if len(rules) == 0 {
		err = fmt.Errorf("no such rule")
		return
	}
	erules := Build(
		rules,
		conditions,
		ruleActions,
		ruleConditionOrActionArguments,
		ruleActionArgumentMappings,
	)
	res = erules[0]
	return
}

func (s rulesService) Get() ([]engine.Rule, error) {
	rules,
		conditions,
		ruleActions,
		ruleConditionOrActionArguments,
		ruleActionArgumentMappings,
		err := s.repository.Get(sql.NullInt32{})
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

// takes flat db representaion of buils hierarchic [engine.Rule]
// (opposite to ToDb)
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
			Actions:   BuildActions(r.Id, allRuleActions, allArgs, allMappings),
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
	allMappings []DbRuleActionArgumentMapping,
) (result []engine.Action) {
	actions := lo.Filter(allRuleActions, func(a DbRuleAction, i int) bool {
		return a.RuleId == ruleId
	})
	for _, a := range actions {
		args := lo.Filter(allArgs, func(arg DbRuleConditionOrActionArgument, i int) bool {
			return arg.ActionId.Valid && arg.ActionId.Int32 == a.Id
		})
		mappings := lo.Filter(allMappings, func(mapping DbRuleActionArgumentMapping, i int) bool {
			return lo.SomeBy(args, func(arg DbRuleConditionOrActionArgument) bool {
				return arg.Id == mapping.ArgumentId
			})
		})
		result = append(result, engine.Action{
			Id:   int(a.Id),
			Fn:   engine.ActionFn(a.FunctionType.Int32),
			Args: BuildArguments(args),
			// DeviceId: engine.DeviceId(a.DeviceId.String),
			Mapping: BuildMappings(mappings, args),
		})
	}
	return
}

func BuildMappings(
	mappings []DbRuleActionArgumentMapping,
	args []DbRuleConditionOrActionArgument,
) (result engine.Mapping) {
	result = engine.Mapping{}
	for _, mapping := range mappings {
		arg, _ := lo.Find(args, func(arg DbRuleConditionOrActionArgument) bool {
			return arg.Id == mapping.ArgumentId
		})
		_, exist := result[arg.ArgumentName]
		if !exist {
			result[arg.ArgumentName] = make(map[string]string)
		}
		result[arg.ArgumentName][mapping.Key] = mapping.Value
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
			value = types.DeviceId(a.DeviceId.String)
		} else if a.DeviceClassId.Valid {
			value = types.DeviceClass(a.DeviceClassId.Int32)
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

// transforms [engine.Rule] to flat representtion for db
// (opposite to Build)
func ToDb(inrule engine.Rule, seq utils.Seq) (
	DbRule,
	[]DbRuleCondition,
	[]DbRuleAction,
	[]DbRuleConditionOrActionArgument,
	[]DbRuleActionArgumentMapping,
) {
	outrule := ToDbRule(inrule, seq)
	args := make([]DbRuleConditionOrActionArgument, 0)
	mappings := make([]DbRuleActionArgumentMapping, 0)
	outconds := ToDbConditions(outrule.Id, nil, inrule.Condition, seq, &args)
	outactions := ToDbActions(outrule.Id, inrule.Actions, seq, &args, &mappings)
	return outrule, outconds, outactions, args, mappings
}

func ToDbRule(rule engine.Rule, seq utils.Seq) DbRule {
	return DbRule{
		Id:         int32(seq.Next()),
		Comments:   rule.Comments,
		IsDisabled: db.NewNullInt32FromBool(rule.Disabled),
		Throttle:   db.NewNullInt32(int32(rule.Throttle.Seconds())),
	}
}

func ToDbActions(
	ruleId int32,
	actions []engine.Action,
	seq utils.Seq,
	args *[]DbRuleConditionOrActionArgument,
	mappings *[]DbRuleActionArgumentMapping,
) (res []DbRuleAction) {
	for _, action := range actions {
		// deviceId := sql.NullString{
		// 	String: string(action.DeviceId),
		// 	Valid:  len(action.DeviceId) > 0,
		// }
		node := DbRuleAction{
			Id:           int32(seq.Next()),
			RuleId:       ruleId,
			FunctionType: db.NewNullInt32(int32(action.Fn)),
			// DeviceId:     deviceId,
		}
		res = append(res, node)
		argNameToId := make(map[string]int32, len(action.Args))
		for key, value := range action.Args {
			newargs := ToDbArguments(
				ruleId,
				nil,
				&node,
				key,
				value,
				seq,
				false,
			)
			for _, arg := range newargs {
				argNameToId[arg.ArgumentName] = arg.Id
			}
			*args = append(*args, newargs...)
		}
		for argName, argMapping := range action.Mapping {
			for key, value := range argMapping {
				*mappings = append(*mappings, DbRuleActionArgumentMapping{
					Id:         int32(seq.Next()),
					RuleId:     ruleId,
					ArgumentId: argNameToId[argName],
					Key:        key,
					Value:      value,
				})
			}
		}
	}
	return
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
			newargs := ToDbArguments(
				ruleId,
				&node,
				nil,
				key,
				value,
				seq,
				false,
			)
			*args = append(*args, newargs...)
		}
	}
	return
}

func ToDbArguments(
	ruleId int32,
	condition *DbRuleCondition,
	action *DbRuleAction,
	key string,
	value any,
	seq utils.Seq,
	islist bool,
) (res []DbRuleConditionOrActionArgument) {
	// fmt.Printf("%v, %v, %T\n", key, value, value)
	if listArg, ok := value.([]any); ok {
		for _, vi := range listArg {
			args := ToDbArguments(ruleId, condition, action, key, vi, seq, true)
			res = append(res, args...)
		}
	} else {
		arg := DbRuleConditionOrActionArgument{
			Id:           int32(seq.Next()),
			RuleId:       ruleId,
			ArgumentName: key,
			IsList:       db.NewNullInt32FromBool(islist),
		}
		if condition != nil {
			arg.ConditionId = db.NewNullInt32(condition.Id)
		}
		if action != nil {
			arg.ActionId = db.NewNullInt32(action.Id)
		}
		if deviceId, ok := value.(types.DeviceId); ok {
			arg.DeviceId = db.NewNullString(string(deviceId))
		} else if deviceClass, ok := value.(types.DeviceClass); ok {
			arg.DeviceClassId = db.NewNullInt32(int32(deviceClass))
		} else {
			arg.Value = db.NewNullString(fmt.Sprintf("%v", value))
		}
		res = append(res, arg)
	}
	return
}

//
