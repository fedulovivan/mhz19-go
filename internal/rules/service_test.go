package rules

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
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
	// case 00
	{
		// no rules, conditions, anything
		expectedJson: "null",
	},
	// case 01
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
			},
			{
				Id:            4,
				ConditionId:   db.NewNullInt32(1),
				ArgumentName:  "Fourth",
				DeviceClassId: db.NewNullInt32(1),
			},
			{
				Id:           5,
				ConditionId:  db.NewNullInt32(1),
				ArgumentName: "Fifth",
				DeviceId:     db.NewNullString("0x00158d0004244bda"),
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
				ArgumentName: "ListIds",
				DeviceId:     db.NewNullString("10011cec96"),
				IsList:       db.NewNullInt32FromBool(true),
			},
			{
				Id:           8,
				ActionId:     db.NewNullInt32(1),
				ArgumentName: "ListIds",
				DeviceId:     db.NewNullString("78345aaa67"),
				IsList:       db.NewNullInt32FromBool(true),
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
		ruleActionArgumentMappings: []DbRuleActionArgumentMapping{
			{
				Id:         1,
				ArgumentId: 6,
				Key:        "lorem-3",
				Value:      "dolor-4",
			},
			{
				Id:         2,
				ArgumentId: 6,
				Key:        "sit-5",
				Value:      "amet-6",
			},
		},
		expectedLen:  1,
		expectedJson: `[{"id":1,"comments":"case 2","condition":{"fn":1,"args":{"Fifth":"DeviceId(0x00158d0004244bda)","Fourth":"DeviceClass(1)","Left":"foo","Right":"bar","Third":"baz"}},"actions":[{"fn":1,"args":{"ListIds":["DeviceId(10011cec96)","DeviceId(78345aaa67)"],"Value":"$message.action"},"mapping":{"Value":{"lorem-3":"dolor-4","sit-5":"amet-6"}},"deviceId":"DeviceId(0x00158d0004244bda)"}]}]`,
	},
	// case 02
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
		expectedJson: `[{"id":2,"comments":"case 3","condition":{"list":[{"fn":1},{"fn":2},{"list":[{"fn":3},{"fn":4}],"or":true}]}}]`,
	},
}

func (s *ServiceSuite) Test10() {
	row := testDataTable[0]
	result := Build(
		row.rules,
		row.conditions,
		row.ruleActions,
		row.ruleConditionOrActionArguments,
		row.ruleActionArgumentMappings,
	)
	data, _ := json.Marshal(result)
	s.JSONEq(string(data), row.expectedJson)
}

func (s *ServiceSuite) Test11() {
	row := testDataTable[1]
	result := Build(
		row.rules,
		row.conditions,
		row.ruleActions,
		row.ruleConditionOrActionArguments,
		row.ruleActionArgumentMappings,
	)
	data, _ := json.Marshal(result)
	s.JSONEq(string(data), row.expectedJson)
	// fmt.Println(string(data))
}

func (s *ServiceSuite) Test12() {
	row := testDataTable[2]
	result := Build(
		row.rules,
		row.conditions,
		row.ruleActions,
		row.ruleConditionOrActionArguments,
		row.ruleActionArgumentMappings,
	)
	data, _ := json.Marshal(result)
	s.JSONEq(string(data), row.expectedJson)
	fmt.Println(string(data))
}

// fmt.printf("row: %+v\n", row)
// fmt.printf("result: %+v\n", result)
// fmt.println("data", string(data))
// data2, _ := json.marshal(row)
// fmt.println("data2", string(data2))
// json, _ = json.marshal(row)
// s.jsoneq(string(json), row.expectedjson)
// fmt.Println(string(json))
// json, _ := json.MarshalIndent(result, "", "  ")
// s.Len(result, row.expectedLen)

func (s *ServiceSuite) Test20() {
	defer func() { _ = recover() }()
	BuildArguments([]DbRuleConditionOrActionArgument{
		{},
	})
	s.Fail("expected to panic")
}

func (s *ServiceSuite) Test30() {
	res := BuildCondition(0, []DbRuleCondition{}, []DbRuleConditionOrActionArgument{})
	s.Zero(res)
}

func (s *ServiceSuite) Test40() {
	res := BuildCondition(
		1,
		[]DbRuleCondition{{}},
		[]DbRuleConditionOrActionArgument{},
	)
	s.Zero(res)
}

func (s *ServiceSuite) Test50() {
	ToDbConditions(1, nil, engine.Condition{}, nil, nil)
	// defer func() { _ = recover() }()
	// s.Fail("expected to panic")
}

func (s *ServiceSuite) Test51() {
	defer func() { _ = recover() }()
	ToDbConditions(1, nil, engine.Condition{
		Fn:   engine.COND_CHANGED,
		List: []engine.Condition{{Fn: engine.COND_EQUAL}},
	}, nil, nil)
	s.Fail("expected to panic")
}

