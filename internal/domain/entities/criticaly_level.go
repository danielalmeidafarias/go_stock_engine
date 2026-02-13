package entities

type CriticalityLevel int

const (
	Low CriticalityLevel = iota + 1
	Moderate
	High
	VeryHigh
	Critical
)

func IsValidCriticalityLevel(c CriticalityLevel) bool {
	switch c {
	case Low, Moderate, High, VeryHigh, Critical:
		return true
	default:
		return false
	}
}
