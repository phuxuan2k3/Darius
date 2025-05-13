package handler

import (
	"context"
	"darius/models"
	"darius/pkg/proto/suggest"
	"encoding/json"
	"fmt"
	"log"
)

func (h *handler) SuggestOutlines(ctx context.Context, req *suggest.SuggestOutlinesRequest) (*suggest.SuggestOutlinesResponse, error) {

	prompt := fmt.Sprintf(`
	You are an assistant that generates topic outlines for multiple-choice tests.
Given the following test information, suggest 1 to 3 concise and relevant topic outlines that are not already present in the list.
The new outlines should match the test's title, description, difficulty, and tags. Avoid duplicates or near-duplicates of the existing outlines.

Test information:

Title: %v

Description: %v

Difficulty: %v

Tags: %v

Existing Outlines:
%v

Output format:
{
  "outlines": [
    "First new outline idea",
    "Second new outline idea (optional)",
    "Third new outline idea (optional)"
  ]
}`, req.GetTitle(), req.GetDescription(), req.GetDifficulty(), req.GetTags(), req.GetOutlines())
	llmResponse, err := h.llmManager.Generate(ctx, models.F1, prompt)
	if err != nil {
		return nil, err
	}
	log.Println("[SuggestOutlines] LLM response:", llmResponse)
	parsedResponse, err := sanitizeJSON(llmResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}
	// Convert the parsed response to the expected format
	var outlines []string
	err = json.Unmarshal([]byte(parsedResponse), &outlines)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling outlines: %v", err)
	}

	return &suggest.SuggestOutlinesResponse{
		Outlines: outlines,
	}, nil
}
