package app

type App struct {
	basename    string
	name        string
	description string
	runFunc     RunFunc
	silence     bool
	noVersion   bool
	noConfig    bool
}

type Option func(*App)

func NewApp(name string, basename string, opts ...Option) *App {
	a := &App{
		name:     name,
		basename: basename,
	}

	for _, o := range opts {
		o(a)
	}

	return a
}

func WithSilence() Option {
	return func(a *App) {
		a.silence = true
	}
}

// WithNoVersion set the application does not provide version flag.
func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

// WithNoConfig set the application does not provide config flag.
func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

type RunFunc func(basename string) error

func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}


