package entities

type ProductCategory string

const (
	Engine ProductCategory = "engine"
)

func IsValidProductCategory(c ProductCategory) bool {
	switch c {
	case Engine:
		return true
	default:
		return false
	}
}
