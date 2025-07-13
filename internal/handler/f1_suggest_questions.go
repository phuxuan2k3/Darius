package handler

import (
	"context"
	"darius/internal/constants"
	"darius/internal/converters"
	"darius/internal/errors"
	llm "darius/internal/services/llm"
	"darius/pkg/proto/suggest"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"
)

func generateOptionsPrompt(questionsContent interface{}) string {
	return fmt.Sprintf(`
	You are an expert in designing high-quality standardized multiple-choice exam content. Your task is to generate **exactly 4 answer options** for each of the provided questions, identifying one correct answer and ensuring all options are meaningful and unique.
	---
	📥 Input Format:
	You will receive a JSON object with the following structure:
	{
	  "questions": [
		"What is the capital of France?",
		"Which data structure uses LIFO order?",
		...
	  ]
	}
	
	---
	
	🧠 Thought Process (Chain-of-Thought Required):
	1. **Understand the concept** behind each question and identify the accurate correct answer.
	2. **Generate three incorrect but plausible options** that are:
	   - Contextually related to the topic.
	   - Grammatically and semantically consistent with the question.
	   - Distinct from the correct answer and from each other.
	3. Before finalizing each set of options:
	   - Check for **semantic duplication** (e.g., “Rome” vs. “Roma”) or repetition (e.g., two identical options).
	   - Validate that no two options have the same meaning, phrasing, or structure.
	4. Ensure the correct answer is placed at a random index (from 0 to 3), and record that index in the "correctOption" field.
	5. Use the "points" field to reflect question difficulty (e.g., Easy = 1–3, Medium = 4–6, Hard = 7–10).
	
	---
	
	📤 Output Format:
	Respond with a valid **strict JSON object** that adheres to the following Protobuf-compatible schema:
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
	
	Exactly 4 answer options per question.
	Only one correct option, indicated by the correctOption index (0–3).
	Incorrect options must be plausible but clearly incorrect.
	No two options may be textually or semantically identical.
	Options must be diverse in content and style, not simple variations.
	All options must be grammatically compatible with the question.
	The order of questions in the output must match the input.
	Output must be valid JSON without comments, explanations, or formatting errors.
	
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
		  "points": 3,
		  "correctOption": 0
		}
	  ]
	}
	  🧪 Validation Tip:
	Before returning the result, double-check:
	Each question has exactly one correct answer.
	There is no overlap in the options.
	The format is strictly valid JSON.
	Now, based on the following input, generate the answer options:
	%v
		`, questionsContent)
}

