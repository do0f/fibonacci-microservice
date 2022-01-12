package service_test

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"fibonacci_service/pkg/cache"
	"fibonacci_service/pkg/service"
	mock_cache "fibonacci_service/pkg/service/mocks"
)

func TestFibService_FibSequence(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCache := mock_cache.NewMockCache(mockCtrl)
	//cache miss for any value except two for error checking case
	mockCache.EXPECT().GetFibonacci(gomock.All(gomock.Any(), gomock.Not(6), gomock.Not(7))).Return(cache.FibNumber{}, cache.ErrKeyDoesntExist).AnyTimes()
	//cache hit for count 6
	mockCache.EXPECT().GetFibonacci(6).Return(cache.FibNumber{Count: 6, Value: big.NewInt(13)}, nil).AnyTimes()
	//cache error for count 7
	mockCache.EXPECT().GetFibonacci(7).Return(cache.FibNumber{}, cache.ErrParsingValue).AnyTimes()
	//pretend to cache
	mockCache.EXPECT().SetFibonacci(gomock.Any()).Return(nil).AnyTimes()

	svc := service.New(mockCache)

	type input struct {
		first int
		last  int
	}
	type output struct {
		sequence []service.FibNumber
		err      error
	}
	var tests = []struct {
		name      string
		testInput input
		expected  output
	}{
		{name: "0 to 5 numbers", testInput: input{0, 5}, expected: output{sequence: []service.FibNumber{
			{0, big.NewInt(1)},
			{1, big.NewInt(1)},
			{2, big.NewInt(2)},
			{3, big.NewInt(3)},
			{4, big.NewInt(5)},
			{5, big.NewInt(8)},
		}, err: nil}},
		{name: "3 to 4 numbers", testInput: input{3, 4}, expected: output{sequence: []service.FibNumber{
			{3, big.NewInt(3)},
			{4, big.NewInt(5)},
		}, err: nil}},
		{name: "6 to 4 numbers", testInput: input{6, 4}, expected: output{sequence: nil, err: service.ErrFirstLargerThanLast}},
		{name: "-1 to 4 numbers", testInput: input{-1, 4}, expected: output{sequence: nil, err: service.ErrNegativeCount}},
		{name: "check cache hit", testInput: input{0, 6}, expected: output{sequence: []service.FibNumber{
			{0, big.NewInt(1)},
			{1, big.NewInt(1)},
			{2, big.NewInt(2)},
			{3, big.NewInt(3)},
			{4, big.NewInt(5)},
			{5, big.NewInt(8)},
			{6, big.NewInt(13)},
		}, err: nil}},
		{name: "check cache error", testInput: input{6, 7}, expected: output{sequence: nil, err: service.ErrCacheError}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			seq, err := svc.FibSequence(test.testInput.first, test.testInput.last)

			assert.Equal(t, test.expected.err, err)
			assert.Equal(t, len(test.expected.sequence), len(seq))
			for i := 0; i < len(seq); i++ {
				assert.Zero(t, test.expected.sequence[i].Value.Cmp(seq[i].Value))
			}
		})
	}
}
