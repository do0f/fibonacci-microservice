package server

import (
	"context"
	"errors"
	"fibonacci_service/pkg/service"
)

var (
	ErrFirstLargerThanLast = errors.New("first is larger that last")
	ErrNegativeCount       = errors.New("fibonacci numbers count should start with 0")
	ErrBadQuery            = errors.New("invalid query")
)

type FibService interface {
	FibSequence(ctx context.Context, first int, last int) ([]service.FibNumber, error)
}
