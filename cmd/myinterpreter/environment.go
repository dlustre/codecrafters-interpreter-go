package main

type Environment struct {
	// The outer scope, or nil if this is the global environment.
	Enclosing *Environment
	Values    map[string]any
}

func (e *Environment) get(name Token) (any, error) {
	// fmt.Printf("getting %s in current scope: (%v)\n", name.Lexeme, e)
	if value, ok := e.Values[name.Lexeme]; ok {
		if value == nil {
			return "nil", nil
		}
		return value, nil
	}

	if e.Enclosing != nil {
		// fmt.Printf("getting %s in enclosing scope: (%v)\n", name.Lexeme, e.Enclosing)
		return e.Enclosing.get(name)
	}

	// fmt.Println("could not find " + name.Lexeme)
	return nil, RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."}
}

func (e *Environment) assign(name Token, value any) error {
	// fmt.Println("assigning " + name.Lexeme)
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.assign(name, value)
	}

	// fmt.Println("could not find " + name.Lexeme)
	return RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."}
}

func (e *Environment) define(name string, value any) {
	// fmt.Printf("defining %s as: %v\n", name, value)
	e.Values[name] = value
	// for k, v := range e.Values {
	// 	fmt.Printf("%s -> %v\n", k, v)
	// }
}
