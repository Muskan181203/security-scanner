package models

type SemgrepResult struct {
	CheckID string   `json:"check_id"`
	Path    string   `json:"path"`
	Start   Position `json:"start"`
	Extra   Extra    `json:"extra"`
}

type Position struct {
	Line int `json:"line"`
	Col  int `json:"col"`
}
type Extra struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type SemgrepResponse struct {
	Results []SemgrepResult `json:"results"`
}
