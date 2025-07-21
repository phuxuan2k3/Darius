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
	var conversationId *uint64 = nil

	for attempt := 0; attempt <= maxRetries; attempt++ {
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

		_, conversationIdResp, llmResponse, err := h.llmManager.Generate(ctx, entry, prompt, conversationId)
		conversationId = conversationIdResp

		if err != nil {
			log.Printf("[retryCallLLM] LLM call failed on attempt %d: %v", attempt, err)
			if attempt == maxRetries {
				return nil, err
			}

			if conversationId != nil {
				prompt = "Please try to generate the response again. The previous attempt has just met the error: " + err.Error()
			}

			continue
		}

		result, err := parseFunc.Parse(llmResponse)
		if err == nil {
			log.Printf("[retryCallLLM] Successfully parsed response on attempt %d", attempt)
			return result, nil
		}

		log.Printf("[retryCallLLM] Failed to parse response on attempt %d: %v", attempt, err)
		prompt = "I could not parse the response. Please try to generate the response again. The previous attempt has just met the error: " + err.Error()
	}

	log.Printf("[retryCallLLM] All %d retry attempts failed", maxRetries)
	return nil, errors.Error(errors.ErrJSONParsing)
}
