package domain

type App interface {
	Run()
}

type PaginationConfig struct {
	DefaultLimit int
	MaxLimit     int
}

func (config *PaginationConfig) ApplyPaginationConfig(p *Pagination) {
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
