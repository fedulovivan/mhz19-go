package rules

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/types"
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
				Name:       "case 2",
				IsDisabled: db.NewNullInt32(0),
				ThrottleMs: db.NewNullInt32(0),
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
				Id:            1,
				ConditionId:   db.NewNullInt32(1),
				ArgumentName:  "Left",
				Value:         db.NewNullString("foo"),
				ValueDataType: db.NewNullString("string"),
			},
			{
				Id:            2,
				ConditionId:   db.NewNullInt32(1),
				ArgumentName:  "Right",
				Value:         db.NewNullString("bar"),
				ValueDataType: db.NewNullString("string"),
			},
			{
				Id:            3,
				ConditionId:   db.NewNullInt32(1),
				ArgumentName:  "Third",
				Value:         db.NewNullString("baz"),
				ValueDataType: db.NewNullString("string"),
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
				Id:            6,
				ActionId:      db.NewNullInt32(1),
				ArgumentName:  "Value",
				Value:         db.NewNullString("$message.action"),
				ValueDataType: db.NewNullString("string"),
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
		expectedJson: `[{"id":1,"name":"case 2","condition":{"fn":"Changed","args":{"Fifth":"DeviceId(0x00158d0004244bda)","Fourth":"zigbee-device","Left":"foo","Right":"bar","Third":"baz"}},"actions":[{"fn":"PostSonoffSwitchMessage","args":{"ListIds":["DeviceId(10011cec96)","DeviceId(78345aaa67)"],"Value":"$message.action"},"mapping":{"Value":{"lorem-3":"dolor-4","sit-5":"amet-6"}}}],"throttle":null}]`,
	},
	// case 02
	{
		rules: []DbRule{
			{
				Id:         2,
				Name:       "case 3",
				IsDisabled: db.NewNullInt32(0),
				ThrottleMs: db.NewNullInt32(0),
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
		expectedJson: `[{"id":2,"name":"case 3","condition":{"nested":[{"fn":"Changed"},{"fn":"Equal"},{"nested":[{"fn":"InList"},{"fn":""}],"or":true}]},"throttle":null}]`,
	},
}

type mockrepo struct {
	err error
}

func (r mockrepo) Get(ruleId sql.NullInt32) (
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
) (ruleId int64, err error) {
	err = r.err
	return
}
func (r mockrepo) Delete(int32) error {
	return nil
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
	// fmt.Println(string(data))
}

func (s *ServiceSuite) Test20() {
	s.PanicsWithValue("unexpected conditions", func() {
		BuildArguments([]DbRuleConditionOrActionArgument{
			{},
		})
	})
}

func (s *ServiceSuite) Test21() {
	s.Nil(BuildArguments([]DbRuleConditionOrActionArgument{}))
}

func (s *ServiceSuite) Test22() {
	actual := BuildArguments([]DbRuleConditionOrActionArgument{
		{
			ArgumentName:  "Bar",
			Value:         db.NewNullString("foo"),
			ValueDataType: db.NewNullString("string"),
		},
	})
	s.Equal("foo", actual["Bar"])
}

func (s *ServiceSuite) Test23() {
	actual := BuildArguments([]DbRuleConditionOrActionArgument{
		{
			ArgumentName:  "Bar1",
			Value:         db.NewNullString("123"),
			ValueDataType: db.NewNullString("int"),
		},
	})
	s.Equal(int(123), actual["Bar1"])
}

func (s *ServiceSuite) Test24() {
	actual := BuildArguments([]DbRuleConditionOrActionArgument{
		{
			ArgumentName:  "Bar2",
			Value:         db.NewNullString("111"),
			ValueDataType: db.NewNullString("float64"),
		},
	})
	s.Equal(float64(111), actual["Bar2"])
}

func (s *ServiceSuite) Test25() {
	actual := BuildArguments([]DbRuleConditionOrActionArgument{
		{
			ArgumentName:  "Bar3",
			Value:         db.NewNullString("true"),
			ValueDataType: db.NewNullString("bool"),
		},
	})
	s.Equal(true, actual["Bar3"])
}

func (s *ServiceSuite) Test26() {
	s.PanicsWithValue("unexpected value data type types.Rule", func() {
		BuildArguments([]DbRuleConditionOrActionArgument{
			{
				ArgumentName:  "Bar4",
				Value:         db.NewNullString("{}"),
				ValueDataType: db.NewNullString("types.Rule"),
			},
		})
	})
}

func (s *ServiceSuite) Test27() {
	actual := BuildArguments([]DbRuleConditionOrActionArgument{
		{
			ArgumentName: "Bar5",
			DeviceId:     db.NewNullString("lorem111"),
		},
	})
	s.Equal(types.DeviceId("lorem111"), actual["Bar5"])
}

func (s *ServiceSuite) Test28() {
	actual := BuildArguments([]DbRuleConditionOrActionArgument{
		{
			ArgumentName:  "Bar6",
			DeviceClassId: db.NewNullInt32(2),
		},
	})
	s.Equal(types.DEVICE_CLASS_PINGER, actual["Bar6"])
}

func (s *ServiceSuite) Test29() {
	actual := BuildArguments([]DbRuleConditionOrActionArgument{
		{
			ArgumentName:  "Bar7",
			ChannelTypeId: db.NewNullInt32(3),
		},
	})
	s.Equal(types.CHANNEL_DNS_SD, actual["Bar7"])
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
	ToDbConditions(1, nil, types.Condition{}, nil, nil)
}

func (s *ServiceSuite) Test51() {
	s.PanicsWithValue("unexpected conditions", func() {
		ToDbConditions(1, nil, types.Condition{
			Fn:     types.COND_CHANGED,
			Nested: []types.Condition{{Fn: types.COND_EQUAL}},
		}, nil, nil)
	})
}

func (s *ServiceSuite) Test52() {
	seq := &atomic.Int32{}
	actualArgs := []DbRuleConditionOrActionArgument{}
	actualConds := ToDbConditions(1, nil, types.Condition{
		Or: true,
		Nested: []types.Condition{
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1, "Right": 2}},
			{Fn: types.COND_IN_LIST},
		},
	}, seq, &actualArgs)
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
	seq := &atomic.Int32{}
	actual := ToDbConditions(1, nil, types.Condition{
		Or: true,
		Nested: []types.Condition{
			{Fn: types.COND_EQUAL},
			{Nested: []types.Condition{
				{Fn: types.COND_EQUAL, Not: true},
				{Fn: types.COND_IN_LIST},
			}},
		},
	}, seq, nil)
	expected := "[{1 1 {0 false} {1 true} {0 false} {0 false} { false}} {2 1 {2 true} {0 false} {0 true} {1 true} { false}} {3 1 {0 false} {0 true} {0 false} {1 true} { false}} {4 1 {2 true} {0 false} {1 true} {3 true} { false}} {5 1 {3 true} {0 false} {0 true} {3 true} { false}}]"
	s.Equal(expected, fmt.Sprintf("%v", actual))
}

