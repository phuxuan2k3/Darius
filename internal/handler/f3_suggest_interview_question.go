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

func (h *handler) SuggestInterviewQuestion(ctx context.Context, req *suggest.SuggestInterviewQuestionRequest) (*suggest.SuggestInterviewQuestionResponse, error) {
	if req.GetContext() == nil || req.GetSubmissions() == nil {
		log.Println("[SuggestInterviewQuestion] context or submissions is nil")
		return nil, errors.Error(errors.ErrInvalidInput)
	}

	listOfPreviosQuestions := convertSuggestInterviewSubmissionToString(req.GetSubmissions())
	prompt := generateSuggestInterviewQuestionPrompt(req, listOfPreviosQuestions)

	llmResponse, err := h.llmManager.Generate(ctx, constants.F3_SUGGEST_INTERVIEW_QUESTIONS, prompt)
	if err != nil {
		return nil, errors.Error(errors.ErrNetworkConnection)
	}

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
You are an expert in creating high-quality, contextually appropriate interview questions. Your task is to generate the next **two interview questions** based on the provided interview information and previous questions. To ensure the generated questions are aligned and meaningful, follow a chain-of-thought reasoning process.

---

🧠 Chain of Thought Reasoning Instructions:
1. **Understand the Context**: Carefully examine the interview’s field, role, language, and skill focus to determine the purpose and tone of the interview.
2. **Analyze Past Questions**: Review the list of previous questions to avoid repetition and ensure progressive difficulty and topic coverage.
3. **Skill Matching**: Ensure the generated questions evaluate relevant skills from the provided skill set.
4. **Question Quality**: Ensure each question is clear, concise, and suitable for the role and level.
5. **Language**: The questions must be written in the specified language.
6. **Limits**: Only generate two new questions. Do not repeat or modify old ones.

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
Return only a **valid JSON object** as specified below. Do not include explanations or formatting outside this JSON.

{
  "questions": [
    "First generated question here",
    "Second generated question here"
  ]
}
⚠️ Constraints:
Do not include questions that are vague, overly generic, or redundant.
Avoid overlapping with previous questions.
Ensure diversity in format, phrasing, and focus across the two questions.
The output must be strictly valid JSON — no markdown, no comments, no explanations.
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
