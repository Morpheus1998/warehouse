package store

type RemoveProductAndUpdateArticlesRequest struct {
	ProductID string
}

type GetAllProductsResponse struct {
	Products []Product
}

type Product struct {
	ProductName string
	ProductID   string
	Stock       int
	ArticlesID  []string
}

type CreateOrUpdateProductsRequest struct {
	Products []Product
}

type Article struct {
	Stock       int
	ArticleName string
	ArticleID   string
}

type CreateOrUpdateArticlesRequest struct {
	Articles []Article
}
