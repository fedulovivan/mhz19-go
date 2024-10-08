package types

type TableStats struct {
	Rules    int32 `json:"rules"`
	Devices  int32 `json:"devices"`
	Messages int32 `json:"messages"`
	Actions  int32 `json:"actions"`
	Conds    int32 `json:"conditions"`
	Args     int32 `json:"arguments"`
	Mappings int32 `json:"mappings"`
}
