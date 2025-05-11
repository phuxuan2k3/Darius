package handler

import (
	"context"
	llm "darius/internal/services/llm"
	"darius/models"
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
	prompt := fmt.Sprintf(`
You are an expert in creating test questions and answers. Your task is to generate a set of questions and answers based on the provided test information and guidelines. Follow these steps:
1. Input Provided by the User:
   - General Information:
    Name: %v,
	Description: %v,
	Fields: %v,
	Duration: %v,
	Difficulty: %v,
	QuestionType: %v,
	Language: %v,
	Options: %v,
	NumberOfQuestion: %v,
	CandidateSeniority: %v,
	Context: %v
2. Your Task:
   - Review the general information about the test to understand its context, purpose, and constraints.
   - Generate a set of questions and answers that align with the test's context, difficulty level, and format.
   - Ensure the questions are clear, precise, and meaningful.
   - Provide the output in the specified JSON format.
3. Output Format:
   [
     {
       questionContent: "[Question 1]",
       optionList: [
         {
           optionContent: "[Option 1]",
           isCorrect: [true/false]
         },
         {
           optionContent: "[Option 2]",
           isCorrect: [true/false]
         },
         ...
       ]
     },
     {
       questionContent: "[Question 2]",
       optionList: [
         {
           optionContent: "[Option 1]",
           isCorrect: [true/false]
         },
         {
           optionContent: "[Option 2]",
           isCorrect: [true/false]
         },
         ...
       ]
     },
     ...
   ]
Now, based on the user's input, generate the output in the specified format

	`, req.GetName(), req.GetDescription(), req.GetFields(), req.GetDuration(), req.GetDifficulty(), req.GetQuestionType(), req.GetLanguage(), req.GetOptions(), req.GetNumberOfQuestion(), req.GetCandidateSeniority(), req.GetContext())
	// llmResponse, err := h.llmService.Generate(ctx, &llm.LLMRequest{
	// 	Model:   viper.GetString("llm.model"),
	// 	Content: prompt,
	// })
	llmResponse, err := h.llmManager.Generate(ctx, models.F1, prompt)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	input := llmResponse

	jsonStr, err := extractJSONQuestions(input)
	if err != nil {
		fmt.Println("Lỗi:", err)
		return nil, err
	}

	// Parse JSON
	questionListResp, err := parseQuestions(jsonStr)
	if err != nil {
		fmt.Println("Lỗi:", err)
		return nil, err
	}

	return &suggest.SuggestQuestionsResponse{
		QuestionList: questionListResp,
	}, nil
}

func extractJSONQuestions(input string) (string, error) {
	re := regexp.MustCompile(`(?s)\[\s*\{.*\}\s*\]`)
	match := re.FindString(input)
	if match == "" {
		fmt.Println("extractJSONQuestions: Không tìm thấy JSON hợp lệ")
		return "", fmt.Errorf("không tìm thấy JSON hợp lệ")
	}

	match, err := sanitizeJSON(match)
	if err != nil {
		return "", fmt.Errorf("lỗi vệ sinh JSON: %v", err)
	}

	return match, nil
}

func parseQuestions(jsonStr string) ([]*suggest.Question, error) {
	var questions []*suggest.Question
	err := json.Unmarshal([]byte(jsonStr), &questions)
	if err != nil {
		return nil, fmt.Errorf("lỗi giải mã JSON: %v", err)
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
	return "", fmt.Errorf("chuỗi json chứa ký tự không thể vệ sinh")
}
