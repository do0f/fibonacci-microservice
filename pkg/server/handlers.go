package server

import (
	"errors"
	"math/big"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrFirstLargerThanLast = errors.New("first is larger that last")
	ErrNegativeCount       = errors.New("fibonacci numbers count should start with 0")
	ErrBadQuery            = errors.New("invalid query")
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

func (serv *Server) GetFibonacci(ctx echo.Context) error {
	reqData := &fibbonaciReqData{}
	err := ctx.Bind(reqData)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse{ErrBadQuery.Error()})
	}

	if reqData.First > reqData.Last {
		return ctx.JSON(http.StatusBadRequest, errorResponse{ErrFirstLargerThanLast.Error()})
	}
	if reqData.First < 0 || reqData.Last < 0 {
		return ctx.JSON(http.StatusBadRequest, errorResponse{ErrNegativeCount.Error()})
	}

	sequence, err := serv.svc.FibSequence(reqData.First, reqData.Last)
	if err != nil {
		serv.Logger.Error(err)
		return ctx.JSON(http.StatusInternalServerError, errorResponse{"interval error"})
	}

	seqValue := fibbonaciRespData{Sequence: make([]*big.Int, len(sequence))}

	for i, v := range sequence {
		seqValue.Sequence[i] = v.Value
	}

	return ctx.JSON(http.StatusOK, seqValue)
}
