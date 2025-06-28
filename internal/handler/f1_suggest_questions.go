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

func (h *handler) SuggestQuestions(ctx context.Context, req *suggest.SuggestQuestionsRequest) (*suggest.SuggestQuestionsResponse, error) {
	chargeCode, err := h.checkCanCall(ctx, constants.F1_SUGGEST_QUESTIONS)
	if err != nil {
		return nil, err
	}

	questionsContents, err := h.missfortune.GetExamQuestionContent(ctx, converters.ConvertSuggestQuestionRequestToMissFortuneRequest(ctx, req))
	prompt := ""
	if err != nil {
		log.Printf("[SuggestExamQuestion] error getting exam question content: %v", err)

		prompt = fmt.Sprintf(`
		You are an AI that generates multiple-choice questions based on provided metadata.
		Given the following input, generate a list of questions in strict JSON format.
		-Input:
			Title: %v;
			Description: %v;
			Minutes to answer: %v;
			Language: %v;
			Difficulty: %v;
			Tags: %v;
			Outlines: %v;
			Number Of Questions: %v;
			Number Of Options: %v;
		
		
		-Output format (strictly follow this structure):
		{
			questions: {
				text: string;
				options: string[];
				points: number; // positive
				correctOption: number; // index of correct option, starting from 0
			}[];
		}
		
		Requirements:
		The number of questions and options must match numberOfQuestions and numberOfOptions respectively.
		All questions must relate to the provided title, description, tags, and outlines.
		The questions should be appropriate for the given difficulty level.
		All options must be plausible, but only one is correct (correctOption).
		points should be a positive integer (e.g., 1 to 10) assigned to each question based on relevance and depth.
		Ensure the final result is valid JSON and strictly follows the output structure.
			`, req.GetTitle(), req.GetDescription(), req.GetMinutesToAnswer(), req.GetLanguage(), req.GetDifficulty(), req.GetTags(), req.GetOutlines(), req.GetNumberOfQuestions(), req.GetNumberOfOptions())
	} else {

		prompt = fmt.Sprintf(`
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
	}
	llmResponse, err := h.llmManager.Generate(ctx, constants.F1_SUGGEST_QUESTIONS, prompt)
	if err != nil {
		return nil, errors.Error(errors.ErrNetworkConnection)
	}

	input := llmResponse

	jsonStr, err := extractJSONQuestions(input)
	if err != nil {
		fmt.Println("Lỗi:", err)
		return nil, errors.Error(errors.ErrJSONParsing)
	}

	// Parse JSON
	questionListResp, err := parseQuestions(jsonStr)
	if err != nil {
		fmt.Println("Lỗi:", err)
		return nil, errors.Error(errors.ErrJSONUnmarshalling)
	}

	if !h.bulbasaur.ChargeCallingLLM(ctx, chargeCode) {
		log.Printf("[SuggestExamQuestion] Charge Code %s failed to charge for LLM call", chargeCode)
		return nil, errors.Error(errors.ErrChargingFailed)
	}

	return questionListResp, nil
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

func parseQuestions(jsonStr string) (*suggest.SuggestQuestionsResponse, error) {
	var questions *suggest.SuggestQuestionsResponse
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

	fmt.Println("sanitizeJSON: Chuỗi JSON chứa ký tự không thể vệ sinh")
	return "", errors.Error(errors.ErrJSONParsing)
}
