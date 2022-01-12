package rest_test

import (
	fib_service_mock "fibonacci_service/pkg/server/mocks"
	server "fibonacci_service/pkg/server/rest"
	"fibonacci_service/pkg/service"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandlers_GetFibonacciHandler(t *testing.T) {
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
	go serv.StartRest(1324)

	type input struct {
		first string
		last  string
	}
	var tests = []struct {
		name         string
		testInput    input
		expectedCode int
	}{
		{name: "normal count", testInput: input{first: "1", last: "4"}, expectedCode: http.StatusOK},

		{name: "too long count", testInput: input{first: "1231243524634645634123124235353453451231", last: "1231243524634645634123124235353453451232"},
			expectedCode: http.StatusBadRequest},
		{name: "negative count", testInput: input{first: "-5", last: "4"}, expectedCode: http.StatusBadRequest},
		{name: "string count", testInput: input{first: "string", last: "5"}, expectedCode: http.StatusBadRequest},
		{name: "no count", testInput: input{first: "3", last: ""}, expectedCode: http.StatusBadRequest},

		{name: "internal error", testInput: input{first: "6", last: "7"}, expectedCode: http.StatusInternalServerError},
	}

	for _, test := range tests {
		var request *http.Request

		query := make(url.Values)
		query.Set("first", test.testInput.first)
		query.Set("last", test.testInput.last)
		request = httptest.NewRequest(http.MethodGet, "/?"+query.Encode(), nil)

		recorder := httptest.NewRecorder()
		ctx := serv.NewContext(request, recorder)
		ctx.SetPath(server.GetFibbonaciEndpoint)

		t.Run(test.name, func(t *testing.T) {
			if assert.NoError(t, serv.GetFibonacciHandler(ctx)) {
				assert.Equal(t, test.expectedCode, recorder.Code)
			}
		})
	}
}
