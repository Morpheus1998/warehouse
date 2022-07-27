package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/warehouse/app/articles"
	"github.com/warehouse/app/products"
	"github.com/warehouse/app/store"
)

type Server struct {
	ProductsHandler *products.Handler
	ArticlesHandler *articles.Handler
}

func (srv *Server) setHandlers() {
	if srv.ArticlesHandler == nil {
		srv.ArticlesHandler = articles.NewHandler()
	}
	if srv.ProductsHandler == nil {
		srv.ProductsHandler = products.NewHandler()
	}
}

func (srv *Server) setStores(pgDB interface{}) error {
	var ok bool
	srv.setHandlers()
	if srv.ProductsHandler.ProductsStore, ok = pgDB.(store.ProductsStore); !ok {
		return ErrInvalidTypeForStore
	}
	if srv.ArticlesHandler.ArticleStore, ok = pgDB.(store.ArticlesStore); !ok {
		return ErrInvalidTypeForStore
	}
	return nil
}

func StartServer(cfg Configuration) error {
	ctx := context.Background()
	log.Ctx(ctx).Info().Msg("enter StartServer")
	defer log.Ctx(ctx).Info().Msg("exit StartServer")

	server := &Server{}

	if server.ProductsHandler == nil {
		db, err := store.NewPostgresDB(
			cfg.PostgresConfiguration.Port,
			cfg.PostgresConfiguration.Host,
			cfg.PostgresConfiguration.DB,
			cfg.PostgresConfiguration.CredentialsFileName,
		)
		if err != nil {
			log.Error().Msg("failed to get postgres client")
			return err
		}
		err = server.setStores(db)
		if err != nil {
			log.Error().Msg("failed to set postgres client to handlers")
			return err
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	router := NewRouter(makeRoutes(server))
	httpServer := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler:           router,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
		ReadHeaderTimeout: time.Duration(cfg.HTTP.Timeout) * time.Millisecond,
	}
	// start httpServer listening
	httpServerErr := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			httpServerErr <- err
		}
		close(httpServerErr)
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	select {
	case s := <-signalChan:
		log.Info().
			Str("signal", s.String()).
			Msg("signal received: shutting down the server")
	case err, ok := <-httpServerErr:
		if !ok {
			// http server exited without an error
		} else {
			log.Warn().
				AnErr("error", err).
				Msg("unexpected error from the http server")
		}
	}

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), serverGracefulShutdownTime)
	defer cancelShutdown()
	return httpServer.Shutdown(gracefullCtx)
}
