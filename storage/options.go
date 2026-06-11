package storage

type Options struct {
	ContentType string
	Public      bool
}

type Option func(*Options)

func WithContentType(ct string) Option {
	return func(o *Options) {
		o.ContentType = ct
	}
}

func WithPublic() Option {
	return func(o *Options) {
		o.Public = true
	}
}
