package engine

type CondFnName string

type ConditionFunctionImplementation func(mt MessageTuple, args NamedArgs) bool

const (
	Equal    CondFnName = "Equal"
	NotEqual CondFnName = "NotEqual"
	InList   CondFnName = "InList"
	Changed  CondFnName = "Changed"
	NotNil   CondFnName = "NotNil"
)

func EqualFn(mt MessageTuple, args NamedArgs) bool {
	left, e1 := arg(mt[0], args, "Left")
	right, e2 := arg(mt[0], args, "Right")
	if e1 != nil || e2 != nil {
		return false
	}
	return left == right
}

func NotEqualFn(mt MessageTuple, args NamedArgs) bool {
	left, e1 := arg(mt[0], args, "Left")
	right, e2 := arg(mt[0], args, "Right")
	if e1 != nil || e2 != nil {
		return false
	}
	return left != right
}

func InListFn(mt MessageTuple, args NamedArgs) bool {
	v, e1 := arg(mt[0], args, "Value")
	list, e2 := arg(mt[0], args, "List")
	if e1 != nil || e2 != nil {
		return false
	}
	lslice, ok := list.([]any)
	if !ok {
		return false
	}
	for _, el := range lslice {
		if el == v {
			return true
		}
	}
	return false
}

// return false for nil and empty strings
// return true for the rest
func NotNilFn(mt MessageTuple, args NamedArgs) bool {
	v, e1 := arg(mt[0], args, "Value")
	if e1 != nil {
		return false
	}
	switch vTyped := v.(type) {
	case string:
		return len(vTyped) > 0
	case nil:
		return false
	default:
		return true
	}
}

func ChangedFn(mt MessageTuple, args NamedArgs) bool {
	cv, e1 := arg(mt[0], args, "Value")
	pv, e2 := arg(mt[1], args, "Value")
	if e1 != nil || e2 != nil {
		return false
	}
	return pv != cv
}

var conditions = map[CondFnName]ConditionFunctionImplementation{
	"Equal":    EqualFn,
	"NotEqual": NotEqualFn,
	"InList":   InListFn,
	"NotNil":   NotNilFn,
	"Changed":  ChangedFn,
}
