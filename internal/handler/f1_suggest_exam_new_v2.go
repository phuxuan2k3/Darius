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

func (h *handler) SuggestExamQuestionV2(ctx context.Context, req *suggest.SuggestExamQuestionRequest) (resp *suggest.SuggestExamQuestionResponseV2, err error) {
	chargeCode, err := h.checkCanCall(ctx, constants.F1_SUGGEST_EXAM)
	if err != nil {
		return nil, err
	}

	questionsContents, err := h.missfortune.GetExamQuestionContent(ctx, converters.ConvertExamRequestToMissfortuneRequest(ctx, req))
	prompt := ""
	if err != nil {
		log.Printf("[SuggestExamQuestion] error getting exam question content: %v", err)

		prompt = fmt.Sprintf(`
You are an expert exam question designer. Your task is to generate a diverse set of high-quality exam questions based on the structured input below. Each question must be either a multiple-choice question (MCQ) or a long-answer (essay-style) question.

Before generating, follow this step-by-step reasoning to ensure quality and uniqueness:

🧠 Step-by-Step Reasoning:
1. First, generate a conceptual list of question ideas covering different aspects of the provided topics.
2. Check and confirm that each idea is distinct in content and intent (i.e., no duplicate or near-duplicate questions).
3. For each MCQ:
   - Generate 4 answer options that are factually distinct and grammatically consistent.
   - Ensure that no two options are semantically or syntactically identical or too similar.
   - The incorrect options must be plausible and non-trivial.
4. For long-answer questions:
   - Ensure the question targets higher-order thinking (e.g., analysis, explanation, comparison).
   - Provide a clear expected answer and optional visual/image links if applicable.

🔁 Self-Verification (Post-Generation Check):
- Verify that **no question text** is repeated.
- Verify that **within each MCQ**, no two options are the same or nearly the same.
- Verify that the output strictly matches the structure described below and is valid JSON.

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
- Each question must include: id, testId, text, points, type, and detail.
- type must be either "MCQ" or "LONG_ANSWER".
- The detail.type field must match the parent type field.
- For MCQs:
  - Exactly 4 options.
  - Only one correct option (index from 0 to 3).
  - Options must be clear, distinct, and grammatically aligned with the question.
- For long-answer questions:
  - Must include a model answer and optional image links.
  - Must require thoughtful, detailed explanations.
- The points field should match question difficulty: (e.g., Easy = 1–3, Medium = 4–6, Hard = 7–10).
- The output JSON must be valid and must not include explanations, notes, or markdown.
- No two questions should be identical or too similar.
- No two options in any MCQ should be the same or semantically identical.

Now, generate the questions based on the following input:
%v
	`, req)
	} else {
		prompt = generateOptionsPrompt(questionsContents)
	}

	llmResponse, err := h.llmManager.Generate(ctx, constants.F1_SUGGEST_EXAM, prompt)
	if err != nil {
		return nil, errors.Error(errors.ErrNetworkConnection)
	}
	parsedResponse, err := sanitizeJSON(llmResponse)
	if err != nil {
		return nil, errors.Error(errors.ErrJSONParsing)
	}
	// Convert the parsed response to the expected format
	var exam = &suggest.SuggestExamQuestionResponseV2{}
	err = json.Unmarshal([]byte(parsedResponse), &exam)
	if err != nil {
		log.Printf("[SuggestExamQuestion] error unmarshalling JSON: %v", err)
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
	chargeCode, err := h.bulbasaur.CheckCallingLLM(ctx, uid, amount, desc)
	if err != nil {
		return "", errors.Error(errors.ErrGeneral)
	}
	return chargeCode, nil
}
