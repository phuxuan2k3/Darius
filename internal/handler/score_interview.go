package handler

import (
	"context"
	suggest "darius/pkg/proto/suggest"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"google.golang.org/protobuf/proto"
)

func (h *handler) ScoreInterview(ctx context.Context, req *suggest.ScoreInterviewRequest) (*suggest.ScoreInterviewResponse, error) {
	if len(req.GetSubmissions()) == 0 {
		log.Println("[ScoreInterview] submissions is nil")
		return nil, errors.New(" submissions is nil")
	}

	prompt := generateScoreInterviewPrompt(req)

	llmResponse, err := h.llmGRPCService.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}
	log.Println("[ScoreInterview] LLM response:", llmResponse)

	return sanitizeAndParseResponse(llmResponse)
}

func sanitizeAndParseResponse(input string) (*suggest.ScoreInterviewResponse, error) {
	// B1: Lấy JSON từ dấu { đầu tiên đến } cuối cùng. shit hello work
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
		`You are an expert interview evaluator. Your job is to evaluate an interview session of a candidate based on the provided Q&A data. Each submission contains a question asked during the interview and the corresponding answer from the candidate. Your evaluation should include:

For each submission:

Give a short comment (if necessary) on the quality of the answer.

Assign a score from the set {A, B, C, D, F}, where:

A = Excellent

B = Good

C = Fair

D = Poor

F = Unacceptable

Summary:

A totalScore object showing how many of each score (A–F) were assigned.

A positiveFeedback section: list the candidate’s notable strengths in bullet points.

An actionableFeedback section: list areas of improvement in bullet points.

A finalComment: 2–3 sentences summarizing the overall performance.

✅ Your response must strictly follow the JSON format below:

json
Copy
Edit
{
  "result": [
    {
      "index": 1,
      "comment": "Your comment here (optional, can be empty)",
      "score": "A"
    },
    {
      "index": 2,
      "comment": "",
      "score": "B"
    }
    // ...
  ],
  "totalScore": {
    "A": 1,
    "B": 6,
    "C": 2,
    "D": 1,
    "F": 0
  },
  "positiveFeedback": "- ...\n- ...",
  "actionableFeedback": "- ...\n- ...",
  "finalComment": "..."
}
Now, evaluate the following submissions:
%v`, submissionString)
}
