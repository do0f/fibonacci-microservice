package rpc_test

import (
	"context"
	"fibonacci_service/pkg/rpc"
	fib_service_mock "fibonacci_service/pkg/server/mocks"
	server "fibonacci_service/pkg/server/rpc"
	"fibonacci_service/pkg/service"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	grpc "google.golang.org/grpc"
)

func TestServer_GetFibonacci(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := fib_service_mock.NewMockFibService(mockCtrl)

	mockService.EXPECT().FibSequence(1, 4).Return([]service.FibNumber{
		{Count: 1, Value: big.NewInt(1)},
		{Count: 2, Value: big.NewInt(2)},
		{Count: 3, Value: big.NewInt(3)},
		{Count: 4, Value: big.NewInt(5)},
	}, nil).Times(1)
	mockService.EXPECT().FibSequence(6, 7).Return(nil, service.ErrCacheError).Times(1)

	serv := server.New(mockService)

	go func() {
		err := serv.StartRpc(9001)
		if err != nil {
			assert.Fail(t, "server failed to start")
		}
	}()

	time.Sleep(1 * time.Second)

	type input struct {
		first int64
		last  int64
	}
	var tests = []struct {
		name          string
		testInput     input
		expectedError bool
	}{
		{name: "normal count", testInput: input{first: 1, last: 4}, expectedError: false},

		{name: "negative count", testInput: input{first: -5, last: 4}, expectedError: true},

		{name: "first larger than last", testInput: input{first: 8, last: 1}, expectedError: true},

		{name: "internal error", testInput: input{first: 6, last: 7}, expectedError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn, err := grpc.Dial("localhost:9001", grpc.WithInsecure())
			if err != nil {
				log.Fatalf("fail to dial: %v", err)
			}
			defer conn.Close()
			client := rpc.NewFibonacciClient(conn)

			res, err := client.GetFibonacci(context.Background(), &rpc.FibonacciSequenceRequest{
				First: test.testInput.first,
				Last:  test.testInput.last,
			})

			assert.Equal(t, err, nil, "error not nil %v", err)
			assert.Equal(t, res.Error != "", test.expectedError)
		})
	}
}
