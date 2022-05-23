package vo

type MidAndTask struct {
	// 目标id
	Mid  int  `json:"mid"`
	Task Task `json:"task"`
	//-1的意思是插入发生重复
	HasNext int `json:"has_next"`
}

type Paged interface {
	HasNextPage() bool
	Next() []byte
	Store()
}