func (s *ServiceSuite) Test60() {
	s.Panics(func() {
		ToDb(types.Rule{}, nil)
	})
}

func (s *ServiceSuite) Test62() {
	seq := &atomic.Int32{}
	inrule := types.Rule{}
	outrule, outconds, outactions, outargs, mappings := ToDb(inrule, seq)
	expected := "{1  {0 true} {0 true}}"
	s.Equal(expected, fmt.Sprintf("%v", outrule))
	s.Len(outconds, 0)
	s.Len(outactions, 0)
	s.Len(outargs, 0)
	s.Len(mappings, 0)
}

func (s *ServiceSuite) Test63() {
	seq := &atomic.Int32{}
	inrule := types.Rule{
		Name:     "unit test",
		Disabled: true,
		Condition: types.Condition{
			Fn: types.COND_EQUAL, Args: types.Args{"One": 1, "Two": 2, "Three": []any{3, 4}},
		},
		Actions: []types.Action{
			{
				Fn:   types.ACTION_ZIGBEE2_MQTT_SET_STATE,
				Args: types.Args{"Lorem": 100, "Ipsum": "200"},
				Mapping: types.Mapping{
					"Lorem": {"Ipsum": "112233", "Bar": "Baz"},
				},
			},
		},
		Throttle: types.Throttle{
			Duration: time.Duration(100500),
		},
	}
	outrule, outconds, outactions, outargs, mappings := ToDb(inrule, seq)

	expectedRule := "{1 unit test {1 true} {0 true}}"
	s.Equal(expectedRule, fmt.Sprintf("%v", outrule))

	expectedConds := "[{2 1 {2 true} {0 false} {0 true} {0 false} { false}}]"
	s.Len(outconds, 1)
	s.Equal(expectedConds, fmt.Sprintf("%v", outconds))

	expectedActs := "[{7 1 {5 true}}]"
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
	go func() {
		for range service.OnCreated() {
			// noop, just to unblock sender
		}
	}()
	_, err := service.Create(types.Rule{})
	s.Nil(err)
}

