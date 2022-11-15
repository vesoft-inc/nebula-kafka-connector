package logger

var defaultOptions = options{
	level:   InfoLevel,
	console: true,
}

type (
	options struct {
		level   Level
		fields  Fields
		console bool
		files   []string
	}
	Option func(*options)
)

func WithLevel(lvl Level) Option {
	return func(o *options) {
		o.level = lvl
	}
}

func WithLevelText(text string) Option {
	return func(o *options) {
		WithLevel(ParseLevel(text))(o)
	}
}

func WithFields(fields ...Field) Option {
	return func(o *options) {
		o.fields = fields
	}
}

func WithConsole(console bool) Option {
	return func(o *options) {
		o.console = console
	}
}

func WithFiles(files ...string) Option {
	return func(o *options) {
		o.files = files
	}
}
