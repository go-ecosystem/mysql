package mysql

import "time"

// WithMaxOpenConns new maxOpenConns option
func WithMaxOpenConns(maxOpenConns int) ConfigOption {
	return newConfigFuncOption(func(o *configOption) {
		o.maxOpenConns = maxOpenConns
	})
}

// WithMaxIdleConns new maxIdleConns option
func WithMaxIdleConns(maxIdleConns int) ConfigOption {
	return newConfigFuncOption(func(o *configOption) {
		o.maxIdleConns = maxIdleConns
	})
}

// WithConnMaxLifetime new connMaxLifetime option
func WithConnMaxLifetime(connMaxLifetime time.Duration) ConfigOption {
	return newConfigFuncOption(func(o *configOption) {
		o.connMaxLifetime = connMaxLifetime
	})
}

// ConfigOption config option
type ConfigOption interface {
	apply(*configOption)
}

type configOption struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
}

func decodeConfigOpts(opts []ConfigOption) configOption {
	op := defaultConfigOption()
	for _, opt := range opts {
		opt.apply(&op)
	}
	return op
}

func defaultConfigOption() configOption {
	return configOption{
		maxOpenConns:    150,
		maxIdleConns:    150,
		connMaxLifetime: 100,
	}
}

type configFuncOption struct {
	f func(*configOption)
}

func (fo *configFuncOption) apply(do *configOption) {
	fo.f(do)
}

func newConfigFuncOption(f func(*configOption)) *configFuncOption {
	return &configFuncOption{
		f: f,
	}
}
