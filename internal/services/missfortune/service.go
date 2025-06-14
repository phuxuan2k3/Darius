package missfortune

import (
	"bytes"
	"context"
	"darius/pkg/proto/deps/missfortune"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	URL_GetExamQuestionContent = "/SuggestExamQuestion"
)

type Service interface {
	GetExamQuestionContent(ctx context.Context, req *missfortune.SuggestExamQuestionRequest) (*missfortune.SuggestExamQuestionResponse, error)
}

type service struct {
	address string
}

func NewService(address string) Service {
	return &service{
		address: address,
	}
}

func (s *service) GetExamQuestionContent(ctx context.Context, req *missfortune.SuggestExamQuestionRequest) (*missfortune.SuggestExamQuestionResponse, error) {
	jsonBody := []byte(req.String())
	bodyReader := bytes.NewReader(jsonBody)

	requestURL := s.address + URL_GetExamQuestionContent
	httpReq, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	httpResp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		log.Printf("HTTP request failed with status code: %d", httpResp.StatusCode)
		return nil, fmt.Errorf("HTTP request failed with status code: %d", httpResp.StatusCode)
	}

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	questionsContent := &missfortune.SuggestExamQuestionResponse{}

	if err := json.Unmarshal(respBody, questionsContent); err != nil {
		log.Printf("Error unmarshalling response body: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return questionsContent, nil
}
