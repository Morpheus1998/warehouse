package articles

import (
	"net/http"

	"github.com/warehouse/app/store"
)

type Handler struct {
	ArticleStore store.ArticlesStore
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) CreateOrUpdateArticles(w http.ResponseWriter, r *http.Request) {

}
