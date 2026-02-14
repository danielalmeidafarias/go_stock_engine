package entities

type ProductCategory string

const (
	Engine ProductCategory = "engine"
	Oil    ProductCategory = "oil"
)

func IsValidProductCategory(c ProductCategory) bool {
	switch c {
	case Engine, Oil:
		return true
	default:
		return false
	}
}
