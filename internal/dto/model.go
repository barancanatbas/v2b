package dto

type Module struct {
	Path       string  `json:"Path"`
	Version    string  `json:"Version"`
	CommitHash string  `json:"CommitHash"`
	Branch     *string `json:"Branch"`
}
