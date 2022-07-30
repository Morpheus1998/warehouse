package articles

import (
	"encoding/json"
	"net/http"

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
	dbReq := getCreateOrUpdateArticleDBRequest(req)
	err = h.ArticleStore.CreateOrUpdateArticles(ctx, dbReq)
	if err != nil {
		log.Error().AnErr("error", err).Msg("CreateOrUpdateArticles failed to execute database query")
		body := responses.GenerateErrorResponseBody(ctx, responses.DataBaseQueryFailureError, err.Error())
		responses.WriteError(ctx, w, http.StatusInternalServerError, body)
		return
	}
	responses.WriteCreatedResponse(ctx, w, nil)
}

func getCreateOrUpdateArticleDBRequest(req *CreateOrUpdateArticlesRequest) store.CreateOrUpdateArticlesRequest {
	res := store.CreateOrUpdateArticlesRequest{}
	for _, article := range req.Inventory {
		res.Articles = append(res.Articles, store.Article{
			ArticleID:   article.ArticleID,
			ArticleName: article.Name,
			Stock:       article.Stock,
		})
	}
	return res
}
