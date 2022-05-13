package vo

type Subs struct {
	Code int `json:"code"`
	Data struct {
		List []struct {
			Mid int `json:"mid"`
		} `json:"list"`
		Total int `json:"total"`
	} `json:"data"`
	Mid int `json:"mid"`
}
