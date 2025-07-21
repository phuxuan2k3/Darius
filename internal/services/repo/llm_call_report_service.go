package database

import (
	"context"
	"darius/cmd/db"
	"log"
)

type Service interface {
	CreateLLMCallReport(context.Context, string, string, string, string, float64) error
	GetByRequestKey(context.Context, string) (string, error)
}

type service struct {
	db db.Database
}

func NewService(db db.Database) Service {
	return &service{
		db: db,
	}
}

func (s *service) CreateLLMCallReport(ctx context.Context, entry, res, resp, requestKey string, amount float64) error {
	if s.db == nil {
		log.Print("Database service is not initialized")
		return nil
	}
	return s.db.CreateReport(
		entry,
		res,
		resp,
		requestKey,
		amount,
	)
}

func (s *service) GetByRequestKey(ctx context.Context, requestKey string) (string, error) {
	if s.db == nil {
		log.Print("Database service is not initialized")
		return "", nil
	}
	report, err := s.db.GetByRequestKey(requestKey)
	if err != nil {
		log.Printf("Error getting report by request key: %v", err)
		return "", err
	}
	return report.Resp, nil
}
