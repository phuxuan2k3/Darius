package handler

import (
	"context"
	"darius/models"
	suggest "darius/pkg/proto/suggest"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	proto "google.golang.org/protobuf/encoding/protojson"
)

func (h *handler) ScoreInterview(ctx context.Context, req *suggest.ScoreInterviewRequest) (*suggest.ScoreInterviewResponse, error) {
	if len(req.GetSubmissions()) == 0 {
		log.Println("[ScoreInterview] submissions is nil")
		return nil, errors.New(" submissions is nil")
	}

	prompt := generateScoreInterviewPrompt(req)

	llmResponse, err := h.llmManager.Generate(ctx, models.F3, prompt)
	if err != nil {
		return nil, err
	}
	log.Println("[ScoreInterview] LLM response:", llmResponse)

	return sanitizeAndParseResponse(llmResponse)
}

func sanitizeAndParseResponse(input string) (*suggest.ScoreInterviewResponse, error) {
	// B1: Lấy JSON từ dấu { đầu tiên đến } cuối cùng. shit hello work workr
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start == -1 || end == -1 || start > end {
		return nil, errors.New("no JSON object found in input")
	}
	jsonStr := input[start : end+1]

	// B2: Escape các xuống dòng bên trong chuỗi JSON
	lines := strings.Split(jsonStr, "\n")
	var escapedLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Nếu dòng nằm trong một chuỗi JSON (tức là có dấu ":"), thì cần escape
		if strings.Contains(trimmed, ":") && strings.Count(trimmed, "\"") >= 2 {
			escapedLines = append(escapedLines, strings.ReplaceAll(line, "\n", "\\n"))
		} else {
			escapedLines = append(escapedLines, line)
		}
	}
	escapedJSON := strings.Join(escapedLines, "")

	// B3: Unmarshal
	var parsed suggest.ScoreInterviewResponse
	err := json.Unmarshal([]byte(escapedJSON), &parsed)
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
