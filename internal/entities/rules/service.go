package rules

import (
	"database/sql"
	"fmt"
	"slices"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/counters"
	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

type service struct {
	created    chan types.Rule
	deleted    chan int
	repository RulesRepository
}

var instance *service

var knownSimpleTypes = []string{"string", "int", "float64", "bool"}

func NewService(r RulesRepository) *service {
	return &service{
		created:    make(chan types.Rule),
		deleted:    make(chan int),
		repository: r,
	}
}

func ServiceSingleton(r RulesRepository) *service {
	if instance == nil {
		instance = NewService(r)
	}
	return instance
}

func (s service) OnCreated() <-chan types.Rule {
	return s.created
}

func (s service) OnDeleted() <-chan int {
	return s.deleted
}

func (s service) Create(rule types.Rule) (newRuleId int64, err error) {
	seq := &atomic.Int32{}
	dbRule, dbConditions, dbActions, dbArguments, dbMappings := ToDb(
		rule,
		seq,
	)
	newRuleId, err = s.repository.Create(
		dbRule,
		dbConditions,
		dbActions,
		dbArguments,
		dbMappings,
	)
	if err == nil {
		new, _ := s.GetOne(int(newRuleId))
		s.created <- new
	}
	return
}

func (s service) Delete(ruleId int) (err error) {
	err = s.repository.Delete(int32(ruleId))
	if err == nil {
		s.deleted <- ruleId
	}
	return
}

func (s service) GetOne(ruleId int) (res types.Rule, err error) {
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

func (s service) Get() ([]types.Rule, error) {
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

// takes flat db representation and build hierarchic structure of types.Rule
// (opposite to ToDb)
func Build(
	allRules []DbRule,
	allConditions []DbRuleCondition,
	allRuleActions []DbRuleAction,
	allArgs []DbRuleConditionOrActionArgument,
	allMappings []DbRuleActionArgumentMapping,
) (result []types.Rule) {

	defer func(t *prometheus.Timer) {
		t.ObserveDuration()
	}(prometheus.NewTimer(counters.BuildRules))

	actionsByRuleId := make(map[int32][]DbRuleAction, len(allRules))
	for _, action := range allRuleActions {
		actionsByRuleId[action.RuleId] = append(actionsByRuleId[action.RuleId], action)
	}

	argsByActionId := make(map[int32][]DbRuleConditionOrActionArgument, len(allRuleActions))
	argsByConditionId := make(map[int32][]DbRuleConditionOrActionArgument /* , len() */)
	for _, arg := range allArgs {
		if arg.ActionId.Valid {
			argsByActionId[arg.ActionId.Int32] = append(argsByActionId[arg.ActionId.Int32], arg)
		}
		if arg.ConditionId.Valid {
			argsByConditionId[arg.ConditionId.Int32] = append(argsByConditionId[arg.ConditionId.Int32], arg)
		}
	}

	mappingsByArgumentId := make(map[int32][]DbRuleActionArgumentMapping)
	for _, mapping := range allMappings {
		mappingsByArgumentId[mapping.ArgumentId] = append(mappingsByArgumentId[mapping.ArgumentId], mapping)
	}

	conditionsByRuleAndParent := make(map[int32](map[int32][]DbRuleCondition))
	for _, condition := range allConditions {
		if _, exist := conditionsByRuleAndParent[condition.RuleId]; !exist {
			conditionsByRuleAndParent[condition.RuleId] = map[int32][]DbRuleCondition{}
		}
		parent := int32(-1)
		if condition.ParentConditionId.Valid {
			parent = condition.ParentConditionId.Int32
		}
		conditionsByRuleAndParent[condition.RuleId][parent] = append(conditionsByRuleAndParent[condition.RuleId][parent], condition)
	}

	for _, r := range allRules {

		сonditionsByParent := conditionsByRuleAndParent[r.Id]

		if len(сonditionsByParent[-1]) != 1 {
			panic("unexpected conditions")
		}

		rootCond := сonditionsByParent[-1][0]

		var throttle types.Throttle
		if r.ThrottleMs.Valid {
			throttle = types.Throttle{
				Duration: time.Duration(r.ThrottleMs.Int32) * time.Millisecond,
			}
		}
		rule := types.Rule{
			Id:       int(r.Id),
			Name:     r.Name,
			Disabled: r.IsDisabled.Int32 == 1,
			Throttle: throttle,
			Condition: BuildCondition(
				rootCond,
				сonditionsByParent,
				argsByConditionId,
			),
			Actions: BuildActions(
				r.Id,
				actionsByRuleId,
				argsByActionId,
				mappingsByArgumentId,
			),
		}
		result = append(result, rule)
	}
	return
}

func BuildCondition(
	root DbRuleCondition,
	сonditionsByParent map[int32][]DbRuleCondition,
	argsByConditionId map[int32][]DbRuleConditionOrActionArgument,
) (cond types.Condition) {

	isFn := root.FunctionType.Valid
	if isFn {
		// build function node
		args := argsByConditionId[root.Id]
		cond = types.Condition{
			Id:       int(root.Id),
			Fn:       types.CondFn(root.FunctionType.Int32),
			Args:     BuildArguments(args),
			Not:      root.Not.Int32 == 1,
			Disabled: root.IsDisabled.Int32 == 1,
		}
		if root.OtherDeviceId.Valid {
			cond.OtherDeviceId = types.DeviceId(root.OtherDeviceId.String)
		}
	} else {
		// recursively build nested nodes
		nested := []types.Condition{}
		children := сonditionsByParent[root.Id]

		for _, child := range children {
			nested = append(nested, BuildCondition(child, сonditionsByParent, argsByConditionId))
		}
		cond = types.Condition{
			Id:       int(root.Id),
			Nested:   nested,
			Or:       root.LogicOr.Int32 == 1,
			Disabled: root.IsDisabled.Int32 == 1,
		}
	}
	return
}

func BuildActions(
	ruleId int32,
	actionsByRuleId map[int32][]DbRuleAction,
	argsByActionId map[int32][]DbRuleConditionOrActionArgument,
	mappingsByArgumentId map[int32][]DbRuleActionArgumentMapping,
) (result []types.Action) {
	actions := actionsByRuleId[ruleId]
	for _, action := range actions {
		actionArgs := argsByActionId[action.Id]
		result = append(result, types.Action{
			Id:       int(action.Id),
			Fn:       types.ActionFn(action.FunctionType.Int32),
			Args:     BuildArguments(actionArgs),
			Mapping:  BuildMappings(actionArgs, mappingsByArgumentId),
			Disabled: action.IsDisabled.Int32 == 1,
		})
	}
	return
}

func BuildMappings(
	actionArgs []DbRuleConditionOrActionArgument,
	mappingsByArgumentId map[int32][]DbRuleActionArgumentMapping,
) types.Mapping {
	result := make(types.Mapping)
	for _, arg := range actionArgs {
		mappings := mappingsByArgumentId[arg.Id]
		if len(mappings) > 0 {
			result[arg.ArgumentName] = make(map[string]string, len(mappings))
			for _, mapping := range mappings {
				result[arg.ArgumentName][mapping.Key] = mapping.Value
			}
		}
	}
	return result
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
func ToDb(inrule types.Rule, seq *atomic.Int32) (
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

func ToDbRule(rule types.Rule, seq *atomic.Int32) DbRule {
	return DbRule{
		Id:         seq.Add(1),
		Name:       rule.Name,
		IsDisabled: db.NewNullInt32FromBool(rule.Disabled),
		ThrottleMs: db.NewNullInt32(int32(rule.Throttle.Duration.Milliseconds())),
	}
}

func ToDbActions(
	ruleId int32,
	actions []types.Action,
	seq *atomic.Int32,
	args *[]DbRuleConditionOrActionArgument,
	mappings *[]DbRuleActionArgumentMapping,
) (res []DbRuleAction) {
	for _, action := range actions {
		node := DbRuleAction{
			Id:           seq.Add(1),
			RuleId:       ruleId,
			FunctionType: db.NewNullInt32(int32(action.Fn)),
			IsDisabled:   db.NewNullInt32FromBool(action.Disabled),
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
					Id:         seq.Add(1),
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
	seq *atomic.Int32,
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
			Id:         seq.Add(1),
			RuleId:     ruleId,
			LogicOr:    db.NewNullInt32FromBool(condition.Or),
			IsDisabled: db.NewNullInt32FromBool(condition.Disabled),
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
			Id:           seq.Add(1),
			RuleId:       ruleId,
			FunctionType: db.NewNullInt32(int32(condition.Fn)),
			Not:          db.NewNullInt32FromBool(condition.Not),
			IsDisabled:   db.NewNullInt32FromBool(condition.Disabled),
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
	seq *atomic.Int32,
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
			Id:           seq.Add(1),
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
