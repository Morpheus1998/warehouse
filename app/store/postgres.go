package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Credentials struct {
	Username string `json:"USERNAME"` //nolint
	Password string `json:"PASSWORD"` //nolint
}

type PostgresDB struct {
	Database *sql.DB
}

func (pg *PostgresDB) Ping(ctx context.Context) error {
	return pg.Database.Ping()
}

func (pg *PostgresDB) Close() error {
	return pg.Database.Close()
}

func NewPostgresDB(port int, host, dbname, credentialJSONFile string) (*PostgresDB, error) {
	cred, err := credentialsFromFile(credentialJSONFile)
	if err != nil {
		log.Error().AnErr("error", err).Msgf("failed to read credentials from file: %v", credentialJSONFile)
		return nil, err
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, cred.Username, cred.Password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error().AnErr("error", err).Msgf("failed to connect to database with address: %v:%v", host, port)
		return nil, err
	}
	return &PostgresDB{
		Database: db,
	}, nil
}

func (pg *PostgresDB) RemoveProductAndUpdateArticles(
	ctx context.Context,
	req RemoveProductAndUpdateArticlesRequest,
) error {
	productArticles, err := pg.getProductArticlesByProductID(ctx, req.ProductID)
	if err != nil {
		log.Ctx(ctx).Error().AnErr("error", err).Msg("sell product, failed to get product_article by product id")
		return err
	}
	tx, err := pg.Database.Begin()
	if err != nil {
		log.Ctx(ctx).Error().AnErr("error", err).Msg("sell product, failed to start transaction")
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Ctx(ctx).Err(rollbackErr).Msg("error happened when rolling back tx in RemoveProductAndUpdateArticles")
			}
		} else {
			err = tx.Commit()
		}
	}()
	// TODO: enhance query so that all article updates happen in a single query
	for _, productArticle := range productArticles {
		_, err = tx.Exec(updateArticleStockForSellProduct, productArticle.ArticleAmount, productArticle.ArticleID)
		if err != nil {
			log.Ctx(ctx).Error().AnErr("error", err).Msg("sell product, failed to update article")
			return err
		}
	}

	return nil
}

func (pg *PostgresDB) getArticlesByProductID(ctx context.Context, productID string) ([]Article, error) {
	rows, err := pg.Database.Query(getArticlesByProductID, productID)
	if err != nil {
		log.Ctx(ctx).Error().AnErr("error", err).Msg("failed to get articles by product_id")
		return nil, err
	}
	articles := make([]Article, 0)
	for rows.Next() {
		var article Article
		err = rows.Scan(&article.ArticleID, &article.Stock)
		if err != nil {
			log.Ctx(ctx).Error().AnErr("error", err).Msg("failed to scan article by product_id")
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (pg *PostgresDB) getProductArticlesByProductID(ctx context.Context, productID string) ([]ProductArticle, error) {
	rows, err := pg.Database.Query(getProductArticlesByProductID, productID)
	if err != nil {
		log.Ctx(ctx).Error().AnErr("error", err).Msg("failed to get product_article by product_id")
		return nil, err
	}
	productArticles := make([]ProductArticle, 0)
	for rows.Next() {
		var productArticle ProductArticle
		err = rows.Scan(&productArticle.ArticleID, &productArticle.ArticleAmount)
		if err != nil {
			log.Ctx(ctx).Error().AnErr("error", err).Msg("failed to scan product_article by product_id")
			return nil, err
		}
		productArticles = append(productArticles, productArticle)
	}
	return productArticles, nil
}

func (pg *PostgresDB) GetAllProducts(ctx context.Context) (GetAllProductsResponse, error) {
	rows, err := pg.Database.Query(getProductsWithStock)
	if err != nil {
		log.Ctx(ctx).Error().AnErr("error", err).Msg("failed to get all products")
		return GetAllProductsResponse{}, err
	}
	products := make([]Product, 0)
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ProductID, &product.Stock)
		if err != nil {
			log.Ctx(ctx).Error().AnErr("error", err).Msg("failed to scan all products")
			return GetAllProductsResponse{}, err
		}
	}
	return GetAllProductsResponse{
		Products: products,
	}, nil
}

func (pg *PostgresDB) CreateOrUpdateProducts(ctx context.Context, req CreateOrUpdateProductsRequest) error {
	tx, err := pg.Database.Begin()
	if err != nil {
		log.Ctx(ctx).Error().AnErr("error", err).Msg("create or update products, failed to start transaction")
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Ctx(ctx).Err(rollbackErr).Msg("error happened when rolling back tx in CreateOrUpdateProducts")
			}
		} else {
			err = tx.Commit()
		}
	}()
	for _, product := range req.Products {
		var productID string
		err := tx.QueryRow(createProduct, product.ProductName).Scan(&productID)
		if err != nil {
			log.Ctx(ctx).Error().AnErr("error", err).Msg("failed to create product")
			return err
		}
		for _, article := range product.Articles {
			_, err = tx.Exec(createProductArticle, productID, article.ArticleID, article.ArticleAmount)
			if err != nil {
				log.Ctx(ctx).Error().AnErr("error", err).Msg("failed to create product_article")
				return err
			}
		}
	}
	return nil
}

func (pg *PostgresDB) CreateOrUpdateArticles(ctx context.Context, req CreateOrUpdateArticlesRequest) error {
	for _, article := range req.Articles {
		_, err := pg.Database.Exec(createArticle, article.ArticleID, article.Stock, article.ArticleName)
		if err != nil {
			log.Ctx(ctx).Error().AnErr("error", err).Msgf("failed to create article with id %v", article.ArticleID)
			return err
		}
	}
	return nil
}

func credentialsFromFile(filename string) (*Credentials, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	var cr Credentials
	err = json.NewDecoder(f).Decode(&cr)
	if err != nil {
		return nil, fmt.Errorf("decoding json: %w", err)
	}

	return &cr, nil
}
