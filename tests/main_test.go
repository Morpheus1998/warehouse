package tests

import (
	"context"
	"github.com/rs/zerolog"
	log2 "log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/warehouse/app/server"
)

var (
	httpClient *http.Client
)

const (
	serverAddressAndPort          = "127.0.0.1:8080"
	integrationTestURL            = "http://" + serverAddressAndPort
	urlHealth                     = "/health"
	urlReady                      = "/readiness"
	secondsToWaitForServerStartup = 10
)

func SetupHTTPServer() {
	log.Warn().Msg("Starting....")

	go func() {
		cfg := getTestConfig()
		err := server.StartServer(cfg)
		if err != nil {
			log.Warn().
				AnErr("error", err).
				Msg("stopped google-proxy service")
			log2.Fatal("Error setting up HTTPServer: " + err.Error())
		}
	}()
}

func SetupHTTPClient() {
	httpClient = &http.Client{
		Timeout: 15 * time.Second,
	}
}

func CheckHTTPServerHealthyAndReady() {
	CheckHTTPServerEndpoint(urlHealth)
	CheckHTTPServerEndpoint(urlReady)
}

func CheckHTTPServerEndpoint(endpointURL string) {
	// wait for server startup
	var err error
	sec := 1
	for {
		req, errReq := http.NewRequestWithContext(context.Background(), http.MethodGet, integrationTestURL+endpointURL, nil)
		if errReq != nil {
			log2.Fatal("Could not create " + endpointURL + " request " + err.Error())
		}
		var res *http.Response
		res, err = httpClient.Do(req)
		sec++

		if err == nil || sec >= secondsToWaitForServerStartup {
			if err == nil {
				res.Body.Close()
			}
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		log2.Fatal(endpointURL + " not available! Error: " + err.Error())
	}
}

func TestMain(m *testing.M) {
	// Make zerolog print humand friendly logs
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: true,
	})

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal().Msgf("couldn't get local dir: %v", err)
	}
	parent := filepath.Dir(pwd)
	StartDB(parent)
	SetupHTTPServer()
	SetupHTTPClient()
	CheckHTTPServerHealthyAndReady()
	// os.Exit() does not respect defer statements
	code := m.Run()
	os.Exit(code)
}

func getTestConfig() server.Configuration {
	cfg, err := server.GetConfigurationFromEnv()
	if err != nil {
		log.Error().AnErr("error", err).Msg("couldn't read the configuration for the test")
		log2.Fatal("Error setting up HTTPServer: " + err.Error())
	}
	cfg.PostgresConfiguration.CredentialsFileName = "../creds.json"
	return cfg
}
