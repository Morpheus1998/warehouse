package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

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
	return nil
}

func (pg *PostgresDB) GetAllProducts(ctx context.Context) (GetAllProductsResponse, error) {
	return GetAllProductsResponse{}, nil
}

func (pg *PostgresDB) CreateOrUpdateProducts(ctx context.Context, req CreateOrUpdateProductsRequest) error {
	return nil
}

func (pg *PostgresDB) CreateOrUpdateArticles(ctx context.Context, req CreateOrUpdateArticlesRequest) error {
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
