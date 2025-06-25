package constants

type LLMCallAmount struct {
	Amount float32 `json:"amount"`
	Desc   string  `json:"desc"`
}

var AmountMap = map[string]LLMCallAmount{
	F1_SUGGEST_EXAM:                {Amount: 500, Desc: "F1 Suggest Exam"},
	F1_SUGGEST_OUTLINES:            {Amount: 500, Desc: "F1 Suggest Outlines"},
	F1_SUGGEST_QUESTIONS:           {Amount: 500, Desc: "F1 Suggest Questions"},
	F2_SCORE:                       {Amount: 0, Desc: "F2 Score"},
	F3_SUGGEST_INTERVIEW_QUESTIONS: {Amount: 0, Desc: "F3 Suggest Interview Questions"},
	F3_SCORE_INTERVIEW_QUESTIONS:   {Amount: 700, Desc: "F3 Score Interview Questions"},
}

func GetLLMCallAmount(key string) (float32, string) {
	if amount, exists := AmountMap[key]; exists {
		return amount.Amount, amount.Desc
	}
	return 0, "Unknown LLM call"
}
