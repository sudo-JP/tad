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
	Behavior string
	Payload  string
}

type BannedWord struct {
	Word     string
	Command  string
	Behavior string
	Payload  string
}

type LeetcodeOwed struct {
	UserID    string
	Command   string
	OwedCount int
}
