package products

type CreateOrUpdateProductsRequest struct {
	Products []Product `json:"products"`
}

type Product struct {
	Name     string    `json:"name"`
	Articles []Article `json:"contain_articles"` // nolint
}

type Article struct {
	ArticleID string `json:"art_id"`    // nolint
	Amount    string `json:"amount_of"` // nolint
}

type SellProductRequest struct {
	ProductID string `json:"productId"`
}

type GetAllProductsWithStockResponse struct {
	Products []ProductWithStock `json:"products"`
}

type ProductWithStock struct {
	Product
	Stock     int    `json:"stock"`
	ProductID string `json:"productId"`
}
