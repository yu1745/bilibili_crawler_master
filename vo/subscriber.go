package vo

type Subs struct {
	Code int `json:"code"`
	Data struct {
		List []struct {
			Mid int `json:"mid"`
		} `json:"list"`
		Total int `json:"total"`
	} `json:"data"`
	Meta
}

func (this *Subs) Next() []byte {
	//TODO implement me
	panic("implement me")
}

func (this *Subs) Store() {
	//TODO implement me
	panic("implement me")
}
