package dto

type CreateMomentReq struct {
	Content string   `json:"content"`
	Images  []string `json:"images"`
	Video   string   `json:"video"`
}
