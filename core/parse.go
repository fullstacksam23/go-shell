package core

type ParsedCommand struct {
	Command        string
	Args           []string
	StdoutFile     string
	StdoutRedirect bool
}