func (s *ServiceSuite) Test72() {
	repo := mockrepo{}
	service := NewService(repo)
	go func() {
		for range service.OnDeleted() {
			// noop, just to unblock sender
		}
	}()
	err := service.Delete(111)
	s.Nil(err)
}

func (s *ServiceSuite) Test73() {
	repo := mockrepo{errors.New("mock error")}
	service := NewService(repo)
	rr, err := service.Get()
	s.Len(rr, 0)
	s.NotNil(err)
	_, err = service.Create(types.Rule{})
	s.NotNil(err)
}

func (s *ServiceSuite) Test80() {
	seq := &atomic.Int32{}
	aa := ToDbArguments(
		1,
		&DbRuleCondition{},
		nil,
		"key",
		111,
		seq,
		false,
	)
	data, _ := json.Marshal(aa)
	// fmt.Println(string(data))
	expected := `[{"Id":1,"RuleId":1,"ConditionId":{"Int32":0,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"key","IsList":{"Int32":0,"Valid":true},"Value":{"String":"111","Valid":true},"ValueDataType":{"String":"int","Valid":true},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":0,"Valid":false},"ChannelTypeId":{"Int32":0,"Valid":false}}]`
	s.Equal(expected, string(data))
}

func (s *ServiceSuite) Test81() {
	seq := &atomic.Int32{}
	aa := ToDbArguments(
		1,
		&DbRuleCondition{},
		nil,
		"key2",
		[]any{222, 333},
		seq,
		false,
	)
	data, _ := json.Marshal(aa)
	// fmt.Println(string(data))
	expected := `[{"Id":1,"RuleId":1,"ConditionId":{"Int32":0,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"key2","IsList":{"Int32":1,"Valid":true},"Value":{"String":"222","Valid":true},"ValueDataType":{"String":"int","Valid":true},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":0,"Valid":false},"ChannelTypeId":{"Int32":0,"Valid":false}},{"Id":2,"RuleId":1,"ConditionId":{"Int32":0,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"key2","IsList":{"Int32":1,"Valid":true},"Value":{"String":"333","Valid":true},"ValueDataType":{"String":"int","Valid":true},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":0,"Valid":false},"ChannelTypeId":{"Int32":0,"Valid":false}}]`
	s.Equal(expected, string(data))
}

func (s *ServiceSuite) Test82() {
	seq := &atomic.Int32{}
	aa := ToDbArguments(
		1,
		&DbRuleCondition{},
		nil,
		"key3",
		types.DeviceId("0xqwe111111"),
		seq,
		false,
	)
	data, _ := json.Marshal(aa)
	// fmt.Println(string(data))
	expected := `[{"Id":1,"RuleId":1,"ConditionId":{"Int32":0,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"key3","IsList":{"Int32":0,"Valid":true},"Value":{"String":"","Valid":false},"ValueDataType":{"String":"","Valid":false},"DeviceId":{"String":"0xqwe111111","Valid":true},"DeviceClassId":{"Int32":0,"Valid":false},"ChannelTypeId":{"Int32":0,"Valid":false}}]`
	s.Equal(expected, string(data))
}

func (s *ServiceSuite) Test83() {
	seq := &atomic.Int32{}
	aa := ToDbArguments(
		1,
		&DbRuleCondition{},
		nil,
		"key4",
		types.DEVICE_CLASS_ZIGBEE_DEVICE,
		seq,
		false,
	)
	data, _ := json.Marshal(aa)
	// fmt.Println(string(data))
	expected := `[{"Id":1,"RuleId":1,"ConditionId":{"Int32":0,"Valid":true},"ActionId":{"Int32":0,"Valid":false},"ArgumentName":"key4","IsList":{"Int32":0,"Valid":true},"Value":{"String":"","Valid":false},"ValueDataType":{"String":"","Valid":false},"DeviceId":{"String":"","Valid":false},"DeviceClassId":{"Int32":1,"Valid":true},"ChannelTypeId":{"Int32":0,"Valid":false}}]`
	s.Equal(expected, string(data))
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
