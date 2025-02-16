package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (env *Environment) Get(name string) (Object, bool) {
	obj, ok := env.store[name]
	if !ok && env.outer != nil {
		obj, ok = env.outer.Get(name)
	}
	return obj, ok
}

func (env *Environment) Set(name string, object Object) Object {
	env.store[name] = object
	return object
}

func NewEnvirnoment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvirnoment()
	env.outer = outer
	return env
}
