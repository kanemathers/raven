package main

type Module interface {
	Init(*IRCClient) error
}

var (
	modules = make(map[string]func() Module)
)

func RegisterModule(name string, f func() Module) {
	modules[name] = f
}

func LoadModule(name string) Module {
	if module, ok := modules[name]; ok {
		return module()
	}

	return nil
}
