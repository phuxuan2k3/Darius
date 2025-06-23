package handler

import (
	"context"
	"darius/internal/errors"
	"darius/models"
	"darius/pkg/proto/suggest"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

func (h *handler) SuggestInterviewQuestion(ctx context.Context, req *suggest.SuggestInterviewQuestionRequest) (*suggest.SuggestInterviewQuestionResponse, error) {
	if req.GetContext() == nil || req.GetSubmissions() == nil {
		log.Println("[SuggestInterviewQuestion] context or submissions is nil")
		return nil, errors.Error(errors.ErrInvalidInput)
	}

	listOfPreviosQuestions := convertSuggestInterviewSubmissionToString(req.GetSubmissions())
	prompt := generateSuggestInterviewQuestionPrompt(req, listOfPreviosQuestions)

	llmResponse, err := h.llmManager.Generate(ctx, models.F3_SUGGEST_INTERVIEW_QUESTIONS, prompt)
	if err != nil {
		return nil, errors.Error(errors.ErrNetworkConnection)
	}
	log.Println("[SuggestInterviewQuestion] LLM response:", llmResponse)

	return convertToInterviewQuestionResponse(llmResponse)
}

func convertToInterviewQuestionResponse(llmResponse string) (*suggest.SuggestInterviewQuestionResponse, error) {
	input := llmResponse

	jsonStr, err := extractAndSanitizeJSON(input)
	if err != nil {
		log.Println("Lỗi:", err)
		return nil, errors.Error(errors.ErrJSONParsing)
	}

	// Parse JSON
	questionListResp, err := parseInterviewQuestions(jsonStr)
	if err != nil {
		log.Println("Lỗi:", err)
		return nil, errors.Error(errors.ErrJSONUnmarshalling)
	}

	return questionListResp, nil
}

func extractAndSanitizeJSON(input string) (string, error) {
	// Bước 1: Chuẩn hóa escape sequences
	normalized := strings.ReplaceAll(input, `\n`, "\n")
	normalized = strings.ReplaceAll(normalized, `\t`, "\t")

	// Bước 2: Regex bắt JSON object hoặc array (không dùng đệ quy)
	re := regexp.MustCompile(`(?s)(\{.*?\}|\[.*?\])`)
	match := re.FindString(normalized)
	if match == "" {
		return "", errors.Error(errors.ErrJSONParsing)
	}

	// Bước 3: Kiểm tra tính hợp lệ
	if !json.Valid([]byte(match)) {
		return "", errors.Error(errors.ErrJSONParsing)
	}

	return match, nil
}

func parseInterviewQuestions(jsonStr string) (*suggest.SuggestInterviewQuestionResponse, error) {
	var questions *suggest.SuggestInterviewQuestionResponse
	err := json.Unmarshal([]byte(jsonStr), &questions)
	if err != nil {
		return nil, errors.Error(errors.ErrJSONUnmarshalling)
	}
	return questions, nil
}

func generateSuggestInterviewQuestionPrompt(req *suggest.SuggestInterviewQuestionRequest, listOfPreviosQuestions string) string {
	return fmt.Sprintf(`
	You are an expert in creating interview questions. Your task is to generate the next two (only 2) interview questions based on the provided interview information and guidelines. Follow these steps:
	1. Provided Input:
	   - General Information:
		Field: %v,
		Position: %v,
		Language: %v,
		Skills: %v,

		Max Question: %v,
		List of previous questions: %v,
	2. Your Task:
	   - Review the general information about the test to understand its context, purpose, and constraints.
	   - Generate the next two (2) questions for the interview that align with the interview's context, feild, language, position, difficulty level, and format.
	   - Ensure the questions are clear, precise, and meaningful.
	   - Provide the output in the specified JSON format.
	3. Output Format:
		   {
		"questions": ["The next question content here"],
		}
	Now, based on the input, generate the output in the specified format
	
		`, req.GetContext().GetPosition(), req.GetContext().GetExperience(), req.GetContext().GetLanguage(), req.GetContext().GetSkills(), req.GetContext().GetMaxQuestions(), listOfPreviosQuestions)
}
func convertSuggestInterviewSubmissionToString(submissions []*suggest.SuggestInterviewQuestionRequest_Submission) string {
	listOfPreviosQuestions := ""
	for index, submission := range submissions {
		listOfPreviosQuestions += fmt.Sprintln("Question ", index+1, ":", submission.GetQuestion())
		listOfPreviosQuestions += fmt.Sprintln("Answer ", index+1, ":", submission.GetAnswer())
	}
	return listOfPreviosQuestions
}
