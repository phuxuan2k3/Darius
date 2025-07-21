package database

import (
	"context"
	"darius/cmd/db"
	"log"
)

type Service interface {
	CreateLLMCallReport(context.Context, string, string, string, float64) (string, error)
}

type service struct {
	db db.Database
}

func NewService(db db.Database) Service {
	return &service{
		db: db,
	}
}

func (s *service) CreateLLMCallReport(ctx context.Context, entry, res, resp string, amount float64) (string, error) {
	if s.db == nil {
		log.Print("Database service is not initialized")
		return "", nil
	}
	return s.db.CreateReport(
		entry,
		res,
		resp,
		amount,
	)
}
