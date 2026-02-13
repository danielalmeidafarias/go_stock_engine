package domain

type PaginationConfig struct {
	DefaultLimit int
	MaxLimit     int
}

type AppConfig struct {
	PaginationConfig PaginationConfig
}

func (config *AppConfig) ApplyPaginationConfig(p *Pagination) {
	if p.Limit > config.PaginationConfig.MaxLimit {
		p.Limit = config.PaginationConfig.MaxLimit
	}

	if p.Limit <= 0 {
		p.Limit = config.PaginationConfig.DefaultLimit
	}

	if p.Page <= 0 {
		p.Page = 1
	}
}
