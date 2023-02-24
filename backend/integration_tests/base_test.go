package integrationtest_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type Config struct {
	data map[string]string
}

func NewTestConfig() Config {
	envPath := os.Getenv("ENV_PATH")
	rawEnvData, err := ioutil.ReadFile(envPath)
	if err != nil {
		log.Fatalf("fail to read env file.")
	}
	var envData map[string]string
	json.Unmarshal(rawEnvData, &envData)
	config := Config{data: envData}
	return config
}

func (config Config) get(key string) string {
	content, ok := config.data[key]
	if !ok {
		log.Fatalf("fail to get %s from config", key)
	}
	return content
}

func (config Config) POSTGRES_LOCATION() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.get("POSTGRES_USER"),
		config.get("POSTGRES_PASSWORD"),
		config.get("POSTGRES_HOST"),
		config.get("POSTGRES_PORT"),
		config.get("POSTGRES_DB"),
		config.get("POSTGRES_SSL"),
	)
}

type BackendTestSuite struct {
	suite.Suite
	ServerBaseURL string
	db            *sql.DB
}

func TestBackendTestSuite(t *testing.T) {
	backendTestSuite := &BackendTestSuite{}
	backendTestSuite.ServerBaseURL = "http://127.0.0.1:8000/api/v1"
	backendTestSuite.SetT(t)
	suite.Run(t, backendTestSuite)
}

func (suite *BackendTestSuite) SetupSuite() {
	go func() {
		backendCMD := exec.Command("sh", "-c", "go run /__w/queue-system-lite/queue-system-lite/backend/main.go")
		_, err := backendCMD.Output()
		if err != nil {
			fmt.Println(err)
		}
	}()

	allDoneChan := make(chan bool)
	go func() {
		httpClient := http.Client{Timeout: 3 * time.Second}
		for {
			response, _ := httpClient.Get("http://127.0.0.1:8000/api/v1")
			if response != nil && response.StatusCode == 405 {
				break
			}
		}
		allDoneChan <- true
	}()
	<-allDoneChan

	testConfig := NewTestConfig()
	os.Setenv("POSTGRES_USER", testConfig.get("POSTGRES_USER"))
	os.Setenv("POSTGRES_PASSWORD", testConfig.get("POSTGRES_PASSWORD"))
	os.Setenv("POSTGRES_HOST", testConfig.get("POSTGRES_HOST"))
	os.Setenv("POSTGRES_PORT", testConfig.get("POSTGRES_PORT"))
	os.Setenv("POSTGRES_DB", testConfig.get("POSTGRES_DB"))

	db, err := sql.Open("postgres", testConfig.POSTGRES_LOCATION())
	if err != nil {
		panic(err)
	}
	suite.db = db
}

func (suite *BackendTestSuite) TearDownSuite() {
}

func (suite *BackendTestSuite) SetupTest() {
	cmd := "/__w/queue-system-lite/queue-system-lite/scripts/migration_tools/migration.sh up"
	dbDown := exec.Command("sh", "-c", cmd)
	_, err := dbDown.Output()
	if err != nil {
		panic(err)
	}
}

func (suite *BackendTestSuite) TearDownTest() {
	cmd := "echo y | /__w/queue-system-lite/queue-system-lite/scripts/migration_tools/migration.sh down"
	dbDown := exec.Command("sh", "-c", cmd)
	_, err := dbDown.Output()
	if err != nil {
		panic(err)
	}
}
