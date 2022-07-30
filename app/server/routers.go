package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

const prefix = ""

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func Ready(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func makeRoutes(srv *Server) Routes {
	generalRoutes := Routes{
		Route{
			"Health",
			http.MethodGet,
			prefix + "/health",
			Health,
		},
		Route{
			"Ready",
			http.MethodGet,
			prefix + "/readiness",
			Ready,
		},
	}
	return union(
		generalRoutes,
		getProductsRoutes(srv),
		getArticlesRoutes(srv),
	)
}

func getProductsRoutes(srv *Server) Routes {
	return Routes{
		{
			"CreateOrUpdateProducts",
			http.MethodPost,
			prefix + "/products",
			srv.ProductsHandler.CreateOrUpdateProducts,
		},
		{
			"SellProduct",
			http.MethodPost,
			prefix + "/products/sell",
			srv.ProductsHandler.SellProduct,
		},
		{
			"GetAllProductsWithStock",
			http.MethodGet,
			prefix + "/products",
			srv.ProductsHandler.GetAllProductsWithStock,
		},
	}
}

func getArticlesRoutes(srv *Server) Routes {
	return Routes{
		{
			"CreateOrUpdateArticles",
			http.MethodPost,
			prefix + "/articles",
			srv.ArticlesHandler.CreateOrUpdateArticles,
		},
	}
}

func union(routes ...Routes) Routes {
	if len(routes) == 0 {
		return Routes{}
	}
	return append(routes[0], union(routes[1:]...)...)
}
