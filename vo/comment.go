package vo

type Comment struct {
	Code int `json:"code"`
	Data struct {
		Page struct {
			Num   int `json:"num"`
			Size  int `json:"size"`
			Count int `json:"count"`
		} `json:"page"`
		Replies []struct {
			Rpid    int64 `json:"rpid"`
			Oid     int   `json:"oid"`
			Mid     int   `json:"mid"`
			Like    int   `json:"like"`
			Ctime   int   `json:"ctime"`
			Content struct {
				Message string `json:"message"`
			} `json:"content,omitempty"`
		} `json:"replies"`
	} `json:"data"`
}
