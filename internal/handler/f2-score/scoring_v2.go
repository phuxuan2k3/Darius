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
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	prompt := generatePromptV2(data)
	llmResponse, err := h.llmManager.Generate(ctx, constants.F2_SCORE, prompt)
	if err != nil {
		log.Printf("Error generating response: %v", err)
		return
	}

	parsedResponse, err := sanitizeAndParseResponseV2(llmResponse)
	if err != nil {
		log.Printf("Error parsing response: %v", err)
		return
	}

	responseByte, err := proto.Marshal(parsedResponse)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		return
	}

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
		log.Printf("Error publishing message: %v", err)
	}
}

func generatePromptV2(data *ekko.EvaluationRequestV2) string {
	return fmt.Sprintf(
		`
		You are an expert AI tutor responsible for evaluating short-answer or essay-style responses. Your task is to grade a user's answer to a given question by comparing it with the correct answer. You will be given the question, the user's answer, the correct answer, and the maximum score for the question.

Your output must be in valid JSON format and must strictly follow the schema below:
{
  "score": number,          // An integer or decimal score between 0 and the maximum value provided in "points"
  "comment": "string"       // A short comment (1–3 sentences) explaining the strengths and weaknesses of the user's answer
}
  📝 Input Format:
You will receive input in this format:
{
  "questionText": "string",        // The question being asked
  "answer": "string",              // The user's submitted answer
  "correctAnswer": "string",       // The ideal answer for full credit
  "points": number,                // Maximum points for the question
  "x-user-id": "string",           // (metadata) user ID
  "x-role-id": "string",           // (metadata) role ID
  "timestamp": "string"            // (metadata) submission time
}
✅ Evaluation Guidelines:
Score from 0 to the maximum points value (e.g., 0–5), based on:

Relevance: How well the answer addresses the question.

Completeness: Whether the key parts of the correct answer are included.

Clarity: Whether the explanation is understandable and coherent.

Accuracy: Whether factual information is correct and matches the correct answer.

Be fair and consistent in applying partial credit. Don’t penalize for minor grammar issues if the answer is clear and correct.

If the answer is blank or irrelevant, assign 0.

The comment must be constructive and clear, offering brief reasoning for the score.

📌 Important Constraints:
Output must be in valid JSON with no trailing commas.

You must only return the JSON object and nothing else (no explanation, no markdown formatting).

Always return both score and comment.

🔧 Example:
{
  "questionText": "Explain the difference between TCP and UDP.",
  "answer": "TCP is connection-based and reliable. UDP is faster but doesn't ensure delivery.",
  "correctAnswer": "TCP is a connection-oriented protocol that guarantees data delivery in order. UDP is connectionless and does not guarantee delivery, making it faster but less reliable.",
  "points": 5,
  "x-user-id": "u123",
  "x-role-id": "student",
  "timestamp": "2025-06-14T12:00:00Z"
}
  {
  "score": 4,
  "comment": "The answer correctly contrasts TCP and UDP and captures the core differences. However, it lacks some technical detail about connection orientation and ordering."
}
  Now, evaluate the following input:

%v
		`, data)
}

func sanitizeAndParseResponseV2(input string) (*ekko.EvaluationResponseV2, error) {
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start == -1 || end == -1 || start > end {
		return nil, errors.New("no JSON object found in input")
	}
	jsonStr := input[start : end+1]

	var parsed ekko.EvaluationResponseV2
	err := json.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	return &parsed, nil
}
