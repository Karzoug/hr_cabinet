package env

import "errors"

const (
	// Production is a constant defining the production environment.
	Production Type = iota
	// Development is a constant defining the development environment.
	Development
	// Staging is a constant defining the staging environment.
	Staging
)

const (
	productionString  = "production"
	developmentString = "development"
	stagingString     = "staging"
	unknownString     = "unknown"
)

// ErrUnknownEnv is an error returned when the environment is unknown.
var ErrUnknownEnv = errors.New("unknown environment mode")

// Type is a type of environment (production or development).
type Type int8

// String returns the string representation of the environment type variable.
func (e Type) String() string {
	switch e {
	case Development:
		return developmentString
	case Production:
		return productionString
	case Staging:
		return stagingString
	default:
		return unknownString
	}
}

func (t *Type) UnmarshalText(data []byte) error {
	switch string(data) {
	case productionString:
		*t = Production
	case developmentString:
		*t = Development
	case stagingString:
		*t = Staging
	default:
		return errors.New("unknown environment mode")
	}
	return nil
}
