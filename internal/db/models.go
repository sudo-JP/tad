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
