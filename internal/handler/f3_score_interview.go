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
		log.Print("[ScoreInterview] error parsing json")
		return nil, errors.Error(errors.ErrJSONParsing)
	}
	jsonStr := input[start : end+1]

	var parsed suggest.ScoreInterviewResponse
	err := proto.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		log.Printf("[ScoreInterview] error unmarshalling json: %v", err)
		return nil, fmt.Errorf("error unmarshalling escaped JSON: %v", err)
	}

	return &parsed, nil
}

func generateScoreInterviewPrompt(req *suggest.ScoreInterviewRequest) string {
	submissionByte, _ := proto.Marshal(req)
	submissionString := string(submissionByte)

	return fmt.Sprintf(`
		You are an expert interview evaluator. Your job is to evaluate an interview session of a candidate based on the provided Q&A data. Each submission contains a question and the candidate’s answer. You will analyze and assign a score with detailed feedback.

---

🧠 Step 1: Primary Evaluation  
For each submission:
1. Think step-by-step:
   - Is the answer relevant to the question?
   - Is it complete and accurate?
   - Is the explanation logically structured and clear?
2. Based on your reasoning, assign a grade:
   - A = Excellent
   - B = Good
   - C = Fair
   - D = Poor
   - F = Unacceptable
3. Write a **comment of 5–10 full sentences** explaining the score.

---

🧪 Skill Evaluation  
Evaluate each skill in the provided list by considering **all answers together**. Assign a grade (A–F) and justify the score in your mind (but only include score in the JSON).

---

📊 Summary  
Provide:
- "totalScore": How many times each letter was used (A–F)
- "positiveFeedback": A paraghaph about points of strengths
- "actionableFeedback": A paraghaph points for improvement
- "finalComment": A paragraph from 3–5 sentence overview of performance

---

🔁 Step 2: Self-Evaluation  
Check if your output is structured correctly and follows the guidelines.
Reflect on the evaluation. If any score or comment seems inconsistent or too harsh/generous, revise.  
---

📌 Output Format (strictly JSON):

{
  "result": [
    {
      "index": 1,
      "comment": "Full evaluation comment (5–10 sentences)",
      "score": "A"
    }
  ],
  "skills": [
    {
      "skill": "Problem Solving",
      "score": "B"
    }
  ],
  "totalScore": {
    "A": 1,
    "B": 0,
    "C": 0,
    "D": 0,
    "F": 0
  },
  "positiveFeedback": "Clearly explains reasoning. Demonstrates technical accuracy",
  "actionableFeedback": "Needs more concise examples. Could improve edge-case handling",
  "finalComment": "The candidate performed well overall, demonstrating confidence and clarity in their answers. Their explanations were logical and showed a good understanding of core concepts. With some refinement in structure and edge-case coverage, they would excel further.",
}
📚 Example Input and Output (Few-shot Prompting):

Example Input:
{
  "submissions": [
    {
      "question": "Explain the difference between an interface and an abstract class in OOP.",
      "answer": "An interface only has method declarations, and classes implement it. An abstract class can have method definitions and fields."
    }
  ],
  "skills": ["Object-Oriented Design"]
}
Example Output:

{
  "result": [
    {
      "index": 1,
      "comment": "The answer correctly identifies the distinction between interface and abstract class. It notes that interfaces contain declarations only, while abstract classes can include implementations. However, it could be more complete by mentioning multiple inheritance limitations or constructor availability. The structure is clear, and the terms are used correctly. Overall, a strong and concise response.",
      "score": "A"
    }
  ],
  "skills": [
    {
      "skill": "Object-Oriented Design",
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
  "positiveFeedback": "Clearly distinguishes concepts. Accurate and concise explanation",
  "actionableFeedback": "Could briefly mention inheritance limitations",
  "finalComment": "The candidate gave an accurate, high-level comparison of interface and abstract class. Their explanation was technically sound and easy to follow. This indicates a solid grasp of OOP fundamentals.",
}
Now evaluate the following interview session:
%v`, submissionString)
}
