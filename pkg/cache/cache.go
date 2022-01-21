package cache

import (
	"context"
	"errors"
	"math/big"
	"strconv"

	redis "github.com/go-redis/redis/v8"
)

var (
	ErrKeyDoesntExist = errors.New("key doesnt exist")
	ErrParsingValue   = errors.New("error parsing value")
)

type FibNumber struct {
	Count int      //key
	Value *big.Int //value
}

type Cache struct {
	c *redis.Client
}

func Connect(address string, password string, db int) (*Cache, error) {
	c := &Cache{
		c: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		}),
	}

	_, err := c.c.Ping(context.Background()).Result()
	if err != nil { // did not connect successfully
		return nil, err
	}

	return c, nil
}

func (c *Cache) GracefulShutdown() error {
	return c.c.Close()
}

func (c *Cache) GetFibonacci(count int) (FibNumber, error) {
	strVal, err := c.c.Get(context.Background(), strconv.Itoa(count)).Result()
	if err != nil {
		if err == redis.Nil {
			return FibNumber{}, ErrKeyDoesntExist
		}

		return FibNumber{}, err
	}

	val := &big.Int{}
	val, ok := val.SetString(strVal, 10)
	if !ok {
		return FibNumber{}, ErrParsingValue
	}

	return FibNumber{count, val}, nil
}

func (c *Cache) SetFibonacci(num FibNumber) error {
	err := c.c.Set(context.Background(), strconv.Itoa(num.Count),
		num.Value.String(), 0).Err()

	return err
}
