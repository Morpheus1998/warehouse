package articles

type CreateOrUpdateArticlesRequest struct {
	Inventory []Article `json:"inventory"`
}

type Article struct {
	ArticleID string `json:"art_id"`
	Name      string `json:"name"`
	Stock     int    `json:"stock"`
}
