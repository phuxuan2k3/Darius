package f2_score

import (
	"context"
	"darius/internal/constants"
	"darius/pkg/proto/deps/ekko"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	proto "google.golang.org/protobuf/encoding/protojson"
)

func (h *scoringHandler) ScoreV2(ctx context.Context, req *ScoreRequest) {
	data := &ekko.EvaluationRequestV2{}
	err := proto.Unmarshal(req.Msg.Body, data)
	if err != nil {
		log.Printf("[ScoreV2] Error unmarshalling message: %v", err)
		return
	}

	prompt := generatePromptV2(data)
	_, _, llmResponse, err := h.llmManager.Generate(ctx, constants.F2_SCORE, prompt, nil)
	if err != nil {
		log.Printf("[ScoreV2] Error generating response: %v", err)
		return
	}

	parsedResponse, err := sanitizeAndParseResponseV2(llmResponse)
	if err != nil {
		log.Printf("[ScoreV2] Error parsing response: %v", err)
		return
	}

	responseByte, err := proto.Marshal(parsedResponse)
	if err != nil {
		log.Printf("[ScoreV2] Error marshalling response: %v", err)
		return
	}

	log.Printf("[ScoreV2] Successfully processed request with ID: %s.\n resp: %s ", req.Msg.MessageId, string(responseByte))

	err = h.queueChannel.Publish(
		"", h.queueQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        responseByte,
			MessageId:   req.Msg.MessageId,
			Timestamp:   req.Msg.Timestamp,
			Type:        "ScoreV2",
		},
	)
	if err != nil {
		log.Printf("[ScoreV2] Error publishing message: %v", err)
	}
}

func generatePromptV2(data *ekko.EvaluationRequestV2) string {
	return fmt.Sprintf(
		`
	You are an expert AI tutor responsible for grading short-answer and essay-style responses in standardized assessments. Your job is to evaluate a user's answer by comparing it to the ideal answer, using a transparent, fair, and detailed reasoning process. You must return the result in a strict JSON format as defined below.
---
📥 Input Format:
You will receive a JSON object with the following structure:
{
  "questionText": "string",
  "answer": "string",
  "correctAnswer": "string",
  "points": number,
  "x-user-id": "string",
  "timestamp": "string",
  "answerId": "string"
  "language": "string"
}
---
🧠 Evaluation Guidelines (Chain-of-Thought Reasoning):
1. **Relevance**: Does the response address the main idea of the question?
2. **Completeness**: Are all key points or required parts of the correct answer covered?
3. **Clarity**: Is the explanation coherent and understandable?
4. **Accuracy**: Are facts presented correctly and aligned with the correct answer?

Assign a score (int) between 0 and the maximum "points", applying partial credit where appropriate. Minor grammar mistakes should not be penalized unless they impact understanding.

---

✍️ Comment Requirements:
- Your evaluation **must** include a "comment" field that is between **3 to 5 full sentences**.
- The comment must explain both the **strengths** and **weaknesses** of the answer.
- Use clear, constructive, and specific language that helps the candidate understand their performance.

---

📤 Output Format (Strictly Required):
{
  "score": number,
  "comment": "string (3–5 full sentences)",
  "timestamp": "string" (keep the same as input),
  "answerId": "string (must keep the same as input)",

}

---

📚 Few-shot Examples:

Example 1:
Input:
{
  "questionText": "What is the capital of France?",
  "answer": "Paris",
  "correctAnswer": "The capital of France is Paris.",
  "points": 5,
  "x-user-id": "u01",
  "timestamp": "2025-06-14T09:30:00Z"
  "answerId": "a01",
  "language": "English"
}
Output:
{
  "score": 5,
  "comment": "Your answer is short but completely correct. It identifies the capital of France accurately and directly. While brief, it leaves no room for confusion. Well done."
  "timestamp": "2025-06-14T09:30:00Z",
  "answerId": "a01"
}

Example 2:
Input:
{
  "questionText": "Explain the difference between TCP and UDP.",
  "answer": "TCP is slower than UDP.",
  "correctAnswer": "TCP is a connection-oriented protocol that guarantees delivery and order. UDP is connectionless and faster but doesn't guarantee delivery.",
  "points": 10,
  "x-user-id": "u02",
  "timestamp": "2025-06-14T09:35:00Z",
  "answerId": "a02",
  "language": "Vietnamese"
}
Output:
{
  "score": 3,
  "comment": "Câu trả lời của bạn cho thấy nhận thức cơ bản về sự khác biệt về hiệu suất giữa TCP và UDP. Tuy nhiên, nó thiếu các chi tiết kỹ thuật quan trọng như định hướng kết nối, đảm bảo phân phối và thứ tự. Tuyên bố này quá mơ hồ và có thể gây hiểu lầm trong bối cảnh kỹ thuật. Hãy cân nhắc giải thích chi tiết về hành vi cốt lõi của từng giao thức. Điều này sẽ chứng minh sự hiểu biết sâu sắc hơn về các nguyên tắc cơ bản của mạng.",
  "timestamp": "2025-06-14T09:35:00Z",
  "answerId": "a02"
}

Example 3:
Input:
{
  "questionText": "Define polymorphism in object-oriented programming.",
  "answer": "It means functions can do different things.",
  "correctAnswer": "Polymorphism allows objects of different classes to be treated through the same interface, enabling methods to behave differently based on the object instance.",
  "points": 3,
  "x-user-id": "u03",
  "timestamp": "2025-06-14T09:40:00Z",
  "answerId": "a03",
  "language": "English"
}
Output:
{
  "score": 1,
  "comment": "Your answer shows some understanding of the core concept behind polymorphism. However, it is too vague and lacks technical accuracy. You did not mention the use of interfaces or the behavior of methods in different object contexts. With more precise language and an example, your answer would be much stronger. Try expanding your definition in future responses.",
  "timestamp": "2025-06-14T09:40:00Z",
  "answerId": "a03"
  }

---

Now, based on the following input, return your evaluation:
%v

		`, data)
}

func sanitizeAndParseResponseV2(input string) (*ekko.EvaluationResponseV2, error) {
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start == -1 || end == -1 || start > end {
		return nil, errors.New("[ScoreV2] no JSON object found in input")
	}
	jsonStr := input[start : end+1]

	var parsed ekko.EvaluationResponseV2
	err := json.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		return nil, fmt.Errorf("[ScoreV2] error unmarshalling JSON: %v", err)
	}
	return &parsed, nil
}
