package main

import (
	"net/http"

	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
)

type JsonResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Content []*metricFamily `json:"content,omitempty"`
}

func BadRequestJSON(c echo.Context, errMsg string) error {
	c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "http://localhost:3000")
	return c.JSON(
		http.StatusBadRequest, &JsonResponse{
			Code:    http.StatusBadRequest,
			Message: errMsg,
		},
	)
}

func fileServeHandler(c echo.Context) error {
	file, err := ioutil.ReadFile(fileToServe)
	if err != nil {
		BadRequestJSON(c, err.Error())
	}

	return c.String(http.StatusOK, string(file))
}


func ServeHTTP(port string) {
	fmt.Printf("Starting server on port %s \n", port)
	// Echo instance
	e := echo.New()
	e.Use(middleware.CORS(), middleware.Logger())

	// Set a default error handler to be JSON as well.
	echo.NotFoundHandler = func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, &JsonResponse{
			Code:    http.StatusNotFound,
			Message: "Endpoint not found",
		})
	}

	// Routes.
	e.GET("/_metrics.json", prom2jsonHandler)
	e.GET("/_metrics", fileServeHandler)

	// Start server.
	e.Run(fasthttp.New(fmt.Sprintf(":%s", port)))
}
