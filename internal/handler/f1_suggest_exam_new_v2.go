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

		templateRule := "Generate %v questions of level %v for topic %v.\n"
		questionCount := 0
		instruction := ""
		for _, topic := range req.Topics {
			if topic.GetDifficultyDistribution().GetIntern() > 0 {
				instruction += fmt.Sprintf(templateRule, topic.GetDifficultyDistribution().GetIntern(), "Intern", topic.GetName())
				questionCount += int(topic.GetDifficultyDistribution().GetIntern())
			}
			if topic.GetDifficultyDistribution().GetJunior() > 0 {
				instruction += fmt.Sprintf(templateRule, topic.GetDifficultyDistribution().GetJunior(), "Junior", topic.GetName())
				questionCount += int(topic.GetDifficultyDistribution().GetJunior())
			}
			if topic.GetDifficultyDistribution().GetMiddle() > 0 {
				instruction += fmt.Sprintf(templateRule, topic.GetDifficultyDistribution().GetMiddle(), "Middle", topic.GetName())
				questionCount += int(topic.GetDifficultyDistribution().GetMiddle())
			}
			if topic.GetDifficultyDistribution().GetSenior() > 0 {
				instruction += fmt.Sprintf(templateRule, topic.GetDifficultyDistribution().GetSenior(), "Senior", topic.GetName())
				questionCount += int(topic.GetDifficultyDistribution().GetSenior())
			}
			if topic.GetDifficultyDistribution().GetLead() > 0 {
				instruction += fmt.Sprintf(templateRule, topic.GetDifficultyDistribution().GetLead(), "Lead", topic.GetName())
				questionCount += int(topic.GetDifficultyDistribution().GetLead())
			}
			if topic.GetDifficultyDistribution().GetExpert() > 0 {
				instruction += fmt.Sprintf(templateRule, topic.GetDifficultyDistribution().GetExpert(), "Expert", topic.GetName())
				questionCount += int(topic.GetDifficultyDistribution().GetExpert())
			}
		}

		instruction += fmt.Sprintf("Total questions you MUST generate: %d.\n", questionCount)
		req.Topics = nil // Clear topics to avoid duplication in the prompt
		prompt = fmt.Sprintf(`
You are an expert exam question designer. Your task is to generate a diverse set of high-quality exam questions based on the structured input below. Each question must be either a multiple-choice question (MCQ) or a long-answer (essay-style) question.
Here is your requirement:
%v
Before generating, follow this step-by-step reasoning to ensure conceptual diversity, difficulty alignment, and uniqueness:

---

🧠 Step-by-Step Reasoning (Chain-of-Thought Inspired):
1. Parse the input to extract all topics and their difficulty distributions.
2. For each topic:
   - For each difficulty level (Intern to Expert), generate the required number of questions.
   - The total number of questions per topic must match the sum of all difficulty levels defined in its "difficultyDistribution".
3. For each question:
   - Ensure it is relevant to the topic and appropriate for the intended difficulty level.
   - Vary the question types (MCQ vs. Long Answer) where reasonable.
4. For MCQ-type questions:
   - Create exactly 4 answer options per question.
   - Ensure all 4 options are:
     - Semantically distinct.
     - Grammatically aligned with the question.
     - Plausible (but only one is correct).
     - Free from duplication or close paraphrases.
5. For Long Answer questions:
   - Ensure the question requires analysis, comparison, or in-depth reasoning.
   - Provide a complete, detailed model answer.
   - Optionally include one or more imageLinks and/or guiding instructions via "extraText".

---

🔁 Self-Verification (Post-Generation Check):
- Verify that **no question text is repeated or similar across topics/difficulties**.
- For MCQs:
  - Ensure no two options are identical or overly similar.
- Confirm the number of questions exactly matches the total count derived from all difficulty distributions across topics.
- Validate JSON structure before returning.

---
📤 Output Format:
Return a valid JSON object with an array of questions. Each question must strictly follow this schema:

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

The type must be either "MCQ" or "LONG_ANSWER".

The detail.type must match the parent type.

For MCQ questions:

Exactly 4 options.

Only one correct option (indexed 0–3).

All options must be grammatically correct and distinct in both meaning and phrasing.

For Long Answer:

Must include correctAnswer.

Must target higher-order thinking skills (e.g., explanation, comparison).

imageLinks can be empty or populated.

The points field must reflect the difficulty:

Intern = 1–2

Junior = 2–3

Middle = 4–5

Senior = 6–7

Lead = 8–9

Expert = 10

The total number of questions generated per topic must match the sum of all difficulty levels in its difficultyDistribution.

Output must be clean, valid JSON with no markdown, explanations, or comments.

Now, generate the questions based on the following input:
%v
	`, instruction, req)
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
