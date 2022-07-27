package store

import (
	"context"
	"errors"
)

type ProductsStore interface {
	RemoveProductAndUpdateArticles(ctx context.Context, req RemoveProductAndUpdateArticlesRequest) error
	GetAllProducts(ctx context.Context) (GetAllProductsResponse, error)
	CreateOrUpdateProducts(ctx context.Context, req CreateOrUpdateProductsRequest) error
}

type ArticlesStore interface {
	CreateOrUpdateArticles(ctx context.Context, req CreateOrUpdateArticlesRequest) error
}

var ErrNotFoundProduct = errors.New("product not found")
