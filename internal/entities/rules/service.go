package rules

import (
	"database/sql"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/samber/lo"
)

type rulesService struct {
	created    chan types.Rule
	deleted    chan int
	repository RulesRepository
}

var instance types.RulesService

var knownSimpleTypes = []string{"string", "int", "float64", "bool"}

func NewService(r RulesRepository) types.RulesService {
	return rulesService{
		created:    make(chan types.Rule, 100),
		deleted:    make(chan int, 100),
		repository: r,
	}
}

func ServiceSingleton(r RulesRepository) types.RulesService {
	if instance == nil {
		instance = NewService(r)
	}
	return instance
}

func (s rulesService) OnCreated() chan types.Rule {
	return s.created
}

func (s rulesService) OnDeleted() chan int {
	return s.deleted
}

func (s rulesService) Create(rule types.Rule) (newRuleId int64, err error) {
	dbRule, dbConditions, dbActions, dbArguments, dbMappings := ToDb(
		rule,
		utils.NewSeq(0),
	)
	newRuleId, err = s.repository.Create(
		dbRule,
		dbConditions,
		dbActions,
		dbArguments,
		dbMappings,
	)
	rule.Id = int(newRuleId)
	if err == nil {
		s.created <- rule
	}
	return
}

func (s rulesService) Delete(ruleId int) (err error) {
	err = s.repository.Delete(int32(ruleId))
	if err == nil {
		s.deleted <- ruleId
	}
	return
}

