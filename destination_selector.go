package kratix

type DestinationSelectorModifier interface{}

type DestinationSelector struct {
	Directory   string
	MatchLabels map[string]any
}
