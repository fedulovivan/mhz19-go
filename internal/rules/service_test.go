package rules

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/stretchr/testify/suite"
)

// func dump(name string, in any) {
// 	json, err := json.MarshalIndent(in, "", "  ")
// 	if err != nil {
// 		fmt.Println(name, err, in)
// 		return
// 	}
// 	fmt.Println(name, string(json))
// }

type MappingsSuite struct {
	suite.Suite
}

type TableRow struct {
	rules                          []DbRule
	conditions                     []DbRuleCondition
	ruleActions                    []DbRuleAction
	ruleConditionOrActionArguments []DbRuleConditionOrActionArgument
	ruleActionArgumentMappings     []DbRuleActionArgumentMapping
	expectedLen                    int
	expectedJson                   string
}

var testDataTable = []TableRow{
	// case 1
	{
		// no rules, conditions, anything
		expectedJson: "null",
	},
	// case 2
	{
		rules: []DbRule{
			{
				Id:         1,
				Comments:   "case 2",
				IsDisabled: db.NewNullInt32(0),
				Throttle:   db.NewNullInt32(0),
			},
		},
		conditions: []DbRuleCondition{
			{
				Id:           1,
				RuleId:       1,
				FunctionType: db.NewNullInt32(1),
			},
		},
		ruleConditionOrActionArguments: []DbRuleConditionOrActionArgument{
			{
				Id:           1,
				ConditionId:  db.NewNullInt32(1),
				ArgumentName: "Left",
				Value:        db.NewNullString("foo"),
			},
			{
				Id:           2,
				ConditionId:  db.NewNullInt32(1),
				ArgumentName: "Right",
				Value:        db.NewNullString("bar"),
			},
			{
				Id:           3,
				ConditionId:  db.NewNullInt32(1),
				ArgumentName: "Third",
				Value:        db.NewNullString("baz"),
				IsList:       db.NewNullInt32(1),
			},
			{
				Id:            4,
				ConditionId:   db.NewNullInt32(1),
				ArgumentName:  "Fourth",
				DeviceClassId: db.NewNullInt32(1),
				IsList:        db.NewNullInt32(1),
			},
			{
				Id:           5,
				ConditionId:  db.NewNullInt32(1),
				ArgumentName: "Fifth",
				DeviceId:     db.NewNullString("0x00158d0004244bda"),
				IsList:       db.NewNullInt32(1),
			},
			{
				Id:           6,
				ActionId:     db.NewNullInt32(1),
				ArgumentName: "Value",
				Value:        db.NewNullString("$message.action"),
			},
			{
				Id:           7,
				ActionId:     db.NewNullInt32(1),
				ArgumentName: "DeviceId",
				DeviceId:     db.NewNullString("10011cec96"),
			},
		},
		ruleActions: []DbRuleAction{
			{
				Id:           1,
				RuleId:       1,
				FunctionType: db.NewNullInt32(1),
				DeviceId:     db.NewNullString("0x00158d0004244bda"),
			},
		},
		expectedLen:  1,
		expectedJson: `[{"id":1,"disabled":false,"comments":"case 2","condition":{"fn":1,"args":{"Fifth":["0x00158d0004244bda"],"Fourth":[1],"Left":"foo","Right":"bar","Third":["baz"]}},"actions":[{"fn":1,"args":{"DeviceId":"10011cec96","Value":"$message.action"},"mapping":{}}],"throttle":0}]`,
	},
	// case 3
	{
		rules: []DbRule{
			{
				Id:         2,
				Comments:   "case 3",
				IsDisabled: db.NewNullInt32(0),
				Throttle:   db.NewNullInt32(0),
			},
		},
		conditions: []DbRuleCondition{
			{
				Id:     22,
				RuleId: 2,
			},
			{
				Id:                23,
				RuleId:            2,
				FunctionType:      db.NewNullInt32(1),
				ParentConditionId: db.NewNullInt32(22),
			},
			{
				Id:                24,
				RuleId:            2,
				FunctionType:      db.NewNullInt32(2),
				ParentConditionId: db.NewNullInt32(22),
			},
			{
				Id:                25,
				RuleId:            2,
				ParentConditionId: db.NewNullInt32(22),
				LogicOr:           db.NewNullInt32(1),
			},
			{
				Id:                26,
				RuleId:            2,
				FunctionType:      db.NewNullInt32(3),
				ParentConditionId: db.NewNullInt32(25),
			},
			{
				Id:                27,
				RuleId:            2,
				FunctionType:      db.NewNullInt32(4),
				ParentConditionId: db.NewNullInt32(25),
			},
		},
		expectedJson: `[{"id":2,"disabled":false,"comments":"case 3","condition":{"list":[{"fn":1},{"fn":2},{"list":[{"fn":3},{"fn":4}],"or":true}]},"actions":null,"throttle":0}]`,
	},
}

func (s *MappingsSuite) Test10() {
	for _, row := range testDataTable {
		// if i != 2 {
		// 	continue
		// }
		result := Build(
			row.rules,
			row.conditions,
			row.ruleActions,
			row.ruleConditionOrActionArguments,
			row.ruleActionArgumentMappings,
		)
		json, _ := json.Marshal(result)
		s.JSONEq(string(json), row.expectedJson)
		// fmt.Println(string(json))
		// json, _ := json.MarshalIndent(result, "", "  ")
		// s.Len(result, row.expectedLen)
	}
}

func (s *MappingsSuite) Test20() {
	defer func() { _ = recover() }()
	BuildArguments([]DbRuleConditionOrActionArgument{
		{},
	})
	s.Fail("expected to panic")
}

