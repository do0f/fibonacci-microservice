package service

import (
	"context"
	"errors"
	"fibonacci_service/pkg/cache"
	"math/big"
)

var (
	ErrFirstLargerThanLast = errors.New("first is larger that last")
	ErrNegativeCount       = errors.New("fibonacci numbers count should start with 0")
	ErrCacheError          = errors.New("internal error while using cache")
)

type Cache interface {
	GetFibonacci(ctx context.Context, count int) (cache.FibNumber, error)
	SetFibonacci(ctx context.Context, num cache.FibNumber) error
}

type FibService struct {
	c Cache
}

func New(cache Cache) *FibService {
	return &FibService{
		c: cache,
	}
}

type FibNumber struct {
	Count int
	Value *big.Int
}

func (fs *FibService) getFib(ctx context.Context, count int) (FibNumber, error) {
	select {
	case <-ctx.Done():
		return FibNumber{}, ctx.Err()
	default:
	}

	if count < 0 {
		return FibNumber{}, ErrNegativeCount
	}
	if count == 0 || count == 1 {
		return FibNumber{count, big.NewInt(1)}, nil
	}

	cachedNum, err := fs.c.GetFibonacci(ctx, count)
	//number was in cache
	if err == nil {
		return FibNumber{cachedNum.Count, cachedNum.Value}, nil
	}
	if err == context.Canceled {
		return FibNumber{}, err
	}
	//return error only when it is internal
	if err != cache.ErrKeyDoesntExist {
		return FibNumber{}, ErrCacheError
	}

	//number was not in cache
	prev := [2]FibNumber{}
	for i := 0; i < 2; i++ {
		var err error
		prev[i], err = fs.getFib(ctx, count-1-i)

		if err != nil {
			return FibNumber{}, err
		}
	}

	fibNum := FibNumber{count, big.NewInt(0).Add(prev[0].Value, prev[1].Value)}
	//cache new number
	go fs.c.SetFibonacci(ctx, cache.FibNumber{Count: fibNum.Count, Value: fibNum.Value})

	return fibNum, nil
}

func (fs *FibService) FibSequence(ctx context.Context, first int, last int) ([]FibNumber, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if first > last {
		return nil, ErrFirstLargerThanLast
	}

	sequence := make([]FibNumber, last-first+1)

	for i := 0; i < last-first+1; i++ {
		num, err := fs.getFib(ctx, first+i)
		if err != nil {
			return nil, err
		}

		sequence[i] = num
	}

	return sequence, nil
}
