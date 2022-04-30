package mysql

import "time"

const (
	defaultMaxOpenConns    = 150
	defaultMaxIdleConns    = 150
	defaultConnMaxLifetime = 100
)

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

// WithTablePrefix new tablePrefix option
func WithTablePrefix(tablePrefix string) ConfigOption {
	return newConfigFuncOption(func(o *configOption) {
		o.tablePrefix = tablePrefix
	})
}

// WithSingularTable new singularTable option
func WithSingularTable(singularTable bool) ConfigOption {
	return newConfigFuncOption(func(o *configOption) {
		o.singularTable = singularTable
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
	tablePrefix     string
	singularTable   bool
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
		maxOpenConns:    defaultMaxOpenConns,
		maxIdleConns:    defaultMaxIdleConns,
		connMaxLifetime: defaultConnMaxLifetime,
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
