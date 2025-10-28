package dto

type Lists struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
type GetListResponse struct {
	Result []Lists `json:"result"`
}
