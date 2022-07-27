package store

import (
	"context"
	"errors"
)

type ProductsStore interface {
	CreateOrUpdateProducts(ctx context.Context, req CreateOrUpdateProductsRequest) error
	RemoveProductAndUpdateArticles(ctx context.Context, req RemoveProductAndUpdateArticlesRequest) error
	GetAllProducts(ctx context.Context) (GetAllProductsResponse, error)
}

type ArticlesStore interface {
	CreateOrUpdateArticles(ctx context.Context, req CreateOrUpdateArticlesRequest) error
}

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrArticleNotFound      = errors.New("article not found")
	ErrProductStockFinished = errors.New("product stock has finished")
)
