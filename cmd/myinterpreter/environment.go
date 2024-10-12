package main

type Environment struct {
	Values map[string]any
}

func (e *Environment) get(name Token) (any, error) {
	// fmt.Println("getting " + name.Lexeme)
	if value, ok := e.Values[name.Lexeme]; ok {
		if value == nil {
			return "nil", nil
		}
		return value, nil
	}
	// fmt.Println("could not find " + name.Lexeme)
	return nil, RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."}
}

func (e *Environment) define(name string, value any) {
	// fmt.Printf("defining %s as: %v\n", name, value)
	e.Values[name] = value
	// for k, v := range e.Values {
	// 	fmt.Printf("%s -> %v\n", k, v)
	// }
}
