package handler

import (
	"context"
	"darius/internal/errors"
	"log"
	"time"
)

const (
	maxRetries = 3
	baseDelay  = 1 * time.Second
)

type ParseFunction interface {
	Parse(input string) (interface{}, error)
}

func (h *handler) retryCallLLM(ctx context.Context, entry string, prompt string, parseFunc ParseFunction) (interface{}, error) {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		delay := time.Duration(attempt) * baseDelay
		log.Printf("[retryCallLLM] Retry attempt %d/%d after %v delay", attempt, maxRetries, delay)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		llmResponse, err := h.llmManager.Generate(ctx, entry, prompt)
		if err != nil {
			log.Printf("[retryCallLLM] LLM call failed on attempt %d: %v", attempt, err)
			if attempt == maxRetries {
				return nil, err
			}
			continue
		}

		result, err := parseFunc.Parse(llmResponse)
		if err == nil {
			log.Printf("[retryCallLLM] Successfully parsed response on attempt %d", attempt)
			return result, nil
		}

		log.Printf("[retryCallLLM] Failed to parse response on attempt %d: %v", attempt, err)
	}

	log.Printf("[retryCallLLM] All %d retry attempts failed", maxRetries)
	return nil, errors.Error(errors.ErrJSONParsing)
}
