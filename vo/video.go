package vo

type Video struct {
	Code int `json:"code"`
	Data struct {
		List struct {
			Vlist []struct {
				Aid int `json:"aid"`
			} `json:"vlist"`
		} `json:"list"`
		Page struct {
			Pn    int `json:"pn"`
			Ps    int `json:"ps"`
			Count int `json:"count"`
		} `json:"page"`
	} `json:"data"`
}
