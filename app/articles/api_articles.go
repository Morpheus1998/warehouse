package articles

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/warehouse/app/server/responses"
	"github.com/warehouse/app/store"
)

type Handler struct {
	ArticleStore store.ArticlesStore
}

func NewHandler() *Handler {
	return &Handler{}
}

// CreateOrUpdateArticles is http api POST /articles
func (h *Handler) CreateOrUpdateArticles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &CreateOrUpdateArticlesRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error().AnErr("error", err).Msg("CreateOrUpdateArticles failed to unmarshal request")
		body := responses.GenerateErrorResponseBody(ctx, responses.UnMarshalRequestError, err.Error())
		responses.WriteError(ctx, w, http.StatusBadRequest, body)
		return
	}
	dbReq, err := getCreateOrUpdateArticleDBRequest(req)
	if err != nil {
		log.Error().AnErr("error", err).Msg("CreateOrUpdateArticles get database request from http request")
		body := responses.GenerateErrorResponseBody(ctx, responses.InvalidBodyError, err.Error())
		responses.WriteError(ctx, w, http.StatusBadRequest, body)
		return
	}
	err = h.ArticleStore.CreateOrUpdateArticles(ctx, dbReq)
	if err != nil {
		log.Error().AnErr("error", err).Msg("CreateOrUpdateArticles failed to execute database query")
		body := responses.GenerateErrorResponseBody(ctx, responses.DataBaseQueryFailureError, err.Error())
		responses.WriteError(ctx, w, http.StatusInternalServerError, body)
		return
	}
	responses.WriteCreatedResponse(ctx, w, nil)
}

func getCreateOrUpdateArticleDBRequest(req *CreateOrUpdateArticlesRequest) (store.CreateOrUpdateArticlesRequest, error) {
	res := store.CreateOrUpdateArticlesRequest{}
	for _, article := range req.Inventory {
		stock, err := strconv.Atoi(article.Stock)
		if err != nil {
			log.Error().AnErr("error", err).Msg("failed to parse inventory stock to integer")
			return store.CreateOrUpdateArticlesRequest{}, err
		}
		res.Articles = append(res.Articles, store.Article{
			ArticleID:   article.ArticleID,
			ArticleName: article.Name,
			Stock:       stock,
		})
	}
	return res, nil
}
