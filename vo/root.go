package vo

type MidAndTask struct {
	Mid     int  `json:"mid"`
	Task    Task `json:"task"`
	hasNext int  //-1就是插入发生重复
}

type Paged interface {
	HasNextPage() bool
	Next() []byte
	Store()
}
