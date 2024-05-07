package server

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
)

func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHTTPServer()
	e := echo.New()

	e.POST("/", httpsrv.handleProduce)
	e.GET("/", httpsrv.handleConsume)

	return &http.Server{
		Addr:    addr,
		Handler: e,
	}
}

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func (s *httpServer) handleProduce(c echo.Context) error {
	var req ProduceRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	off, err := s.Log.Append(req.Record)
	if err != nil {
		return err
	}
	res := ProduceResponse{Offset: off}
	return c.JSON(http.StatusCreated, res)
}

func (s *httpServer) handleConsume(c echo.Context) error {
	var req ConsumeRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	record, err := s.Log.Read(req.Offset)
	if errors.Is(err, ErrOffsetNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return err
	}
	res := ConsumeResponse{Record: record}
	return c.JSON(http.StatusOK, res)
}
