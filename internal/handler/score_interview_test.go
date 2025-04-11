package handler

import (
	"darius/pkg/proto/suggest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sanitizeAndParseResponse(t *testing.T) {
	t.Run("Test sanitizeAndParseResponse", func(t *testing.T) {
		input := `{
  "result": [
    {
      "index": 1,
      "comment": "The introduction is clear but very brief. Including more details such as major, interests, or relevant experience would provide a fuller picture.",
      "score": "C"
    },
    {
      "index": 2,
      "comment": "The answer identifies basic strengths and a weakness, but lacks specificity and depth. It would benefit from examples or elaboration.",
      "score": "C"
    },
    {
      "index": 3,
      "comment": "The answer gives a glimpse into a project of interest, but is vague. Describing the project's purpose, technologies used, or the impact would strengthen it.",
      "score": "C"
    }
  ],
  "totalScore": {
    "A": 0,
    "B": 0,
    "C": 3,
    "D": 0,
    "F": 0
  },
  "positiveFeedback": "- Demonstrates awareness of strengths and weaknesses.\n- Shows interest in meaningful projects like building an AI chatbot.",
  "actionableFeedback": "- Expand answers with specific examples, technical details, and outcomes.\n- Provide a more comprehensive self-introduction, including background and goals.\n- Use clear structure and elaboration to strengthen communication.",
  "finalComment": "The candidate shows potential and a genuine interest in technology, but responses are too brief and lack sufficient detail. With more preparation and elaboration, the candidate could present a much stronger impression."
}

`
		expected := &suggest.ScoreInterviewResponse{
			Result: []*suggest.ScoreInterviewResponse_Submission{
				{
					Index:   1,
					Comment: "The introduction is clear but very brief. Including more details such as major, interests, or relevant experience would provide a fuller picture.",
					Score:   "C",
				},
				{
					Index:   2,
					Comment: "The answer identifies basic strengths and a weakness, but lacks specificity and depth. It would benefit from examples or elaboration.",
					Score:   "C",
				},
				{
					Index:   3,
					Comment: "The answer gives a glimpse into a project of interest, but is vague. Describing the project's purpose, technologies used, or the impact would strengthen it.",
					Score:   "C",
				},
			},
			TotalScore: map[string]int32{
				"A": 0,
				"B": 0,
				"C": 3,
				"D": 0,
				"F": 0,
			},
			PositiveFeedback:   "- Demonstrates awareness of strengths and weaknesses.\n- Shows interest in meaningful projects like building an AI chatbot.",
			ActionableFeedback: "- Expand answers with specific examples, technical details, and outcomes.\n- Provide a more comprehensive self-introduction, including background and goals.\n- Use clear structure and elaboration to strengthen communication.",
			FinalComment:       "The candidate shows potential and a genuine interest in technology, but responses are too brief and lack sufficient detail. With more preparation and elaboration, the candidate could present a much stronger impression.",
		}

		result, err := sanitizeAndParseResponse(input)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		assert.Equal(t, expected, result)
	})
}
