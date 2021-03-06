package rest

import (
	"context"
	"fibonacci_service/pkg/server"
	"math/big"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

type fibbonaciReqData struct {
	First int `query:"first"`
	Last  int `query:"last"`
}

type fibbonaciRespData struct {
	Sequence []*big.Int `json:"slice of sequence"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func (serv *Server) GetFibonacciHandler(ctx echo.Context) error {
	reqData := &fibbonaciReqData{}
	err := ctx.Bind(reqData)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{server.ErrBadQuery.Error()})
	}

	if reqData.First > reqData.Last {
		return ctx.JSON(http.StatusBadRequest, errorResponse{server.ErrFirstLargerThanLast.Error()})
	}
	if reqData.First < 0 || reqData.Last < 0 {
		return ctx.JSON(http.StatusBadRequest, errorResponse{server.ErrNegativeCount.Error()})
	}

	sequence, err := serv.svc.FibSequence(ctx.Request().Context(), reqData.First, reqData.Last)
	if err != nil {
		if err == context.Canceled {
			serv.e.Logger.Info(err)
			return ctx.JSON(http.StatusBadRequest, errorResponse{err.Error()})
		}

		serv.e.Logger.Error(err)
		return ctx.JSON(http.StatusInternalServerError, errorResponse{"interval error"})
	}

	seqValue := fibbonaciRespData{Sequence: make([]*big.Int, len(sequence))}

	for i, v := range sequence {
		seqValue.Sequence[i] = v.Value
	}

	return ctx.JSON(http.StatusOK, seqValue)
}
