package handlers_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/MikeRez0/ypmetrics/internal/handlers"
	"github.com/MikeRez0/ypmetrics/internal/storage"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func ExampleHandler() {
	mh := &handlers.MetricsHandler{
		Store: storage.NewMemStorage(),
	}

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	router := handlers.SetupRouter(mh, l, nil)
	srv := httptest.NewServer(router)

	req := resty.New().R()
	req.Method = http.MethodPost
	reqPath := "/update/gauge/test/5"
	req.URL = srv.URL + reqPath

	res, err := req.Send()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Request: %v %v\n", req.Method, reqPath)
	fmt.Printf("Response status: %v\n", res.StatusCode())

	req = resty.New().R()
	req.Method = http.MethodGet
	reqPath = "/value/gauge/test"
	req.URL = srv.URL + reqPath

	res, err = req.Send()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Request: %v %v\n", req.Method, reqPath)
	fmt.Printf("Response status: %v\n", res.StatusCode())
	fmt.Printf("Response body: %v\n", string(res.Body()))

	// Output:
	// Request: POST /update/gauge/test/5
	// Response status: 200
	// Request: GET /value/gauge/test
	// Response status: 200
	// Response body: 5
}
