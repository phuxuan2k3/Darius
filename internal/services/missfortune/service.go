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

	"go.uber.org/zap"
)

const (
	URL_GetExamQuestionContent = "/generate"
)

type Service interface {
	GetExamQuestionContent(ctx context.Context, req *missfortune.SuggestExamQuestionRequest) (*missfortune.SuggestExamQuestionResponse, error)
}

type service struct {
	address    string
	httpClient *http.Client
}

func NewService(address string, httpClient *http.Client) Service {
	return &service{
		address:    address,
		httpClient: httpClient,
	}
}

func (s *service) GetExamQuestionContent(ctx context.Context, req *missfortune.SuggestExamQuestionRequest) (*missfortune.SuggestExamQuestionResponse, error) {

	jsonBody, err := json.Marshal(req)
	if err != nil {
		log.Printf("[MFT][GetExamQuestionContent] Error marshalling request body: %v, \n MFT body: %v", err, req)
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	zap.L().Error(
		"[MFT][GetExamQuestionContent] Marshal request failed",
		zap.Error(err),
		zap.Reflect("request", req),
		zap.String("jsonBody", string(jsonBody)),
	)

	bodyReader := bytes.NewReader(jsonBody)

	requestURL := s.address + URL_GetExamQuestionContent
	httpReq, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		log.Printf("[MFT][GetExamQuestionContent] Error creating HTTP request: %v, \n MFT body: %v", err, httpReq)
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	httpResp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("[MFT][GetExamQuestionContent] Error making HTTP request: %v,\n MFT body: %v", err, httpReq)
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		log.Printf("[MFT][GetExamQuestionContent] HTTP request failed with status code: %d\n MFT body: %v", httpResp.StatusCode, httpReq)
		return nil, fmt.Errorf("HTTP request failed with status code: %d", httpResp.StatusCode)
	}

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		log.Printf("[MFT][GetExamQuestionContent] Error reading response body: %v, \n MFT body: %v", err, httpReq)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	questionsContent := &missfortune.SuggestExamQuestionResponse{}

	if err := json.Unmarshal(respBody, questionsContent); err != nil {
		log.Printf("[MFT][GetExamQuestionContent] Error unmarshalling response body: %v, \n MFT body: %v", err, httpReq)
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return questionsContent, nil
}
