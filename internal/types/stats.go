package types

type TableStats struct {
	Rules    int32 `json:"rules"`
	Devices  int32 `json:"devices"`
	Messages int32 `json:"messages"`
}
