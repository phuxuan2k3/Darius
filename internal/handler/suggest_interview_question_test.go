package handler

import (
	"darius/pkg/proto/suggest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convertSubmissionToString(t *testing.T) {
	t.Run("Test convertSubmissionToString", func(t *testing.T) {
		submissions := []*suggest.SuggestInterviewQuestionRequest_Submission{
			{
				Question: "What is the difference between AI and Machine Learning?",
				Answer:   "AI is a broader concept, while ML is a subset of AI.",
			},
			{
				Question: "Explain the concept of overfitting in machine learning.",
				Answer:   "Overfitting occurs when a model learns the training data too well, including noise and outliers.",
			},
		}

		expected :=  "Question  1 : What is the difference between AI and Machine Learning?\nAnswer  1 : AI is a broader concept, while ML is a subset of AI.\nQuestion  2 : Explain the concept of overfitting in machine learning.\nAnswer  2 : Overfitting occurs when a model learns the training data too well, including noise and outliers.\n"
		result := convertSubmissionToString(submissions)

		assert.Equal(t, expected, result, "The converted string should match the expected output")
	},)
}

func Test_promtGenerate(t *testing.T) {
	t.Run("Test promtGenerate", func(t *testing.T) {
		req := &suggest.SuggestInterviewQuestionRequest{
			Context: &suggest.SuggestInterviewQuestionRequest_Context{
				Field: 		   "Software Engineering",
				Position: 	   "Backend Developer",
				Language: 	   "Go",
				Level: 		   "Intermediate",
				MaxQuestions:   5,
			},
			Submissions: []*suggest.SuggestInterviewQuestionRequest_Submission{
				{
					Question: "What is the difference between AI and Machine Learning?",
					Answer:   "AI is a broader concept, while ML is a subset of AI.",
				},
			},
		}
		listOfPreviosQuestions := convertSubmissionToString(req.GetSubmissions())
		promtGenerate(req, listOfPreviosQuestions)
		expected := "\n\tYou are an expert in creating interview questions. Your task is to generate the next interview question (only 1 question) based on the provided interview information and guidelines. Follow these steps:\n\t1. Provided Input:\n\t   - General Information:\n\t\tField: Software Engineering,\n\t\tPosition: Backend Developer,\n\t\tLanguage: Go,\n\t\tLevel: Intermediate,\n\t\tMax Question: 5,\n\t\tList of previous questions: Question  1 : What is the difference between AI and Machine Learning?\nAnswer  1 : AI is a broader concept, while ML is a subset of AI.\n,\n\t2. Your Task:\n\t   - Review the general information about the test to understand its context, purpose, and constraints.\n\t   - Generate the next question for the interview that align with the interview's context, feild, language, position, difficulty level, and format.\n\t   - Ensure the questions are clear, precise, and meaningful.\n\t   - Provide the output in the specified JSON format.\n\t3. Output Format:\n\t\t   {\n\t\t\"question\": [\"The next question content here\"],\n\t\t}\n\tNow, based on the input, generate the output in the specified format\n\t\n\t\t"
		assert.Equal(t, expected, promtGenerate(req, listOfPreviosQuestions), expected)
	},)}