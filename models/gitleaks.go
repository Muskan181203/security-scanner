package models

type GitleaksFinding struct {
	RuleID      string `json:RuleID`
	Description string `json:Description`
	StartLine   int    `json:StartLine`
	File        string `json:File`
}
