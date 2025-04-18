package f2_score

import (
	"context"
	llm_grpc "darius/internal/llm-grpc"
	ekko "darius/pkg/proto/deps/ekko"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	proto "google.golang.org/protobuf/encoding/protojson"
)

type ScoringHandler interface {
	Score(ctx context.Context, req *ScoreRequest)
}

type scoringHandler struct {
	llmGRPCService llm_grpc.Service
	queueChannel   *amqp.Channel
	queueQueue     *amqp.Queue
}

func NewScoringHandler(llmGRPCService llm_grpc.Service, queueChannel *amqp.Channel, queueQueue *amqp.Queue) ScoringHandler {
	return &scoringHandler{
		llmGRPCService: llmGRPCService,
		queueChannel:   queueChannel,
		queueQueue:     queueQueue,
	}
}
func (h *scoringHandler) Score(ctx context.Context, req *ScoreRequest) {
	data := &ekko.EvaluationRequest{}
	err := proto.Unmarshal(req.Msg.Body, data)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	prompt := generatePrompt(data)
	llmResponse, err := h.llmGRPCService.Generate(ctx, prompt)
	if err != nil {
		log.Printf("Error generating response: %v", err)
		return
	}
	log.Printf("LLM response: %s", llmResponse)
	parsedResponse, err := sanitizeAndParseResponse(llmResponse)
	if err != nil {
		log.Printf("Error parsing response: %v", err)
		return
	}

	responseByte, err := proto.Marshal(parsedResponse)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		return
	}

	err = h.queueChannel.Publish("", h.queueQueue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        responseByte,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func sanitizeAndParseResponse(input string) (*ekko.EvaluationResponse, error) {
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start == -1 || end == -1 || start > end {
		return nil, errors.New("no JSON object found in input")
	}
	jsonStr := input[start : end+1]

	var parsed ekko.EvaluationResponse
	err := json.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	return &parsed, nil
}

func generatePrompt(req *ekko.EvaluationRequest) string {
	submissionByte, _ := proto.Marshal(req)
	submissionString := string(submissionByte)
	return fmt.Sprintf(`
You are an evaluation AI. Given a scenario description and a list of questions with user answers and criteria, your task is to evaluate each answer based on the following criteria:

Relevance (relevance): How well the answer relates to the question and the expected content based on the criteria.

Clarity & Completeness (clarityCompleteness): How clearly the answer is written and whether it is complete.

Accuracy (accuracy): How factually correct the answer is in the context of the question and criteria.

For each answer, provide:

A score from 1 to 10 (1 being very poor, 10 being excellent) for each criterion above.

An overall score from 1 to 10 (can be the average or your weighted judgment).

A status field: return "SUBMISSION_STATUS_SUCCESS" if the overall score is greater than or equal to 6, otherwise return "SUBMISSION_STATUS_FAILED".

Return the result in the following JSON format:
{
  "result": [
    {
      "id": <uint64 question_id>,
      "relevance": <float64 1-10>,
      "clarityCompleteness": <float64 1-10>,
      "accuracy": <float64 1-10>,
      "overall": <float64 1-10>,
      "status": "SUBMISSION_STATUS_SUCCESS" or "SUBMISSION_STATUS_FAILED"
    },
    ...
  ]
}
Here is the input data to evaluate:
%v
`, submissionString)
}
