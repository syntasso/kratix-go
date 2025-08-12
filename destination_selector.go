package kratix

type DestinationSelector struct {
	Directory   string            `json:"directory"`
	MatchLabels map[string]string `json:"matchLabels"`
}
