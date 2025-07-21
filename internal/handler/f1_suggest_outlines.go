package handler

import (
	"context"
	"darius/internal/constants"
	"darius/internal/errors"
	"darius/pkg/proto/suggest"
	"encoding/json"
	"fmt"
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
	_, llmResponse, err := h.llmManager.Generate(ctx, constants.F1_SUGGEST_OUTLINES, prompt, "", nil)
	if err != nil {
		return nil, errors.Error(errors.ErrNetworkConnection)
	}
	parsedResponse, err := sanitizeJSON(llmResponse)
	if err != nil {
		return nil, errors.Error(errors.ErrJSONParsing)
	}
	// Convert the parsed response to the expected format
	var outlines = &suggest.SuggestOutlinesResponse{}
	err = json.Unmarshal([]byte(parsedResponse), &outlines)
	if err != nil {
		return nil, errors.Error(errors.ErrJSONUnmarshalling)
	}

	return outlines, nil
}
