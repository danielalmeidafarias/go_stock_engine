package domain

type Pagination struct {
	Page  int
	Limit int
}

type PaginationConfig struct {
	DefaultLimit int
	MaxLimit     int
}

func ApplyPaginationRules(p *Pagination, config PaginationConfig) {
	if p.Limit > config.MaxLimit {
		p.Limit = config.MaxLimit
	}

	if p.Limit <= 0 {
		p.Limit = config.DefaultLimit
	}

	if p.Page <= 0 {
		p.Page = 1
	}
}

func PaginatedSlice[T any](slice []T, p *Pagination) []T {
	offset := (p.Page - 1) * p.Limit

	if offset >= len(slice) {
		return []T{}
	}

	end := min(offset+p.Limit, len(slice))

	return slice[offset:end]
}
