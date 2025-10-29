package dto

type Lists struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
type GetListResponse struct {
	Result []Lists `json:"result"`
}

type RegisterWebHookRequest struct {
	Destination string   `json:"destination"`
	Settings    []string `json:"settings"`
	Sort        int      `json:"sort"`
}
