package handler

import (
	"context"
	"darius/internal/constants"
	"darius/internal/errors"
	suggest "darius/pkg/proto/suggest"
	"fmt"
	"log"
	"strings"

	proto "google.golang.org/protobuf/encoding/protojson"
)

func (h *handler) ScoreInterview(ctx context.Context, req *suggest.ScoreInterviewRequest) (*suggest.ScoreInterviewResponse, error) {
	if len(req.GetSubmissions()) == 0 {
		log.Println("[ScoreInterview] submissions is nil")
		return nil, errors.Error(errors.ErrInvalidInput)
	}

	prompt := generateScoreInterviewPrompt(req)

	llmResponse, err := h.llmManager.Generate(ctx, constants.F3_SCORE_INTERVIEW_QUESTIONS, prompt)
	if err != nil {
		return nil, errors.Error(errors.ErrNetworkConnection)
	}

	return sanitizeAndParseResponse(llmResponse)
}

func sanitizeAndParseResponse(input string) (*suggest.ScoreInterviewResponse, error) {
	// B1: Lấy JSON từ dấu { đầu tiên đến } cuối cùng. shit hello work workr
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start == -1 || end == -1 || start > end {
		return nil, errors.Error(errors.ErrJSONParsing)
	}
	jsonStr := input[start : end+1]

	var parsed suggest.ScoreInterviewResponse
	err := proto.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling escaped JSON: %v", err)
	}

	return &parsed, nil
}

func generateScoreInterviewPrompt(req *suggest.ScoreInterviewRequest) string {
	submissionByte, _ := proto.Marshal(req)
	submissionString := string(submissionByte)

	return fmt.Sprintf(
		`You are an expert interview evaluator. Your job is to evaluate an interview session of a candidate based on the provided Q&A data. Each submission contains a question asked during the interview and the corresponding answer from the candidate.  

Your evaluation must include the following sections and strictly follow the provided JSON structure.

---

For each submission in the interview:
- Assign a score from the following set:  
  - A = Excellent  
  - B = Good  
  - C = Fair  
  - D = Poor  
  - F = Unacceptable  
- Provide a short comment about the answer’s quality (optional — leave empty if not needed).

Skill evaluation:
- You are also given a list of skills to evaluate.
- For **each skill provided**, you must assign a score (A–F) that reflects the candidate’s overall demonstrated ability in that skill based on all the submissions.
- ✅ **Always include the 'skills' field** in your response, even if there's only one submission.

Summary section must include:
- 'totalScore': A count of how many times each score (A, B, C, D, F) was given across all submissions.
- 'positiveFeedback': A bullet-point list of notable strengths.
- 'actionableFeedback': A bullet-point list of areas the candidate needs to improve.
- 'finalComment': A 2–3 sentence summary of the candidate's overall performance.

✅ Your response must strictly follow this JSON structure:
{
  "result": [
    {
      "index": 1,
      "comment": "Your comment here (optional, can be empty)",
      "score": "A"
    }
  ],
  "skills": [
    {
      "skill": "Accuracy",
      "score": "A"
    }
  ],
  "totalScore": {
    "A": 1,
    "B": 0,
    "C": 0,
    "D": 0,
    "F": 0
  },
  "positiveFeedback": "- ...\n- ...",
  "actionableFeedback": "- ...\n- ...",
  "finalComment": "..."
}

Now, evaluate the following submissions:

%v`, submissionString)
}
