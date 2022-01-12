package rpc

import (
	"context"
	"fibonacci_service/pkg/server"
)

func (serv *Server) GetFibonacci(ctx context.Context, req *FibonacciSequenceRequest) (*FibonacciSequenceResponse, error) {
	if req.First > req.Last {
		return &FibonacciSequenceResponse{Number: nil, Error: server.ErrFirstLargerThanLast.Error()}, nil
	}
	if req.First < 0 || req.Last < 0 {
		return &FibonacciSequenceResponse{Number: nil, Error: server.ErrNegativeCount.Error()}, nil
	}

	seq, err := serv.svc.FibSequence(int(req.First), int(req.Last))
	if err != nil {
		return &FibonacciSequenceResponse{Number: nil, Error: err.Error()}, nil
	}

	seqStr := make([]string, len(seq))
	for i, v := range seq {
		seqStr[i] = v.Value.String()
	}

	return &FibonacciSequenceResponse{Number: seqStr, Error: ""}, nil
}