func (s rulesService) GetOne(ruleId int) (res types.Rule, err error) {
	rules,
		conditions,
		ruleActions,
		ruleConditionOrActionArguments,
		ruleActionArgumentMappings,
		err := s.repository.Get(db.NewNullInt32(int32(ruleId)))
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

func (s rulesService) Get() ([]types.Rule, error) {
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

// takes flat db representaion of buils hierarchic [types.Rule]
// (opposite to ToDb)
func Build(
	allRules []DbRule,
	allConditions []DbRuleCondition,
	allRuleActions []DbRuleAction,
	allArgs []DbRuleConditionOrActionArgument,
	allMappings []DbRuleActionArgumentMapping,
) (result []types.Rule) {
	for _, r := range allRules {
		rootCond, rootCondFound := lo.Find(allConditions, func(c DbRuleCondition) bool {
			return c.RuleId == r.Id && !c.ParentConditionId.Valid
		})
		cond := types.Condition{}
		if rootCondFound {
			cond = BuildCondition(rootCond.Id, allConditions, allArgs)
		}
		var throttle types.Throttle
		if r.ThrottleMs.Valid {
			throttle = types.Throttle{
				Duration: time.Duration(r.ThrottleMs.Int32) * time.Millisecond,
			}
		}
		rule := types.Rule{
			Id:        int(r.Id),
			Name:      r.Name,
			Disabled:  r.IsDisabled.Int32 == 1,
			Condition: cond,
			Throttle:  throttle,
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
) (cond types.Condition) {
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
		cond = types.Condition{
			Id:   int(root.Id),
			Fn:   types.CondFn(root.FunctionType.Int32),
			Args: BuildArguments(args),
			Not:  root.Not.Int32 == 1,
		}
		if root.OtherDeviceId.Valid {
			cond.OtherDeviceId = types.DeviceId(root.OtherDeviceId.String)
		}
	} else {
		// recursively build nested nodes
		nested := []types.Condition{}
		children := lo.Filter(conditions, func(c DbRuleCondition, i int) bool {
			return c.ParentConditionId.Valid && c.ParentConditionId.Int32 == rootConditionId
		})
		for _, child := range children {
			nested = append(nested, BuildCondition(child.Id, conditions, allArgs))
		}
		cond = types.Condition{
			Id:     int(root.Id),
			Nested: nested,
			Or:     root.LogicOr.Int32 == 1,
		}
	}
	return
}

func BuildActions(
	ruleId int32,
	allRuleActions []DbRuleAction,
	allArgs []DbRuleConditionOrActionArgument,
	allMappings []DbRuleActionArgumentMapping,
) (result []types.Action) {
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
		result = append(result, types.Action{
			Id:      int(a.Id),
			Fn:      types.ActionFn(a.FunctionType.Int32),
			Args:    BuildArguments(args),
			Mapping: BuildMappings(mappings, args),
		})
	}
	return
}

func BuildMappings(
	mappings []DbRuleActionArgumentMapping,
	args []DbRuleConditionOrActionArgument,
) (result types.Mapping) {
	result = types.Mapping{}
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

func BuildArguments(args []DbRuleConditionOrActionArgument) (result types.Args) {
	lists := make(map[string][]any)
	for _, a := range args {
		islist := a.IsList.Valid && a.IsList.Int32 == 1
		var value any
		if a.Value.Valid && a.ValueDataType.Valid {
			switch a.ValueDataType.String {
			case "string":
				value = a.Value.String
			case "int":
				value, _ = strconv.Atoi(a.Value.String)
			case "float64":
				value, _ = strconv.ParseFloat(a.Value.String, 64)
			case "bool":
				value = a.Value.String == "true"
			default:
				panic(fmt.Sprintf("unexpected value data type %s", a.ValueDataType.String))
			}
		} else if a.DeviceId.Valid {
			value = types.DeviceId(a.DeviceId.String)
		} else if a.DeviceClassId.Valid {
			value = types.DeviceClass(a.DeviceClassId.Int32)
		} else if a.ChannelTypeId.Valid {
			value = types.ChannelType(a.ChannelTypeId.Int32)
		} else {
			panic("unexpected conditions")
		}
		if islist {
			lists[a.ArgumentName] = append(lists[a.ArgumentName], value)
		} else {
			if result == nil {
				result = make(types.Args)
			}
			result[a.ArgumentName] = value
		}
	}
	for k, v := range lists {
		if result == nil {
			result = make(types.Args)
		}
		result[k] = v
	}
	return
}

// transforms [types.Rule] to flat representtion for db
// (opposite to Build)
func ToDb(inrule types.Rule, seq utils.Seq) (
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

func ToDbRule(rule types.Rule, seq utils.Seq) DbRule {
	return DbRule{
		Id:         int32(seq.Inc()),
		Name:       rule.Name,
		IsDisabled: db.NewNullInt32FromBool(rule.Disabled),
		ThrottleMs: db.NewNullInt32(int32(rule.Throttle.Duration.Milliseconds())),
	}
}

func ToDbActions(
	ruleId int32,
	actions []types.Action,
	seq utils.Seq,
	args *[]DbRuleConditionOrActionArgument,
	mappings *[]DbRuleActionArgumentMapping,
) (res []DbRuleAction) {
	for _, action := range actions {
		node := DbRuleAction{
			Id:           int32(seq.Inc()),
			RuleId:       ruleId,
			FunctionType: db.NewNullInt32(int32(action.Fn)),
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
					Id:         int32(seq.Inc()),
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
	condition types.Condition,
	seq utils.Seq,
	args *[]DbRuleConditionOrActionArgument,
) (res []DbRuleCondition) {
	withList := len(condition.Nested) > 0
	withFn := condition.Fn > 0
	if !withList && !withFn {
		return
	}
	if withList && withFn {
		panic("unexpected conditions")
	}
	if withList {
		cond := DbRuleCondition{
			Id:      int32(seq.Inc()),
			RuleId:  ruleId,
			LogicOr: db.NewNullInt32FromBool(condition.Or),
		}
		if parent != nil {
			cond.ParentConditionId = db.NewNullInt32(parent.Id)
		}
		res = append(res, cond)
		for _, childIn := range condition.Nested {
			res = append(res, ToDbConditions(ruleId, &cond, childIn, seq, args)...)
		}
	} else if withFn {
		node := DbRuleCondition{
			Id:           int32(seq.Inc()),
			RuleId:       ruleId,
			FunctionType: db.NewNullInt32(int32(condition.Fn)),
			Not:          db.NewNullInt32FromBool(condition.Not),
		}
		if len(condition.OtherDeviceId) > 0 {
			node.OtherDeviceId = db.NewNullString(string(condition.OtherDeviceId))
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
	if listArg, ok := value.([]any); ok {
		for _, vi := range listArg {
			res = append(
				res,
				ToDbArguments(ruleId, condition, action, key, vi, seq, true)...,
			)
		}
	} else {
		arg := DbRuleConditionOrActionArgument{
			Id:           int32(seq.Inc()),
			RuleId:       ruleId,
			ArgumentName: key,
			IsList:       db.NewNullInt32FromBool(islist),
		}
		valueDataType := fmt.Sprintf("%T", value)
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
		} else if channelType, ok := value.(types.ChannelType); ok {
			arg.ChannelTypeId = db.NewNullInt32(int32(channelType))
		} else if slices.Contains(knownSimpleTypes, valueDataType) {
			arg.Value = db.NewNullString(fmt.Sprintf("%v", value))
			arg.ValueDataType = db.NewNullString(valueDataType)
		} else {
			panic(fmt.Sprintf("unexpected value data type %s", valueDataType))
		}
		res = append(res, arg)
	}
	return
}

//