func (s *ServiceSuite) Test52() {
	actualArgs := []DbRuleConditionOrActionArgument{}
	actualConds := ToDbConditions(1, nil, engine.Condition{
		Or: true,
		List: []engine.Condition{
			{Fn: engine.COND_EQUAL, Args: engine.Args{"Left": 1, "Right": 2}},
			{Fn: engine.COND_IN_LIST},
		},
	}, utils.NewSeq(), &actualArgs)
	s.Len(actualArgs, 2)
	s.Len(actualConds, 3)
	// expectedArgs := "[{3 {2 true} {0 false} Left {0 false} {1 true} { false} {0 false}} {4 {2 true} {0 false} Right {0 false} {2 true} { false} {0 false}}]"
	// s.Equal(expectedArgs, fmt.Sprintf("%v", actualArgs))
	// expectedConds := "[{1 1 {0 false} {1 true} {0 false}} {2 1 {2 true} {0 false} {1 true}} {5 1 {3 true} {0 false} {1 true}}]"
	// s.Equal(expectedConds, fmt.Sprintf("%v", actualConds))
	// dump(actualArgs)
	// dump(actualConds)
}

func (s *ServiceSuite) Test53() {
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
	expected := "[{1 1 {0 false} {1 true} {0 false}} {2 1 {2 true} {0 false} {1 true}} {3 1 {0 false} {0 true} {1 true}} {4 1 {4 true} {0 false} {3 true}} {5 1 {3 true} {0 false} {3 true}}]"
	s.Equal(expected, fmt.Sprintf("%v", actual))
	// dump(actual)
}

func (s *ServiceSuite) Test60() {
	defer func() { _ = recover() }()
	ToDb(engine.Rule{}, nil)
	s.Fail("expected to panic")
	/* _, _, _, _, err :=  */
	// s.Nil(err)
}

// func (s *MappingsSuite) Test61() {
// 	ToDb(engine.Rule{}, utils.NewSeq())
// 	/* _, _, _, _, err :=  */
// 	// s.Nil(err)
// }

func (s *ServiceSuite) Test62() {
	inrule := engine.Rule{}
	outrule, outconds, outactions, outargs, mappings := ToDb(inrule, utils.NewSeq())
	expected := "{1  {0 true} {0 true}}"
	s.Equal(expected, fmt.Sprintf("%v", outrule))
	s.Len(outconds, 0)
	s.Len(outactions, 0)
	s.Len(outargs, 0)
	s.Len(mappings, 0)
	// dump("outrule", outrule)
	// dump("outconds", outconds)
	// dump("outargs", outargs)
	// dump("err", err)
	// s.Nil(err)
}

func (s *ServiceSuite) Test63() {
	inrule := engine.Rule{
		Comments: "unit test",
		Disabled: true,
		Condition: engine.Condition{
			Fn: engine.COND_EQUAL, Args: engine.Args{"One": 1, "Two": 2, "Three": []any{3, 4}},
		},
		Actions: []engine.Action{
			{
				Fn:       engine.ACTION_ZIGBEE2_MQTT_SET_STATE,
				DeviceId: "0x00158d0004244bda",
				Args:     engine.Args{"Lorem": 100, "Ipsum": "200"},
				Mapping: engine.Mapping{
					"Lorem": {"Ipsum": "112233", "Bar": "Baz"},
				},
			},
		},
		Throttle: 100500,
	}
	outrule, outconds, outactions, outargs, mappings := ToDb(inrule, utils.NewSeq())

	expectedRule := "{1 unit test {1 true} {0 true}}"
	s.Equal(expectedRule, fmt.Sprintf("%v", outrule))

	expectedConds := "[{2 1 {2 true} {0 false} {0 false}}]"
	s.Len(outconds, 1)
	s.Equal(expectedConds, fmt.Sprintf("%v", outconds))

	expectedActs := "[{7 1 {5 true} {0x00158d0004244bda true}}]"
	s.Len(outactions, 1)
	s.Equal(expectedActs, fmt.Sprintf("%v", outactions))

	s.Len(outargs, 6)
	s.Len(mappings, 2)

	// 1 we have to use here, since json.Marshal(...) gives sorted keys unlike fmt.Sprintf("%v"...)
	// 2 we have to use here, slices.SortFunc
	// expectedMappings := "[{8 6 Ipsum 112233} {9 6 Bar Baz}]"
	// s.Equal(expectedMappings, fmt.Sprintf("%v", mappings))
	// fmt.Println(outargs)
	// dump("outrule", outrule)
	// dump("outconds", outconds)
	// dump("outargs", outargs)
	// dump("err", err)
	// expectedArgs := `[{"Id":3,"ConditionId":{"Int32":2,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"One","IsList":{"Int32":0,"Valid":false},"Value":{"String":"1","Valid":true},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":0,"Valid":false}},{"Id":4,"ConditionId":{"Int32":2,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"Two","IsList":{"Int32":0,"Valid":false},"Value":{"String":"2","Valid":true},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":0,"Valid":false}},{"Id":6,"ConditionId":{"Int32":0,"Valid":false},"ActionId":{"Int32":5,"Valid":true},"ArgumentName":"Lorem","IsList":{"Int32":0,"Valid":false},"Value":{"String":"100","Valid":true},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":0,"Valid":false}},{"Id":7,"ConditionId":{"Int32":0,"Valid":false},"ActionId":{"Int32":5,"Valid":true},"ArgumentName":"Ipsum","IsList":{"Int32":0,"Valid":false},"Value":{"String":"200","Valid":true},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":0,"Valid":false}}]`
	// outargsj, _ := json.Marshal(outargs)
	// s.Equal(expectedArgs, string(outargsj))
	// expectedArgs := "[{3 {2 true} {0 false} Two {0 false} {2 true} { false} {0 false}} {4 {2 true} {0 false} One {0 false} {1 true} { false} {0 false}}]"
	// s.Equal(expectedArgs, fmt.Sprintf("%v", outargs))
	// s.Len(mappings, 0)
	// s.Nil(err)

}