func (h *handler) SuggestQuestions(ctx context.Context, req *suggest.SuggestQuestionsRequest) (*suggest.SuggestExamQuestionResponseV2, error) {
	chargeCode, err := h.checkCanCall(ctx, constants.F1_SUGGEST_QUESTIONS)
	if err != nil {
		return nil, err
	}

	questionsContents, err := h.missfortune.GetExamQuestionContent(ctx, converters.ConvertSuggestQuestionRequestToMissFortuneRequest(ctx, req))
	prompt := ""
	if err != nil {
		log.Printf("[SuggestQuestions] error getting exam question content: %v", err)

		prompt = fmt.Sprintf(`
	You are an expert exam question designer. Your task is to generate a diverse, non-redundant, high-quality set of exam questions based on the structured input below.

📥 Input Metadata:
%v

🧠 Generation Process (Chain-of-Thought Required):
1. Carefully analyze the title, description, tags, and outlines to understand the full context and intended coverage.
2. Brainstorm a diverse pool of possible questions (both MCQ and LONG_ANSWER) aligned with the specified difficulty level and key topics.
3. Filter questions to ensure:
   - ✅ Questions language must matches the input language.
   - ✅ The number of questions matches the requested count.
   - ✅ The number of MCQ and LONG_ANSWER questions must be balanced.
   - ❌ No repeated or semantically similar questions.
   - ✅ Broad and representative coverage of all outlines and tags.
4. For MCQ:
   - Generate **exactly 4 unique options**.
   - Ensure only **1 correct answer** is clearly identifiable (correctOption index: 0–3).
   - Ensure all options are **plausible**, **grammatically consistent**, and **non-overlapping in meaning**.
5. For LONG_ANSWER:
   - Provide at least one illustrative image link if applicable.
   - Include clear instructions and an ideal sample answer.
6. ✅ Final validation steps:
   - Questions language must match the input language.
   - The number of MCQ and LONG_ANSWER questions must be balanced.
   - No question or option duplication.
   - All JSON fields are present and conform strictly to the expected schema.
   - Output is a valid JSON object (no markdown, no explanations).

📤 Output Format:
Return a single **valid JSON object** with the following format:
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
📌 Constraints Recap:

Questions language must match the input language.
The number of MCQ and LONG_ANSWER questions must be balanced.
Return exactly the number of questions.
Each MCQ has exactly 4 unique options.
No repeated questions or options allowed.
The points field should reflect difficulty level (e.g., Easy = 1–3, Medium = 4–6, Hard = 7–10).
All outputs must be raw JSON without any explanations or comments.
Now, generate the questions based on the provided metadata.
			`, req)
	} else {
		prompt = generateOptionsPrompt(questionsContents)
	}
	_, llmResponse, err := h.llmManager.Generate(ctx, constants.F1_SUGGEST_QUESTIONS, prompt, nil)
	if err != nil {
		return nil, errors.Error(errors.ErrNetworkConnection)
	}

	input := llmResponse

	jsonStr, err := extractJSONQuestions(input)
	if err != nil {
		fmt.Println("[SuggestQuestions] error extract Json questions:", err)
		return nil, errors.Error(errors.ErrJSONParsing)
	}

	// Parse JSON
	questionListResp, err := parseQuestions(jsonStr)
	if err != nil {
		fmt.Println("[SuggestQuestions] error parse questions", err)
		return nil, errors.Error(errors.ErrJSONUnmarshalling)
	}

	if !h.bulbasaur.ChargeCallingLLM(ctx, chargeCode) {
		log.Printf("[SuggestQuestions] Charge Code %s failed to charge for LLM call", chargeCode)
		return nil, errors.Error(errors.ErrChargingFailed)
	}

	return questionListResp, nil
}

func (h *handler) SuggestOptions(ctx context.Context, req *suggest.SuggestOptionsRequest) (*suggest.SuggestOptionsResponse, error) {

	return &suggest.SuggestOptionsResponse{
		CriteriaList: &suggest.CriteriaEleResponse{
			Criteria:   "criteria1",
			OptionList: []string{"option1", "option2", "option3"},
		},
	}, nil
}

