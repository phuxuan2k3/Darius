package handler

import (
	"context"
	"darius/pkg/proto/suggest"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func (h *handler) SuggestInterviewQuestion(ctx context.Context, req *suggest.SuggestInterviewQuestionRequest) (*suggest.SuggestInterviewQuestionResponse, error) {
	if req.GetContext() == nil || req.GetSubmissions() == nil {
		log.Println("[SuggestInterviewQuestion] context or submissions is nil")
		return nil, errors.New("context or submissions is nil")
	}

	listOfPreviosQuestions := convertSubmissionToString(req.GetSubmissions())
	prompt := promtGenerate(req, listOfPreviosQuestions)
	
	llmResponse, err := h.llmGRPCService.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}
	log.Println("[SuggestInterviewQuestion] LLM response:", llmResponse)

	input := llmResponse

	jsonStr, err := extractJSONQuestions(input)
	if err != nil {
		log.Println("Lỗi:", err)
		return nil, err
	}

	// Parse JSON
	questionListResp, err := parseInterviewQuestions(jsonStr)
	if err != nil {
		log.Println("Lỗi:", err)
		return nil, err
	}

	return questionListResp, nil
}

func promtGenerate(req *suggest.SuggestInterviewQuestionRequest, listOfPreviosQuestions string) string {
	return fmt.Sprintf(`
	You are an expert in creating interview questions. Your task is to generate the next interview question (only 1 question) based on the provided interview information and guidelines. Follow these steps:
	1. Provided Input:
	   - General Information:
		Field: %v,
		Position: %v,
		Language: %v,
		Level: %v,
		Max Question: %v,
		List of previous questions: %v,
	2. Your Task:
	   - Review the general information about the test to understand its context, purpose, and constraints.
	   - Generate the next question for the interview that align with the interview's context, feild, language, position, difficulty level, and format.
	   - Ensure the questions are clear, precise, and meaningful.
	   - Provide the output in the specified JSON format.
	3. Output Format:
		   {
		"question": ["The next question content here"],
		}
	Now, based on the input, generate the output in the specified format
	
		`, req.GetContext().GetField(), req.GetContext().GetPosition(), req.GetContext().GetLanguage(), req.GetContext().GetLevel(), req.GetContext().GetMaxQuestions(), listOfPreviosQuestions)
}	

func convertSubmissionToString(submissions []*suggest.SuggestInterviewQuestionRequest_Submission) string {
	listOfPreviosQuestions := ""
	for index, submission := range submissions {
		listOfPreviosQuestions += fmt.Sprintln("Question ", index+1,":", submission.GetQuestion())
		listOfPreviosQuestions += fmt.Sprintln("Answer ", index+1,":", submission.GetAnswer())
	}
	return listOfPreviosQuestions
}

func parseInterviewQuestions(jsonStr string) (*suggest.SuggestInterviewQuestionResponse, error) {
	var questions *suggest.SuggestInterviewQuestionResponse
	err := json.Unmarshal([]byte(jsonStr), &questions)
	if err != nil {
		return nil, fmt.Errorf("lỗi giải mã JSON: %v", err)
	}
	return questions, nil
}