package articles

type CreateOrUpdateArticlesRequest struct {
	Inventory []Article `json:"inventory"`
}

type Article struct {
	ArticleID string `json:"art_id"` // nolint
	Name      string `json:"name"`
	Stock     string `json:"stock"`
}
