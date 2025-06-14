package handler

// import (
// 	"context"
// 	"darius/models"
// 	"darius/pkg/proto/suggest"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// )

// func (h *handler) SuggestExamQuestion(ctx context.Context, req *suggest.SuggestExamQuestionRequest) (*suggest.SuggestExamQuestionResponse, error) {
// 	prompt := fmt.Sprintf(`
// 	You are an expert exam question designer. Generate high-quality multiple-choice questions (MCQs) based on the structured input below. The output must strictly follow the specified format and constraints.

// ---

// Exam Metadata
// Title: %v
// Description: %v
// Language: %v (All content must be in this language)
// Target Seniority Level: %v
// Creativity Level: %v (Scale from 1–10. Higher values mean more original or less conventional question styles.)

// ---

// Topics and Difficulty Distribution
// %v

// Note: Difficulty levels may include "Easy", "Medium", "Hard", "Very hard", etc. These levels are dynamically defined per topic.

// ---

// Context
// Textual Context:
// %v

// Relevant Links:
// %v
// ---

// Output Format
// Respond with a JSON object containing an array of question objects:
// {
//   "questions": [
//     {
//       "text": "Question goes here",
//       "options": ["Option A", "Option B", "Option C", "Option D"],
//       "points": number, // Use a scale based on difficulty level (e.g., Easy = 1, Medium = 2, etc.)
//       "correctOption": number // Index (0–3) of the correct option
//     },
//     ...
//   ]
// }
// 	`, req.GetTitle(), req.GetDescription(), req.GetLanguage(), req.GetSeniority(), req.GetCreativity(), req.GetTopics(), req.GetContext().GetText(), req.GetContext().GetLinks())

// 	llmResponse, err := h.llmManager.Generate(ctx, models.F1_SUGGEST_EXAM, prompt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Println("[SuggestExamQuestion] LLM response:", llmResponse)
// 	parsedResponse, err := sanitizeJSON(llmResponse)
// 	if err != nil {
// 		return nil, fmt.Errorf("[SuggestExamQuestion] error parsing response: %v", err)
// 	}
// 	// Convert the parsed response to the expected format
// 	var exam = &suggest.SuggestExamQuestionResponse{}
// 	err = json.Unmarshal([]byte(parsedResponse), &exam)
// 	if err != nil {
// 		return nil, fmt.Errorf("[SuggestExamQuestion] error unmarshalling outlines: %v", err)
// 	}

// 	return exam, nil
// }
