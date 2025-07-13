package handler

import (
	"context"
	"darius/internal/constants"
	"darius/internal/errors"
	"darius/pkg/proto/suggest"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type SuggestInterviewQuestionParseFunc struct{}

func (p SuggestInterviewQuestionParseFunc) Parse(input string) (interface{}, error) {
	return convertToInterviewQuestionResponse(input)
}

func (h *handler) SuggestInterviewQuestion(ctx context.Context, req *suggest.SuggestInterviewQuestionRequest) (*suggest.SuggestInterviewQuestionResponse, error) {
	if req.GetContext() == nil || req.GetSubmissions() == nil {
		log.Println("[SuggestInterviewQuestion] context or submissions is nil")
		return nil, errors.Error(errors.ErrInvalidInput)
	}

	listOfPreviosQuestions := convertSuggestInterviewSubmissionToString(req.GetSubmissions())
	prompt := generateSuggestInterviewQuestionPrompt(req, listOfPreviosQuestions)

	parseFunc := SuggestInterviewQuestionParseFunc{}
	result, err := h.retryCallLLM(ctx, constants.F3_SCORE_INTERVIEW_QUESTIONS, prompt, parseFunc)
	if err != nil {
		return nil, err
	}

	if scoreResp, ok := result.(*suggest.SuggestInterviewQuestionResponse); ok {
		return scoreResp, nil
	}

	return nil, errors.Error(errors.ErrJSONParsing)
}

func convertToInterviewQuestionResponse(llmResponse string) (*suggest.SuggestInterviewQuestionResponse, error) {
	input := llmResponse

	jsonStr, err := extractAndSanitizeJSON(input)
	if err != nil {
		log.Println("[SuggestInterviewQuestion] error json parsing", err)
		return nil, errors.Error(errors.ErrJSONParsing)
	}

	// Parse JSON
	questionListResp, err := parseInterviewQuestions(jsonStr)
	if err != nil {
		log.Println("[SuggestInterviewQuestion] error json unmarshalling", err)
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
You are an expert in creating high-quality, contextually appropriate interview questions. Your task is to generate the next **two interview questions** based on the provided interview information and previous questions. To ensure the questions are pedagogically sound, logically structured, and role-appropriate, follow a chain-of-thought process with controlled output logic.

---

🧠 Chain-of-Thought Reasoning Process:
1. **Understand the Interview Context**: Review the field, position, language, and required skills.
2. **Avoid Repetition**: Analyze the previous questions to ensure novelty and progression.
3. **Align to Skills**: Ensure each question focuses on one or more of the provided skills.
4. **Controlled Structure**:
   - Question 1 must test **conceptual knowledge** or a **key theory** related to the role/field.
   - Question 2 must be a **realistic application** or **practical scenario** that allows the candidate to demonstrate applied understanding of the concept from Question 1.
5. **Language Matching**: Generate both questions in the specified language.
6. **Strict Quantity**: Only return 2 new questions. Do not duplicate or rephrase previous ones.

---

📥 Provided Input:
- Field: %v  
- Position: %v  
- Language: %v  
- Skills: %v  
- Max Questions Allowed: %v  
- Previous Questions: %v  

---

📤 Output Format:
Return only a **valid JSON object** using the following structure:

{
  "questions": [
    "First conceptual/theoretical question here",
    "Second real-world application question here"
  ]
}
🧪 Example (Few-Shot Prompting Guide):

📥 Input Example:

Field: Backend Development

Position: Junior Backend Developer

Language: English

Skills: REST API Design, HTTP Methods

Max Questions Allowed: 10

Previous Questions: ["What is the difference between PUT and POST in RESTful APIs?"]

📤 Output Example:
{
  "questions": [
    "Can you explain the core principles of REST architecture and how they guide API design?",
    "Imagine you're tasked with designing an API for a bookstore. How would you use REST principles to define endpoints for managing books and authors?"
  ]
}
⚠️ Constraints:

Do not generate vague, repetitive, or off-topic questions.

Ensure that the second question meaningfully builds on the first.

Maintain clarity, precision, and technical relevance.

Return valid JSON only. Do not include explanations, formatting, or markdown.

Now, generate the next two questions based on the input above.
	
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
