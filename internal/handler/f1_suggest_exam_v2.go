package handler

import (
	"context"
	ctxdata "darius/ctx"
	"darius/internal/constants"
	"darius/internal/converters"
	"darius/internal/errors"
	"darius/pkg/proto/suggest"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

func (h *handler) SuggestExamQuestion(ctx context.Context, req *suggest.SuggestExamQuestionRequest) (*suggest.SuggestExamQuestionResponse, error) {
	chargeCode, err := h.checkCanCall(ctx, constants.F1_SUGGEST_EXAM)
	if err != nil {
		return nil, err
	}

	questionsContents, err := h.missfortune.GetExamQuestionContent(ctx, converters.ConvertExamRequestToMissfortuneRequest(ctx, req))
	if err != nil {
		log.Printf("[SuggestExamQuestion] error getting exam question content: %v", err)
		return h.SuggestExamQuestionLegacy(ctx, req)
	}

	prompt := fmt.Sprintf(`
	You are an expert question generator for standardized multiple-choice exams. Your task is to generate answer options for a list of exam questions. You will receive a list of questions (text-only), and for each question, you must generate exactly 4 answer options, clearly indicating which one is the correct answer.
	📥 Input Format:
You will be given an object in the following format:
{
  "questions": [
    "What is the capital of France?",
    "Which data structure uses LIFO order?",
    ...
  ]
}
  📤 Output Format:
Return your response as a JSON object that strictly follows the Protobuf schema below:
{
  "questions": [
    {
      "text": "Question text here",
      "options": [
        "Option A",
        "Option B",
        "Option C",
        "Option D"
      ],
      "points": 1,
      "correctOption": 2
    },
    ...
  ]
}
📌 Constraints and Rules:
There must be exactly 4 options per question.

Only one option must be correct; the other three must be reasonable but clearly incorrect.

The correctOption field must be an integer from 0 to 3, representing the index of the correct option.

Use a scale based on difficulty level (e.g., Easy = 1, Medium = 2, etc.) for the points field.

Make sure that:

The incorrect options are plausible and not obviously wrong.

Options are diverse in content and style, not just variations of the same word.

All options must be grammatically consistent with the question.

The output must be valid JSON with no trailing commas or formatting errors.

The order of questions must match the order of the input.

Keep answers factually accurate, and where applicable, follow commonly accepted academic standards.

📝 Example:
Input:
{
  "questions": [
    "What is the boiling point of water at sea level?"
  ]
}
Output:
{
  "questions": [
    {
      "text": "What is the boiling point of water at sea level?",
      "options": [
        "100°C",
        "90°C",
        "120°C",
        "80°C"
      ],
      "points": 1,
      "correctOption": 0
    }
  ]
}
Now, generate the answer options for the following questions:
%v
	`, questionsContents)

	llmResponse, err := h.llmManager.Generate(ctx, constants.F1_SUGGEST_EXAM, prompt)
	if err != nil {
		return nil, errors.Error(errors.ErrNetworkConnection)
	}
	log.Println("[SuggestExamQuestion] LLM response:", llmResponse)
	parsedResponse, err := sanitizeJSON(llmResponse)
	if err != nil {
		return nil, errors.Error(errors.ErrJSONParsing)
	}
	// Convert the parsed response to the expected format
	var exam = &suggest.SuggestExamQuestionResponse{}
	err = json.Unmarshal([]byte(parsedResponse), &exam)
	if err != nil {
		return nil, errors.Error(errors.ErrJSONUnmarshalling)
	}

	// Charge the user for the LLM call
	if !h.bulbasaur.ChargeCallingLLM(ctx, chargeCode) {
		log.Printf("[SuggestExamQuestion] Charge Code %s failed to charge for LLM call", chargeCode)
		return nil, errors.Error(errors.ErrChargingFailed)
	}

	return exam, nil
}

func (h *handler) checkCanCall(ctx context.Context, llmCaller string) (string, error) {
	amount, desc := constants.GetLLMCallAmount(llmCaller)
	uidStr, _ := ctxdata.GetUserIdFromContext(ctx)
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		log.Printf("[SuggestExamQuestion] error parsing user ID: %v", err)
		return "", errors.Error(errors.ErrInvalidInput)
	}
	ok, chargeCode := h.bulbasaur.CheckCallingLLM(ctx, uid, amount, desc)
	if !ok {
		log.Printf("[SuggestExamQuestion] user %d does not have enough credits to call LLM", uid)
		return "", errors.Error(errors.ErrNotEnoughCredits)
	}
	return chargeCode, nil
}
