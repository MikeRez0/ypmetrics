package http_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	handlers "github.com/MikeRez0/ypmetrics/internal/api/http"
	"github.com/MikeRez0/ypmetrics/internal/logger"
	"github.com/MikeRez0/ypmetrics/internal/service"
	"github.com/MikeRez0/ypmetrics/internal/storage"
	"github.com/go-resty/resty/v2"
)

func ExampleHandler() {
	l := logger.GetLogger("info")
	serv, err := service.NewMetricService(storage.NewMemStorage(), l)
	if err != nil {
		log.Fatal(err)
	}

	mh, err := handlers.NewMetricsHandler(serv, l)
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
