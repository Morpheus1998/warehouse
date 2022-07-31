package products

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/warehouse/app/server/responses"
	"github.com/warehouse/app/store"
)

type Handler struct {
	ProductsStore store.ProductsStore
}

func NewHandler() *Handler {
	return &Handler{}
}

// CreateOrUpdateProducts is http api POST /products
func (h *Handler) CreateOrUpdateProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &CreateOrUpdateProductsRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error().AnErr("error", err).Msg("CreateOrUpdateProducts failed to unmarshal request")
		body := responses.GenerateErrorResponseBody(ctx, responses.UnMarshalRequestError, err.Error())
		responses.WriteError(ctx, w, http.StatusBadRequest, body)
		return
	}
	dbReq, err := getCreateOrUpdateProductsDBRequest(req)
	if err != nil {
		log.Error().AnErr("error", err).Msg("CreateOrUpdateProducts get database request from http request")
		body := responses.GenerateErrorResponseBody(ctx, responses.InvalidBodyError, err.Error())
		responses.WriteError(ctx, w, http.StatusBadRequest, body)
		return
	}
	err = h.ProductsStore.CreateOrUpdateProducts(ctx, dbReq)
	if err != nil {
		if errors.Is(err, store.ErrArticleNotFound) {
			log.Error().AnErr("error", err).Msg("CreateOrUpdateProducts failed to execute database query, article not found")
			body := responses.GenerateErrorResponseBody(ctx, responses.ResourceNotFound, err.Error())
			responses.WriteError(ctx, w, http.StatusNotFound, body)
			return
		}
		log.Error().AnErr("error", err).Msg("CreateOrUpdateProducts failed to execute database query")
		body := responses.GenerateErrorResponseBody(ctx, responses.DataBaseQueryFailureError, err.Error())
		responses.WriteError(ctx, w, http.StatusInternalServerError, body)
		return
	}
	responses.WriteCreatedResponse(ctx, w, nil)
}

// SellProduct is http api POST /products/sell
func (h *Handler) SellProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &SellProductRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error().AnErr("error", err).Msg("SellProduct failed to unmarshal request")
		body := responses.GenerateErrorResponseBody(ctx, responses.UnMarshalRequestError, err.Error())
		responses.WriteError(ctx, w, http.StatusBadRequest, body)
		return
	}
	err = h.ProductsStore.RemoveProductAndUpdateArticles(
		ctx,
		store.RemoveProductAndUpdateArticlesRequest{
			ProductID: req.ProductID,
		})
	if err != nil {
		if errors.Is(err, store.ErrProductNotFound) {
			log.Error().AnErr("error", err).Msg("SellProduct failed to execute database query, product not found")
			body := responses.GenerateErrorResponseBody(ctx, responses.ResourceNotFound, err.Error())
			responses.WriteError(ctx, w, http.StatusNotFound, body)
			return
		}
		if errors.Is(err, store.ErrProductStockFinished) {
			log.Error().AnErr("error", err).Msg("SellProduct failed to execute database query, product stock finished")
			body := responses.GenerateErrorResponseBody(ctx, responses.ResourceFinished, err.Error())
			responses.WriteError(ctx, w, http.StatusBadRequest, body)
			return
		}
		log.Error().AnErr("error", err).Msg("SellProduct failed to execute database query")
		body := responses.GenerateErrorResponseBody(ctx, responses.DataBaseQueryFailureError, err.Error())
		responses.WriteError(ctx, w, http.StatusInternalServerError, body)
		return
	}
	responses.WriteNoContentResponse(ctx, w)
}

// GetAllProductsWithStock is http api GET /products
func (h *Handler) GetAllProductsWithStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := h.ProductsStore.GetAllProducts(ctx)
	if err != nil {
		log.Error().AnErr("error", err).Msg("GetAllProductsWithStock failed to execute database query")
		body := responses.GenerateErrorResponseBody(ctx, responses.DataBaseQueryFailureError, err.Error())
		responses.WriteError(ctx, w, http.StatusInternalServerError, body)
		return
	}
	response := getProductsResponseFromDBResult(res)
	responses.WriteOkResponse(ctx, w, response)
}

func getProductsResponseFromDBResult(dbResult store.GetAllProductsResponse) *GetAllProductsWithStockResponse {
	response := &GetAllProductsWithStockResponse{}
	for _, product := range dbResult.Products {
		productArticles := make([]Article, 0)
		for _, productArticle := range product.Articles {
			productArticles = append(productArticles, Article{
				Amount:    strconv.Itoa(productArticle.ArticleAmount),
				ArticleID: productArticle.ArticleID,
			})
		}
		response.Products = append(response.Products, ProductWithStock{
			Stock:     product.Stock,
			ProductID: product.ProductID,
			Product: Product{
				Name:     product.ProductName,
				Articles: productArticles,
			},
		})
	}
	return response
}

func getCreateOrUpdateProductsDBRequest(req *CreateOrUpdateProductsRequest) (store.CreateOrUpdateProductsRequest, error) {
	res := store.CreateOrUpdateProductsRequest{
		Products: []store.Product{},
	}
	for _, product := range req.Products {
		productArticles := make([]store.ProductArticle, 0)
		for _, productArticle := range product.Articles {
			articleAmount, err := strconv.Atoi(productArticle.Amount)
			if err != nil {
				log.Error().AnErr("error", err).Msg("failed to parse product article amount to integer")
				return store.CreateOrUpdateProductsRequest{}, err
			}
			productArticles = append(productArticles, store.ProductArticle{
				ArticleID:     productArticle.ArticleID,
				ArticleAmount: articleAmount,
			})
		}
		storeProduct := store.Product{
			ProductName: product.Name,
			Articles:    productArticles,
		}
		res.Products = append(res.Products, storeProduct)
	}
	return res, nil
}
