package request

type Request struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Src     string  `json:"src"`
	StoreID string  `json:"store_id"`
	Alt     *string `json:"alt"`
}
