package config

import "fmt"

type Environment int

const (
	Local Environment = iota
	Development
	QA
	Staging
	Production
)

func (e Environment) String() string {
	switch e {
	case Local:
		return "local"
	case Development:
		return "development"
	case Staging:
		return "staging"
	case QA:
		return "qa"
	case Production:
		return "production"
	default:
		panic("Unhandled config environment")
	}
}

func EnvironmentFromString(e string) Environment {
	for _, env := range []Environment{Local, Development, QA, Staging, Production} {
		if e == env.String() {
			return env
		}
	}

	panic(fmt.Sprintf("Unhandled config environment: %s", e))
}