type mockrepo struct {
	err error
}

func (r mockrepo) Get() (
	rules []DbRule,
	conditions []DbRuleCondition,
	ruleActions []DbRuleAction,
	args []DbRuleConditionOrActionArgument,
	mappings []DbRuleActionArgumentMapping,
	err error,
) {
	err = r.err
	return
}
func (r mockrepo) Create(
	rule DbRule,
	conditions []DbRuleCondition,
	actions []DbRuleAction,
	arguments []DbRuleConditionOrActionArgument,
	mappings []DbRuleActionArgumentMapping,
) (err error) {
	err = r.err
	return
}

func (s *ServiceSuite) Test70() {
	repo := mockrepo{}
	service := NewService(repo)
	rules, err := service.Get()
	s.Len(rules, 0)
	s.Nil(err)
}

func (s *ServiceSuite) Test71() {
	repo := mockrepo{}
	service := NewService(repo)
	err := service.Create(engine.Rule{})
	s.Nil(err)
}

func (s *ServiceSuite) Test72() {
	repo := mockrepo{errors.New("mock error")}
	service := NewService(repo)
	rr, err := service.Get()
	s.Len(rr, 0)
	s.NotNil(err)
	err = service.Create(engine.Rule{})
	s.NotNil(err)
}

func (s *ServiceSuite) Test80() {
	aa := ToDbArguments(
		&DbRuleCondition{},
		nil,
		"key",
		111,
		utils.NewSeq(),
		false,
	)
	data, _ := json.Marshal(aa)
	s.Equal(`[{"Id":1,"ConditionId":{"Int32":0,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"key","IsList":{"Int32":0,"Valid":true},"Value":{"String":"111","Valid":true},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":0,"Valid":false}}]`, string(data))
}

func (s *ServiceSuite) Test81() {
	aa := ToDbArguments(
		&DbRuleCondition{},
		nil,
		"key2",
		[]any{222, 333},
		utils.NewSeq(),
		false,
	)
	data, _ := json.Marshal(aa)
	expected := `[{"Id":1,"ConditionId":{"Int32":0,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"key2","IsList":{"Int32":1,"Valid":true},"Value":{"String":"222","Valid":true},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":0,"Valid":false}},{"Id":2,"ConditionId":{"Int32":0,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"key2","IsList":{"Int32":1,"Valid":true},"Value":{"String":"333","Valid":true},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":0,"Valid":false}}]`
	s.Equal(expected, string(data))
}

func (s *ServiceSuite) Test82() {
	aa := ToDbArguments(
		&DbRuleCondition{},
		nil,
		"key3",
		engine.DeviceId("0xqwe111111"),
		utils.NewSeq(),
		false,
	)
	data, _ := json.Marshal(aa)
	expected := `[{"Id":1,"ConditionId":{"Int32":0,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"key3","IsList":{"Int32":0,"Valid":true},"Value":{"String":"","Valid":false},"DeviceId":{"String":"0xqwe111111","Valid":true},"DeviceClassId":{"Int32":0,"Valid":false}}]`
	s.Equal(expected, string(data))
}

func (s *ServiceSuite) Test83() {
	aa := ToDbArguments(
		&DbRuleCondition{},
		nil,
		"key4",
		engine.DEVICE_CLASS_ZIGBEE_DEVICE,
		utils.NewSeq(),
		false,
	)
	data, _ := json.Marshal(aa)
	expected := `[{"Id":1,"ConditionId":{"Int32":0,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"key4","IsList":{"Int32":0,"Valid":true},"Value":{"String":"","Valid":false},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":1,"Valid":true}}]`
	s.Equal(expected, string(data))
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
