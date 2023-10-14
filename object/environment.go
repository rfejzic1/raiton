package object

type Environment struct {
	enclosing *Environment
	symbols   map[string]Object
}

func NewEnvironment() *Environment {
	return NewEnclosedEnvironment(nil)
}

func NewEnclosedEnvironment(env *Environment) *Environment {
	return &Environment{
		enclosing: env,
		symbols:   map[string]Object{},
	}
}

func (e *Environment) Define(name string, value Object) Object {
	e.symbols[name] = value
	return value
}

func (e *Environment) Lookup(name string) (Object, bool) {
	obj, ok := e.symbols[name]

	if !ok && e.enclosing != nil {
		return e.enclosing.Lookup(name)
	}

	return obj, ok
}

func (e *Environment) Enclosing() *Environment {
	return e.enclosing
}
