package service

import (
	"errors"
	"fibonacci_service/pkg/cache"
	"math/big"
)

var (
	ErrNegativeCount = errors.New("fibonacci numbers count should start with 0")
)

type ICache interface {
	GetFibonacci(count int) (cache.FibNumber, error)
	SetFibonacci(num cache.FibNumber) error
}

type FibService struct {
	c ICache
}

func New(cache ICache) *FibService {
	return &FibService{
		c: cache,
	}
}

type FibNumber struct {
	Count int
	Value *big.Int
}

func (fs *FibService) getFib(count int) (FibNumber, error) {
	if count < 0 {
		return FibNumber{}, ErrNegativeCount
	}
	if count == 0 || count == 1 {
		return FibNumber{count, big.NewInt(1)}, nil
	}

	cachedNum, err := fs.c.GetFibonacci(count)
	//number was in cache
	if err == nil {
		return FibNumber{cachedNum.Count, cachedNum.Value}, nil
	}
	//cache error
	if err != cache.ErrKeyDoesntExist {
		return FibNumber{}, err
	}

	//number was not in cache
	prev := [2]FibNumber{}
	for i := 0; i < 2; i++ {
		var err error
		prev[i], err = fs.getFib(count - 1 - i)

		if err != nil {
			return FibNumber{}, err
		}
	}

	fibNum := FibNumber{count, big.NewInt(0).Add(prev[0].Value, prev[1].Value)}
	//cache new number
	go fs.c.SetFibonacci(cache.FibNumber{Count: fibNum.Count, Value: fibNum.Value})

	return fibNum, nil
}

func (fs *FibService) FibSequence(first int, last int) ([]FibNumber, error) {
	sequence := make([]FibNumber, last-first+1)

	for i := 0; i < last-first+1; i++ {
		num, err := fs.getFib(first + i)
		if err != nil {
			return nil, err
		}

		sequence[i] = num
	}

	return sequence, nil
}
