package storage_test

import (
	"compress/gzip"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/MikeRez0/ypmetrics/internal/handlers"
	"github.com/MikeRez0/ypmetrics/internal/storage"
)

var dbtest *TestDBInstance
var l *zap.Logger

func setup() {
	var err error
	dbtest, err = NewTestDBInstance()
	if err != nil {
		log.Fatal(err)
	}
	l, err = zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
}
func shutdown() {
	if dbtest != nil {
		dbtest.Down()
	}
	err := os.Remove("test.js")
	if err != nil {
		log.Println(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func runHandlerTests(t *testing.T, router *gin.Engine) {
	t.Helper()
	srv := httptest.NewServer(router)

	tests := getTestData()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.URL = srv.URL + tt.request
			if tt.requestBody != "" {
				req.Body = tt.requestBody
			}
			if tt.contentType != "" {
				req.SetHeader("Content-Type", tt.contentType)
			}

			res, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode())

			if tt.want.body != "" {
				body := string(res.Body())
				if res.Header().Get("Content-Encoding") == "gzip" {
					gr, err := gzip.NewReader(strings.NewReader(body))
					assert.NoError(t, err)
					bodyBytes, err := io.ReadAll(gr)
					assert.NoError(t, err)
					body = string(bodyBytes)
					assert.NoError(t, gr.Close())
				}
				if tt.contentType != "application/json" {
					assert.Equal(t, tt.want.body, body)
				} else {
					assert.JSONEq(t, tt.want.body, body)
				}
			}
			if tt.want.contentType != "" {
				assert.Contains(t, res.Header().Get("Content-Type"), tt.want.contentType)
			}
		})
	}
}

func TestServerDB_Handlers(t *testing.T) {
	repo, err := storage.NewDBStorage(dbtest.DSN, l)
	assert.NoError(t, err)

	mh := &handlers.MetricsHandler{
		Store: repo,
		Log:   l,
	}

	router := handlers.SetupRouter(mh, l)
	runHandlerTests(t, router)
}

func TestServerMem_Handlers(t *testing.T) {
	repo := storage.NewMemStorage()

	mh := &handlers.MetricsHandler{
		Store: repo,
		Log:   l,
	}

	router := handlers.SetupRouter(mh, l)
	runHandlerTests(t, router)
}

func TestServerFS_Handlers(t *testing.T) {
	repo, err := storage.NewFileStorage("test.js", 0, false, l)
	assert.NoError(t, err)

	mh := &handlers.MetricsHandler{
		Store: repo,
		Log:   l,
	}

	router := handlers.SetupRouter(mh, l)
	runHandlerTests(t, router)
}
