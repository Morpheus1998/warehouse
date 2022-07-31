package tests

import (
	"context"
	"fmt"
	"github.com/warehouse/app/store"
	"os"
	"strconv"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog/log"
)

var (
	testDB *store.PostgresDB
)

const (
	dbUser      = "postgres"
	dbPassword  = "password"
	dbPort      = "5432"
	dbName      = "warehouse"
	credentials = "../creds.json"
)

func StartDB(liquidasePath string) (*dockertest.Resource, *dockertest.Pool, *docker.Network) {
	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost != "" {
		return nil, nil, nil
	}
	postgresPort := "5432"
	os.Setenv("POSTGRES_HOST", postgresHost)
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_DATABASE", "warehouse")
	var container *dockertest.Resource
	var pool *dockertest.Pool
	var network *docker.Network
	var err error

	os.Setenv("POSTGRES_CREDENTIALS_FILENAME", credentials)
	if postgresHost == "" {
		postgresHost = "localhost"
		containerName := "postgres-dockertest-warehouse"
		pool, err = dockertest.NewPool("")
		log.Print("Starting postgres dockertest")
		if err != nil {
			log.Fatal().Msgf("Could not connect to docker: %s", err)
		}

		network, err = pool.Client.CreateNetwork(docker.CreateNetworkOptions{Name: "database_migration_network"})
		if err != nil {
			log.Fatal().Msgf("could not create a network to liquidbase: %s", err)
		}

		container, postgresPort = initPostgres(pool, postgresHost, postgresPort, network.ID, containerName)
		if err = pool.Retry(func() error {
			if err = createDB(postgresHost, postgresPort, credentials); err != nil {
				return err
			}
			return testDB.Ping(context.Background())
		}); err != nil {
			log.Fatal().Msgf("Could not connect to database: %s", err)
		}

		log.Info().Msgf("start migrations")
		err = runMigrations(pool, containerName, network.ID, liquidasePath)
		if err != nil {
			log.Fatal().Msgf("Could not run postgres migration: %s", err)
		}
		log.Info().Msg("migrations done")
	}

	if testDB == nil {
		log.Info().Msgf("starting database with %s:%s and %s", postgresHost, postgresPort, credentials)
		err = createDB(postgresHost, postgresPort, credentials)
		if err != nil {
			log.Fatal().Msgf("Could not create postgres connection: %s", err)
		}
	}
	return container, pool, network
}

func runMigrations(pool *dockertest.Pool, host string, networkID string, liquidasePath string) error {
	options := dockertest.RunOptions{
		Repository: "liquibase/liquibase",
		Tag:        "4.3.5",
		NetworkID:  networkID,
		Cmd: []string{
			"--url",
			fmt.Sprintf("jdbc:postgresql://%s:%s/%s", host, dbPort, dbName),
			"--changeLogFile",
			"changelog/master.yml",
			"--username",
			dbUser,
			"--password",
			dbPassword,
			"--logLevel",
			"debug",
			"update",
		},
	}
	resource, err := pool.RunWithOptions(&options, func(config *docker.HostConfig) {
		config.Mounts = []docker.HostMount{
			{
				Target:   "/liquibase/changelog",
				Source:   fmt.Sprintf("%s/liquibase/changelog/", liquidasePath),
				Type:     "bind",
				ReadOnly: true,
			},
		}
	})
	if err != nil {
		return err
	}
	err = resource.Expire(60) // nolint
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second) // nolint
	return nil
}

func createDB(host string, sPort string, credFile string) error {
	port, err := strconv.Atoi(sPort)
	if err != nil {
		log.Fatal().Msgf("Could not get postgres port: %s", err)
	}
	testDB, err = store.NewPostgresDB(port, host, dbName, credFile)
	return err
}

func initPostgres(pool *dockertest.Pool, host string, port string, networkID string, containerName string) (*dockertest.Resource, string) {
	options := dockertest.RunOptions{
		Repository: "postgres",
		Name:       containerName,
		Tag:        "12.9",
		NetworkID:  networkID,
		Hostname:   host,
		Env: []string{
			"POSTGRES_USER=" + dbUser,
			"POSTGRES_PASSWORD=" + dbPassword,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{dbPort},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(port): {{HostIP: "0.0.0.0", HostPort: port}},
		},
	}
	resource, err := pool.RunWithOptions(&options)
	if err != nil {
		log.Fatal().Msgf("Could not start resource reason: %v", err)
	}
	err = resource.Expire(60) // nolint
	if err != nil {
		log.Fatal().Msgf("Could not set expire timeout to postgres: %s", err)
	}
	connectPort := resource.GetPort(port + "/tcp")
	return resource, connectPort
}
