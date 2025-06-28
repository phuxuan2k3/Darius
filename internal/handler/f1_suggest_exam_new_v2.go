package handler

import (
	"context"
	ctxdata "darius/ctx"
	"darius/internal/constants"
	"darius/internal/errors"
	"darius/pkg/proto/suggest"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

func (h *handler) SuggestExamQuestionV2(ctx context.Context, req *suggest.SuggestExamQuestionRequest) (resp *suggest.SuggestExamQuestionResponseV2, err error) {
	chargeCode, err := h.checkCanCall(ctx, constants.F1_SUGGEST_EXAM)
	if err != nil {
		return nil, err
	}

	// questionsContents, err := h.missfortune.GetExamQuestionContent(ctx, converters.ConvertExamRequestToMissfortuneRequest(ctx, req))
	// if err != nil {
	// 	log.Printf("[SuggestExamQuestion] error getting exam question content: %v", err)
	// 	return h.SuggestExamQuestionLegacy(ctx, req)
	// }

	prompt := fmt.Sprintf(`
You are an expert exam question designer. Your task is to generate a diverse set of high-quality exam questions based on the structured input below. Each question must be either a multiple-choice question (MCQ) or a long-answer (essay-style) question.

📤 Output Format
Return a **valid JSON object** with an array of questions. Each question must strictly follow the schema below:

{
  "questions": [
    {
      "id": 1,
      "testId": "test123",
      "text": "Question text goes here",
      "points": 5,
      "type": "MCQ",
      "detail": {
        "type": "MCQ",
        "options": ["Option A", "Option B", "Option C", "Option D"],
        "correctOption": 2
      }
    },
    {
      "id": 2,
      "testId": "test123",
      "text": "Question text goes here",
      "points": 10,
      "type": "LONG_ANSWER",
      "detail": {
        "type": "LONG_ANSWER",
        "imageLinks": ["https://example.com/image1.png"],
        "extraText": "Additional instruction for candidates",
        "correctAnswer": "The ideal answer should explain..."
      }
    },
	...
  ]
}
📌 Constraints and Rules:

Each question must include: id, testId, text, points, type, and detail.

type must be either "MCQ" or "LONG_ANSWER".

The detail.type field must match the parent type field.

For "MCQ" questions:

Exactly 4 options are required.

Only one correct option (index from 0 to 3).

Options must be meaningful and grammatically consistent with the question.

For "LONG_ANSWER" questions:

Provide at least one image link or set "imageLinks": [] if not applicable.

The extraText field may include additional instructions or constraints.

The correctAnswer must clearly explain the expected full answer to earn full points.

Use the points field to reflect question difficulty (e.g., Easy = 1–3, Medium = 4–6, Hard = 7–10).

The number and type of questions should reflect the difficulty distribution in the topics input.

The JSON must be valid and contain no comments, markdown, or explanatory text — only the JSON object as specified.

Now, generate the questions based on the following input.
%v
	`, req)

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
	var exam = &suggest.SuggestExamQuestionResponseV2{}
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
