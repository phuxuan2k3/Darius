package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type LLM interface {
	GenerateCriteria(ctx context.Context, req *LLMRequest) (*LLMResponse, error)
	GenerateQuestion(ctx context.Context, req *LLMRequest) (*LLMResponse, error)
	Generate(ctx context.Context, req *LLMRequest) (*LLMResponse, error)
}

type llmInstance struct {
	config *Config
}

func NewLLM(config *Config) *llmInstance {
	if config == nil {
		config = &Config{
			Host:  viper.GetString("llm.host"),
			Model: viper.GetString("llm.model"),
		}
	}
	return &llmInstance{
		config: config,
	}
}

func (l *llmInstance) Generate(ctx context.Context, req *LLMRequest) (*LLMResponse, error) {
	url := l.config.Host
	req.Model = l.config.Model
	data, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling request:", err)
		return nil, err
	}

	llmReq, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	llmReq.Header.Set("Content-Type", "application/json")

	// Gửi request
	client := &http.Client{}
	llmResp, err := client.Do(llmReq)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer llmResp.Body.Close()

	// Đọc response
	fmt.Println("response body:", llmResp.Body)
	var response LLMResponse
	if err := json.NewDecoder(llmResp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding response:", err)
		return nil, err
	}

	log.Println("response", response)

	return &response, nil
}

func (l *llmInstance) GenerateCriteria(ctx context.Context, req *LLMRequest) (*LLMResponse, error) {

	//xtodo : implement the logic to generate LLM

	return &LLMResponse{
		Created_at: "",
		Content: fmt.Sprintf(`
				[
		  {
		    "criteria": "Question Type",
		    "optionList": [
		      "Short Answer",
		      "True or False",
		      "Fill in the Blank",
		      "Matching",
		      "Word Problems"
		    ]
		  },
		  {
		    "criteria": "Topics Covered",
		    "optionList": [
		      "Geometry",
		      "Trigonometry",
		      "Statistics and Probability",
		      "Number Theory",
		      "Functions and Graphs"
		    ]
		  }
		]
				`),
	}, nil
}
func (l *llmInstance) GenerateQuestion(ctx context.Context, req *LLMRequest) (*LLMResponse, error) {

	//xtodo : implement the logic to generate LLM

	return &LLMResponse{
		Content: fmt.Sprintf(`
		[
    {
        "questionContent": "What is the solution to the equation 2x + 3 = 7?",
        "optionList": [
            {
                "optionContent": "x = 2",
                "isCorrect": true
            },
            {
                "optionContent": "x = 3",
                "isCorrect": false
            },
            {
                "optionContent": "x = 4",
                "isCorrect": false
            },
            {
                "optionContent": "x = 5",
                "isCorrect": false
            }
        ]
    },
    {
        "questionContent": "Which of the following is a quadratic equation?",
        "optionList": [
            {
                "optionContent": "x + 2 = 5",
                "isCorrect": false
            },
            {
                "optionContent": "x² + 3x + 2 = 0",
                "isCorrect": true
            },
            {
                "optionContent": "2x + 3 = 7",
                "isCorrect": false
            },
            {
                "optionContent": "3x³ + 2x² + x = 0",
                "isCorrect": false
            }
        ]
    },
    {
        "questionContent": "What is the slope of the line y = 3x + 2?",
        "optionList": [
            {
                "optionContent": "2",
                "isCorrect": false
            },
            {
                "optionContent": "3",
                "isCorrect": true
            },
            {
                "optionContent": "1",
                "isCorrect": false
            },
            {
                "optionContent": "0",
                "isCorrect": false
            }
        ]
    }
]
		`),
	}, nil
}
