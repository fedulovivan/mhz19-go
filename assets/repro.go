// repoduction of https://github.com/goccy/go-json/issues/526
package main

import (
	"github.com/goccy/go-json"
)

type CondFn byte
type Args map[string]any
type Rule struct {
	Id        int       `json:"id"`
	Disabled  bool      `json:"disabled,omitempty"`
	Name      string    `json:"name,omitempty"`
	Condition Condition `json:"condition,omitempty"`
}
type Condition struct {
	Id     int         `json:"-"`
	Fn     CondFn      `json:"fn,omitempty"`
	Args   Args        `json:"args,omitempty"`
	Nested []Condition `json:"nested,omitempty"`
}

func main() {
	rule := Rule{
		Condition: Condition{
			Nested: []Condition{
				{
					Args: Args{
						"Value": 111,
					},
				},
			},
		},
	}
	_, _ = json.Marshal(rule)
}
