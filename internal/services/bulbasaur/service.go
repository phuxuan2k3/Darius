package bulbasaur

import (
	"context"
	"darius/internal/errors"
	"darius/pkg/proto/deps/bulbasaur"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	CheckCallingLLM(ctx context.Context, uid uint64, amount float32, description string) (string, error)
	ChargeCallingLLM(ctx context.Context, code string) bool
}

type service struct {
	client bulbasaur.VenusaurClient
}

func NewService(client bulbasaur.VenusaurClient) Service {
	return &service{
		client: client,
	}
}

func (s *service) CheckCallingLLM(ctx context.Context, uid uint64, amount float32, description string) (string, error) {
	res, err := s.client.StartTransaction(ctx, &bulbasaur.StartTransactionRequest{
		UserId: uid,
		Amount: amount,
		Note:   description,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Internal {
			log.Printf("[CheckCallingLLM] 500 Internal Server Error: %v", err)
			return "", errors.Error(errors.ErrNotEnoughCredits)
		}
		log.Printf("[CheckCallingLLM] Error checking calling LLM: %v", err)
		return "", errors.Error(errors.ErrNotEnoughCredits)
	}

	log.Printf("[CheckCallingLLM]: UserID: %d, Amount: %f, Description: %s, Response: %+v", uid, amount, description, res)

	if res.GetTransactionCode() == "" || len(res.GetTransactionCode()) == 0 {
		log.Printf("[CheckCallingLLM] Empty transaction code received for UserID: %d with amount: %f", uid, amount)
		return "", errors.Error(errors.ErrGeneral)
	}

	return res.GetTransactionCode(), nil
}
func (s *service) ChargeCallingLLM(ctx context.Context, code string) bool {
	res, err := s.client.CommitTransaction(ctx, &bulbasaur.CommitTransactionRequest{
		TransactionCode: code,
	})

	if err != nil {
		log.Printf("[ChargeCallingLLM] Error charging calling LLM: %v", err)
		return false
	}

	log.Printf("[ChargeCallingLLM]: Transaction Code: %s, Response: %+v", code, res)

	return res.GetSuccess()
}
