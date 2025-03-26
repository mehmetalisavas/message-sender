package models

type ListOptions struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Page   int `json:"page"`
}

var DefaultListOptions = ListOptions{
	Limit:  20,
	Offset: 0,
	Page:   1,
}

func InitWithDefaultListOptions(opts ListOptions) ListOptions {
	if opts.Limit <= 0 {
		opts.Limit = DefaultListOptions.Limit
	}
	if opts.Offset <= 0 {
		opts.Offset = DefaultListOptions.Offset
	}
	if opts.Page <= 0 {
		opts.Page = DefaultListOptions.Page
	}

	if opts.Page > 1 {
		opts.Offset = (opts.Page - 1) * opts.Limit
	}

	return opts
}
