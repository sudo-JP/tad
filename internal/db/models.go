package db

type User struct {
	UserID   string
	Username string
}

type SudoUser struct {
	UserID string
}

type Command struct {
	Command  string
	Behavior CommandBehavior
	Spam     *SpamCommand
	Media    *MediaCommand
}

type CommandBehavior string

const (
	BehaviorSpam  CommandBehavior = "spam"
	BehaviorMedia CommandBehavior = "media"
)

type SpamCommand struct {
	Command   string
	Payload   string
	SpamCount *int
}

type MediaCommand struct {
	Command string
	Payload string
}

type LeetcodeOwed struct {
	UserID    string
	OwedCount int
}

// LeetcodeProblemList is the "questionList" result from LeetCode's GraphQL API.
type LeetcodeProblemList struct {
	Total     int               `json:"total"`
	Questions []LeetcodeProblem `json:"questions"`
}

type LeetcodeProblem struct {
	Number     string `json:"frontendQuestionId"`
	Title      string `json:"title"`
	Slug       string `json:"titleSlug"`
	Difficulty string `json:"difficulty"` // "Easy", "Medium" or "Hard"
	PaidOnly   bool   `json:"isPaidOnly"`
}
