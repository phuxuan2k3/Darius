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

		expected := "Question  1 : What is the difference between AI and Machine Learning?\nAnswer  1 : AI is a broader concept, while ML is a subset of AI.\nQuestion  2 : Explain the concept of overfitting in machine learning.\nAnswer  2 : Overfitting occurs when a model learns the training data too well, including noise and outliers.\n"
		result := convertSuggestInterviewSubmissionToString(submissions)

		assert.Equal(t, expected, result, "The converted string should match the expected output")
	})
}

func Test_promtGenerate(t *testing.T) {
	t.Run("Test promtGenerate", func(t *testing.T) {
		req := &suggest.SuggestInterviewQuestionRequest{
			Context: &suggest.SuggestInterviewQuestionRequest_Context{
				Field:        "Software Engineering",
				Position:     "Backend Developer",
				Language:     "Go",
				Level:        "Intermediate",
				MaxQuestions: 5,
			},
			Submissions: []*suggest.SuggestInterviewQuestionRequest_Submission{
				{
					Question: "What is the difference between AI and Machine Learning?",
					Answer:   "AI is a broader concept, while ML is a subset of AI.",
				},
			},
		}
		listOfPreviosQuestions := convertSuggestInterviewSubmissionToString(req.GetSubmissions())
		generateSuggestInterviewQuestionPrompt(req, listOfPreviosQuestions)
		expected := "\n\tYou are an expert in creating interview questions. Your task is to generate the next interview question (only 1 question) based on the provided interview information and guidelines. Follow these steps:\n\t1. Provided Input:\n\t   - General Information:\n\t\tField: Software Engineering,\n\t\tPosition: Backend Developer,\n\t\tLanguage: Go,\n\t\tLevel: Intermediate,\n\t\tMax Question: 5,\n\t\tList of previous questions: Question  1 : What is the difference between AI and Machine Learning?\nAnswer  1 : AI is a broader concept, while ML is a subset of AI.\n,\n\t2. Your Task:\n\t   - Review the general information about the test to understand its context, purpose, and constraints.\n\t   - Generate the next question for the interview that align with the interview's context, feild, language, position, difficulty level, and format.\n\t   - Ensure the questions are clear, precise, and meaningful.\n\t   - Provide the output in the specified JSON format.\n\t3. Output Format:\n\t\t   {\n\t\t\"question\": [\"The next question content here\"],\n\t\t}\n\tNow, based on the input, generate the output in the specified format\n\t\n\t\t"
		assert.Equal(t, expected, generateSuggestInterviewQuestionPrompt(req, listOfPreviosQuestions), expected)
	})

	t.Run("Test promtGenerate with Full submissions", func(t *testing.T) {
		req := &suggest.SuggestInterviewQuestionRequest{
			Context: &suggest.SuggestInterviewQuestionRequest_Context{
				Field:        "AI Engineering",
				Position:     "AI Engineer Intern",
				Language:     "English",
				Models:       "en-BG-male",
				Speed:        -5,
				Level:        "Easy",
				MaxQuestions: 10,
				SkipIntro:    false,
				Coding:       false,
				InterviewId:  "random-id",
			},
			Submissions: []*suggest.SuggestInterviewQuestionRequest_Submission{
				{
					Question: "Could you introduce yourself?",
					Answer:   "I'm John Doe, a fourth-year student in University of Science.",
				},
				{
					Question: "What are your strengths and weaknesses?",
					Answer:   "I'm good at coding and problem-solving, but I'm not good at communication.",
				},
				{
					Question: "What are your favourite projects?",
					Answer:   "I love the project that I built an AI chatbot for my school.",
				},
				{
					Question: "What are your experience with AI?",
					Answer:   "",
				},
			},
		}
		listOfPreviosQuestions := convertSuggestInterviewSubmissionToString(req.GetSubmissions())
		expected := "\n\tYou are an expert in creating interview questions. Your task is to generate the next two (only 2) interview questions based on the provided interview information and guidelines. Follow these steps:\n\t1. Provided Input:\n\t   - General Information:\n\t\tField: AI Engineering,\n\t\tPosition: AI Engineer Intern,\n\t\tLanguage: English,\n\t\tLevel: Easy,\n\t\tMax Question: 10,\n\t\tList of previous questions: Question  1 : Could you introduce yourself?\nAnswer  1 : I'm John Doe, a fourth-year student in University of Science.\nQuestion  2 : What are your strengths and weaknesses?\nAnswer  2 : I'm good at coding and problem-solving, but I'm not good at communication.\nQuestion  3 : What are your favourite projects?\nAnswer  3 : I love the project that I built an AI chatbot for my school.\nQuestion  4 : What are your experience with AI?\nAnswer  4 : \n,\n\t2. Your Task:\n\t   - Review the general information about the test to understand its context, purpose, and constraints.\n\t   - Generate the next two (2) questions for the interview that align with the interview's context, feild, language, position, difficulty level, and format.\n\t   - Ensure the questions are clear, precise, and meaningful.\n\t   - Provide the output in the specified JSON format.\n\t3. Output Format:\n\t\t   {\n\t\t\"question\": [\"The next question content here\"],\n\t\t}\n\tNow, based on the input, generate the output in the specified format\n\t\n\t\t"
		assert.Equal(t, expected, generateSuggestInterviewQuestionPrompt(req, listOfPreviosQuestions))
	})
}

func Test_convertToInterviewQuestionResponse(t *testing.T) {
	t.Run("Test convertToInterviewQuestionResponse", func(t *testing.T) {
		input := "{\n\t\"question\": [\n\t\t\"What programming languages or tools are you familiar with that are relevant to AI engineering?\",\n\t\t\"Can you explain a simple AI concept you learned in your studies or projects?\"\n\t]\n}"
		actual, _ := convertToInterviewQuestionResponse(input)
		expected := &suggest.SuggestInterviewQuestionResponse{}
		assert.Equal(t, expected, actual)
	})
}