func (s *MappingsSuite) Test30() {
	res := BuildCondition(0, []DbRuleCondition{}, []DbRuleConditionOrActionArgument{})
	s.Zero(res)
}

func (s *MappingsSuite) Test40() {
	res := BuildCondition(
		1,
		[]DbRuleCondition{{}},
		[]DbRuleConditionOrActionArgument{},
	)
	s.Zero(res)
}

func (s *MappingsSuite) Test50() {
	// defer func() { _ = recover() }()
	ToDbConditions(1, nil, engine.Condition{}, nil, nil)
	// s.Fail("expected to panic")
}

func (s *MappingsSuite) Test51() {
	defer func() { _ = recover() }()
	ToDbConditions(1, nil, engine.Condition{
		Fn:   engine.COND_CHANGED,
		List: []engine.Condition{{Fn: engine.COND_EQUAL}},
	}, nil, nil)
	s.Fail("expected to panic")
}

func (s *MappingsSuite) Test52() {
	actualArgs := []DbRuleConditionOrActionArgument{}
	actualConds := ToDbConditions(1, nil, engine.Condition{
		Or: true,
		List: []engine.Condition{
			{Fn: engine.COND_EQUAL, Args: engine.Args{"Left": 1, "Right": 2}},
			{Fn: engine.COND_IN_LIST},
		},
	}, utils.NewSeq(), &actualArgs)
	// dump(actualArgs)
	// dump(actualConds)
	expectedArgs := "[{3 {2 true} {0 false} Left {0 false} {1 true} { false} {0 false}} {4 {2 true} {0 false} Right {0 false} {2 true} { false} {0 false}}]"
	s.Equal(expectedArgs, fmt.Sprintf("%v", actualArgs))
	expectedConds := "[{1 1 {0 false} {1 true} {0 false}} {2 1 {2 true} {0 false} {1 true}} {5 1 {3 true} {0 false} {1 true}}]"
	s.Equal(expectedConds, fmt.Sprintf("%v", actualConds))
}

func (s *MappingsSuite) Test53() {
	actual := ToDbConditions(1, nil, engine.Condition{
		Or: true,
		List: []engine.Condition{
			{Fn: engine.COND_EQUAL},
			{List: []engine.Condition{
				{Fn: engine.COND_NOT_EQUAL},
				{Fn: engine.COND_IN_LIST},
			}},
		},
	}, utils.NewSeq(), nil)
	// dump(actual)
	expected := "[{1 1 {0 false} {1 true} {0 false}} {2 1 {2 true} {0 false} {1 true}} {3 1 {0 false} {0 true} {1 true}} {4 1 {4 true} {0 false} {3 true}} {5 1 {3 true} {0 false} {3 true}}]"
	s.Equal(expected, fmt.Sprintf("%v", actual))
}

func (s *MappingsSuite) Test60() {
	defer func() { _ = recover() }()
	_, _, _, err := ToDb(engine.Rule{}, nil)
	s.Nil(err)
	s.Fail("expected to panic")
}

func (s *MappingsSuite) Test61() {
	_, _, _, err := ToDb(engine.Rule{}, utils.NewSeq())
	s.Nil(err)
}

func (s *MappingsSuite) Test62() {
	inrule := engine.Rule{}
	outrule, outconds, outargs, err := ToDb(inrule, utils.NewSeq())
	// dump("outrule", outrule)
	// dump("outconds", outconds)
	// dump("outargs", outargs)
	// dump("err", err)
	expected := "{1  {0 true} {0 true}}"
	s.Equal(expected, fmt.Sprintf("%v", outrule))
	s.Len(outconds, 0)
	s.Len(outargs, 0)
	s.Nil(err)
}

func (s *MappingsSuite) Test63() {
	inrule := engine.Rule{
		Comments: "unit test",
		Disabled: true,
		Condition: engine.Condition{
			Fn: engine.COND_EQUAL, Args: engine.Args{"One": 1, "Two": 2},
		},
		Actions:  []engine.Action{},
		Throttle: 100500,
	}
	outrule, outconds, outargs, err := ToDb(inrule, utils.NewSeq())

	// dump("outrule", outrule)
	// dump("outconds", outconds)
	// dump("outargs", outargs)
	// dump("err", err)

	expectedRule := "{1 unit test {1 true} {0 true}}"
	s.Equal(expectedRule, fmt.Sprintf("%v", outrule))

	expectedConds := "[{2 1 {2 true} {0 false} {0 false}}]"
	s.Equal(expectedConds, fmt.Sprintf("%v", outconds))

	// TODO reduce expectedArgs size, which we have to use here, since json.Marshal(...) gives sorted keys unlike fmt.Sprintf("%v"...)
	expectedArgs := `[{"Id":3,"ConditionId":{"Int32":2,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"One","IsList":{"Int32":0,"Valid":false},"Value":{"String":"1","Valid":true},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":0,"Valid":false}},{"Id":4,"ConditionId":{"Int32":2,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"Two","IsList":{"Int32":0,"Valid":false},"Value":{"String":"2","Valid":true},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":0,"Valid":false}}]`
	outargsj, _ := json.Marshal(outargs)
	s.Equal(expectedArgs, string(outargsj))
	// fmt.Println(string(outargsj))
	// expectedArgs := "[{3 {2 true} {0 false} Two {0 false} {2 true} { false} {0 false}} {4 {2 true} {0 false} One {0 false} {1 true} { false} {0 false}}]"
	// s.Equal(expectedArgs, fmt.Sprintf("%v", outargs))

	s.Nil(err)

}

func TestMappings(t *testing.T) {
	suite.Run(t, new(MappingsSuite))
}
