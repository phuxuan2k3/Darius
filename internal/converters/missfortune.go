package converters

import (
	"context"
	"darius/pkg/proto/deps/missfortune"
	"darius/pkg/proto/suggest"
	"fmt"
)

func ConvertExamRequestToMissfortuneRequest(ctx context.Context, req *suggest.SuggestExamQuestionRequest) *missfortune.SuggestExamQuestionRequest {
	return &missfortune.SuggestExamQuestionRequest{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Language:    req.GetLanguage(),
		Seniority:   req.GetSeniority(),
		Topics:      convertExamRequestTopicsToMissfortuneRequestTopic(req.GetTopics()),

		Creativity: req.GetCreativity(),
		Context: &missfortune.SuggestExamQuestionRequest_Context{
			Text:  "a",
			Links: req.GetContext().GetLinks(),
		},
		QuestionType: "Multiple Choice",
	}
}

func ConvertSuggestQuestionRequestToMissFortuneRequest(ctx context.Context, req *suggest.SuggestQuestionsRequest) *missfortune.SuggestExamQuestionRequest {
	return &missfortune.SuggestExamQuestionRequest{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Language:    req.GetLanguage(),
		Seniority:   req.GetDifficulty(),
		Context: &missfortune.SuggestExamQuestionRequest_Context{
			Text: fmt.Sprintf(
				`
				Minutes to answer: %d,
				Tags : %v,
				Outlines: %v,
				Number of questions: %d,
				Number of options: %d,
				`, req.GetMinutesToAnswer(), req.GetTags(), req.GetOutlines(), req.GetNumberOfQuestions(), req.GetNumberOfOptions()),
		},
		QuestionType: "Multiple Choice",
	}
}

func convertExamRequestTopicsToMissfortuneRequestTopic(topics []*suggest.Topic) []*missfortune.Topic {
	mfTopics := make([]*missfortune.Topic, len(topics))
	for i, topic := range topics {
		mfTopics[i] = &missfortune.Topic{
			Name: topic.GetName(),
			DifficultyDistribution: &missfortune.DifficultyDistribution{
				Intern: topic.GetDifficultyDistribution().GetIntern(),
				Junior: topic.GetDifficultyDistribution().GetJunior(),
				Middle: topic.GetDifficultyDistribution().GetMiddle(),
				Senior: topic.GetDifficultyDistribution().GetSenior(),
				Lead:   topic.GetDifficultyDistribution().GetLead(),
			},
		}
	}
	return mfTopics
}
