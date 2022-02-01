package application

type Engine struct{}

func Build() (*Engine, error) {
	return &Engine{}, nil
}
