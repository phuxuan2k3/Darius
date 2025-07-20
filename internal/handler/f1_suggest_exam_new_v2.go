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

// handleErrorWithStatusCode sets the appropriate HTTP status code and returns the error
func (h *handler) handleErrorWithStatusCode(ctx context.Context, _ error, errorType string) error {
	appError := errors.Error(errorType)
	ctxdata.SetHeaders(ctx, ctxdata.HttpCodeHeader, errors.GetHTTPStatusCode(appError))
	return appError
}

func (h *handler) SuggestExamQuestionV2(ctx context.Context, req *suggest.SuggestExamQuestionRequest) (resp *suggest.SuggestExamQuestionResponseV2, err error) {
	if req.GetQuestionType() == "" || len(req.GetQuestionType()) == 0 {
		req.QuestionType = "MIXED" //default question type
	}

	chargeCode, err := h.checkCanCall(ctx, constants.F1_SUGGEST_EXAM)
	if err != nil {
		return nil, err
	}
	log.Printf("[MFT] req: %+v", converters.ConvertExamRequestToMissfortuneRequest(ctx, req))
	questionsContents, err := h.missfortune.GetExamQuestionContent(ctx, converters.ConvertExamRequestToMissfortuneRequest(ctx, req))
	prompt := ""
	if err != nil {
		log.Printf("[SuggestExamQuestion] error getting exam question content: %v", err)

		templateRule := "- Topic: **%v**, Level: **%v**, Quantity: **%v**\n"
		questionCount := 0
		instruction := ""
		for _, topic := range req.Topics {
			if topic.GetDifficultyDistribution().GetIntern() > 0 {
				instruction += fmt.Sprintf(templateRule, topic.GetName(), "Intern", topic.GetDifficultyDistribution().GetIntern())
				questionCount += int(topic.GetDifficultyDistribution().GetIntern())
			}
			if topic.GetDifficultyDistribution().GetJunior() > 0 {
				instruction += fmt.Sprintf(templateRule, topic.GetName(), "Junior", topic.GetDifficultyDistribution().GetJunior())
				questionCount += int(topic.GetDifficultyDistribution().GetJunior())
			}
			if topic.GetDifficultyDistribution().GetMiddle() > 0 {
				instruction += fmt.Sprintf(templateRule, topic.GetName(), "Middle", topic.GetDifficultyDistribution().GetMiddle())
				questionCount += int(topic.GetDifficultyDistribution().GetMiddle())
			}
			if topic.GetDifficultyDistribution().GetSenior() > 0 {
				instruction += fmt.Sprintf(templateRule, topic.GetName(), "Senior", topic.GetDifficultyDistribution().GetSenior())
				questionCount += int(topic.GetDifficultyDistribution().GetSenior())
			}
			if topic.GetDifficultyDistribution().GetLead() > 0 {
				instruction += fmt.Sprintf(templateRule, topic.GetName(), "Lead", topic.GetDifficultyDistribution().GetLead())
				questionCount += int(topic.GetDifficultyDistribution().GetLead())
			}
			if topic.GetDifficultyDistribution().GetExpert() > 0 {
				instruction += fmt.Sprintf(templateRule, topic.GetName(), "Expert", topic.GetDifficultyDistribution().GetExpert())
				questionCount += int(topic.GetDifficultyDistribution().GetExpert())
			}
		}
		req.Topics = nil // Clear topics to avoid duplication in the prompt
		prompt = fmt.Sprintf(`
You are an expert exam question designer. Your task is to generate exactly **%v diverse and high-quality exam questions** based on the structured requirements below. Each question must be either a multiple-choice question (MCQ) or a long-answer (essay-style) question.

---

📌 Required Question Breakdown (Strict):
%v
✅ Total: **Exactly %v questions** — no more, no less.

---

🧠 Reasoning Steps (Quality Assurance):
1. Generate a distinct and relevant idea for each question based on its topic and level that match with the question type.
	- If questionType is "MCQ", generate all questions as Multiple Choice Questions (MCQ).
	- If questionType is "LONG_ANSWER", generate all questions as Essay-style questions (LONG_ANSWER).
	- If questionType is "MIXED", generate a balanced mix of MCQ and LONG_ANSWER questions.
2. Ensure all %v questions are unique in wording and intent (no duplication).
3. For MCQs:
   - Provide exactly 4 options.
   - All options must be grammatically aligned, factually plausible, and **clearly distinct** from one another.
   - One option must be clearly correct, indicated by "correctOption" (index 0–3).
4. For Long Answer:
   - Require deep reasoning, explanation, or comparison.
   - Include a clear, complete expected answer ("correctAnswer").
   - Use "imageLinks" if relevant, or leave it as an empty array.

---

🔁 Final Validation (Self-Verification):
- Confirm that exactly %v questions are generated, matching the exact breakdown.
- Confirm that **no two questions are identical or overlapping** in content.
- Confirm that all questions are relevant to the specified topics and levels.
- Confirm that all questions must match the specified question type.
- Confirm that output is valid JSON, with no notes, markdown, or trailing commas.

---

📤 Output Format:
Return only a valid JSON object structured like this:

{
  "questions": [
    {
      "id": 1,
      "testId": "test123",
      "text": "Question text here",
      "points": 2,
      "type": "MCQ",
      "detail": {
        "type": "MCQ",
        "options": ["A", "B", "C", "D"],
        "correctOption": 2
      }
    },
    {
      "id": 2,
      "testId": "test123",
      "text": "Question text here",
      "points": 5,
      "type": "LONG_ANSWER",
      "detail": {
        "type": "LONG_ANSWER",
        "imageLinks": [],
        "extraText": "Instructions here",
        "correctAnswer": "Expected answer here"
      }
    }
    ...
  ]
}
Now, generate questions base on the following requirements: %v
	`, questionCount, instruction, questionCount, questionCount, questionCount, req)
	} else {
		prompt = generateOptionsPrompt(questionsContents)
	}

	_, llmResponse, err := h.llmManager.Generate(ctx, constants.F1_SUGGEST_EXAM, prompt, nil)
	if err != nil {
		return nil, h.handleErrorWithStatusCode(ctx, err, errors.ErrNetworkConnection)
	}
	parsedResponse, err := sanitizeJSON(llmResponse)
	if err != nil {
		return nil, h.handleErrorWithStatusCode(ctx, err, errors.ErrJSONParsing)
	}
	// Convert the parsed response to the expected format
	var exam = &suggest.SuggestExamQuestionResponseV2{}
	err = json.Unmarshal([]byte(parsedResponse), &exam)
	if err != nil {
		log.Printf("[SuggestExamQuestion] error unmarshalling JSON: %v", err)
		return nil, h.handleErrorWithStatusCode(ctx, err, errors.ErrJSONUnmarshalling)
	}

	// Charge the user for the LLM call
	if !h.bulbasaur.ChargeCallingLLM(ctx, chargeCode) {
		log.Printf("[SuggestExamQuestion] Charge Code %s failed to charge for LLM call", chargeCode)
		return nil, h.handleErrorWithStatusCode(ctx, err, errors.ErrChargingFailed)
	}

	return exam, nil
}

func (h *handler) checkCanCall(ctx context.Context, llmCaller string) (string, error) {
	amount, desc := constants.GetLLMCallAmount(llmCaller)
	uidStr, _ := ctxdata.GetUserIdFromContext(ctx)
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		log.Printf("[SuggestExamQuestion] error parsing user ID: %v", err)
		return "", h.handleErrorWithStatusCode(ctx, err, errors.ErrInvalidInput)
	}
	chargeCode, err := h.bulbasaur.CheckCallingLLM(ctx, uid, amount, desc)
	if err != nil {
		ctxdata.SetHeaders(ctx, ctxdata.HttpCodeHeader, errors.GetHTTPStatusCode(err))
		return "", err
	}
	return chargeCode, nil
}
