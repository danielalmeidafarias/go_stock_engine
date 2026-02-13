package entities

type CriticalityLevel int

const (
	Low CriticalityLevel = iota + 1
	Moderate
	High
	VeryHigh
	Critical
)