func (h *handler) SuggestCriteria(ctx context.Context, req *suggest.SuggestCriteriaRequest) (*suggest.SuggestCriteriaResponse, error) {
	generalInfo := req.GetGeneralInfo()
	if generalInfo == nil {
		log.Println("generalInfo is nil")
		return nil, nil
	}

	criteriaList := req.GetCriteriaList()
	if criteriaList == nil {
		return &suggest.SuggestCriteriaResponse{
			CriteriaList: []*suggest.CriteriaEleResponse{
				{
					Criteria: "Test Subject Area",
					OptionList: []string{
						"Computer Networks",
						"Hardware",
						"Software Development"},
				},
				{
					Criteria:   "Difficulty Level:",
					OptionList: []string{"Beginner", "Intermediate", "Advanced"},
				},
				{
					Criteria:   "Test Format:",
					OptionList: []string{"Multiple Choice", "True/False", "Essay"},
				},
				{
					Criteria:   "Test Duration:",
					OptionList: []string{"30 minutes", "60 minutes", "90 minutes"},
				},
			},
		}, nil
	}

	prompt := fmt.Sprintf(`
You are an expert in designing tests and assessments. Your task is to analyze the provided input, which includes general information about the test and a list of criteria with the user's chosen options. Based on this input, suggest additional criteria and options that will help the user provide more detailed information for generating test questions. Follow these steps:
1. Input Provided by the User:
   - General Information:
     %v
   - Criteria List:
     %v
2. Your Task:
   - Review the general information about the test to understand its context, purpose, and constraints.
   - Analyze the list of criteria and the user's chosen options for each criterion.
   - Suggest additional criteria and options that will help the user provide more detailed information for generating test questions. Ensure the suggestions are relevant to the test's context and align with the user's chosen options.
   - Provide the output in the specified JSON format.

3. Output Format:
   [
     {
       criteria: '[Suggested Criterion 1]',
       optionList: [
         "[Suggested Option 1 for Criterion 1]",
         "[Suggested Option 2 for Criterion 1]",
       ]
     },
     {
       criteria: '[Suggested Criterion 2]',
       optionList: [
         "[Suggested Option 1 for Criterion 2]",
         "[Suggested Option 2 for Criterion 2]",
       ]
     },
     ...
   ]
Now, based on the user's input, generate the output in the specified format.
`, generalInfo, criteriaList)

	llmResponse, err := h.llmService.Generate(ctx, &llm.LLMRequest{
		Content: prompt,
	})

	input := llmResponse.Content

	jsonStr, err := extractJSONQuestions(input)
	if err != nil {
		fmt.Println("Lỗi:", err)
		return nil, err
	}

	// Parse JSON
	criteriaResp, err := parseCriterias(jsonStr)
	if err != nil {
		fmt.Println("Lỗi:", err)
		return nil, err
	}

	return &suggest.SuggestCriteriaResponse{
		CriteriaList: criteriaResp,
	}, nil
}

func extractJSONQuestions(input string) (string, error) {
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start == -1 || end == -1 || start > end {
		return "", errors.Error(errors.ErrJSONParsing)
	}
	jsonStr := input[start : end+1]

	return jsonStr, nil
}

func parseQuestions(jsonStr string) (*suggest.SuggestExamQuestionResponseV2, error) {
	var questions *suggest.SuggestExamQuestionResponseV2
	err := json.Unmarshal([]byte(jsonStr), &questions)
	if err != nil {
		return nil, errors.Error(errors.ErrJSONUnmarshalling)
	}
	return questions, nil
}

func parseCriterias(jsonStr string) ([]*suggest.CriteriaEleResponse, error) {
	var criterias []*suggest.CriteriaEleResponse
	err := json.Unmarshal([]byte(jsonStr), &criterias)
	if err != nil {
		fmt.Println("parseCriterias: Lỗi giải mã JSON:", err)
		return nil, fmt.Errorf("lỗi giải mã JSON: %v", err)
	}
	return criterias, nil
}

func sanitizeJSON(jsonStr string) (string, error) {
	reComment := regexp.MustCompile(`(?m)^\s*//.*$`)
	cleaned := reComment.ReplaceAllString(jsonStr, "")

	var builder strings.Builder
	for _, r := range cleaned {
		if r < 0x20 && r != '\n' && r != '\r' && r != '\t' {
			continue
		}
		if !unicode.IsPrint(r) && r != '\n' && r != '\r' && r != '\t' {
			continue
		}
		builder.WriteRune(r)
	}
	sanitized := builder.String()

	if json.Valid([]byte(sanitized)) {
		return sanitized, nil
	}

	start := strings.IndexAny(sanitized, "{[")
	end := strings.LastIndexAny(sanitized, "}]")
	if start != -1 && end != -1 && end > start {
		candidate := sanitized[start : end+1]
		if json.Valid([]byte(candidate)) {
			return candidate, nil
		}
	}

	fmt.Println("[SuggestQuestions] sanitizeJSON: Chuỗi JSON chứa ký tự không thể vệ sinh")
	return "", errors.Error(errors.ErrJSONParsing)
}
