package domain

type Pagination struct {
	Page  int
	Limit int
}

type PaginationConfig struct {
	DefaultLimit int
	MaxLimit     int
}

func ApplyPaginationConfig(p *Pagination, config PaginationConfig) {
	if p.Limit > config.MaxLimit {
		p.Limit = config.MaxLimit
	}

	if p.Limit < 0 {
		p.Limit = config.DefaultLimit
	}

	if p.Page <= 0 {
		p.Page = 1
	}
}
