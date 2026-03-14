package dto

// ---- 请求 ----

type CreatePostReq struct {
	Title    string   `json:"title"    binding:"required"`
	Content  string   `json:"content"  binding:"required"`
	Summary  string   `json:"summary"`
	Cover    string   `json:"cover"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Status   string   `json:"status"`
}

type UpdatePostReq struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Summary  string   `json:"summary"`
	Cover    string   `json:"cover"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Status   string   `json:"status"`
}

type PostFilter struct {
	Category string
	Tag      string
}

// ---- 响应 ----

type CategoryResp struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type TagResp struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type ArchiveResp struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}
