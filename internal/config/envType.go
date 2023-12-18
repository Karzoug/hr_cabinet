package config

import "errors"

const (
	// EnvProduction is a constant defining the production environment.
	EnvProduction EnvType = iota
	// EnvDevelopment is a constant defining the development environment.
	EnvDevelopment
)

const (
	envProductionString  = "production"
	envDevelopmentString = "development"
	envUnknownString     = "unknown"
)

// ErrUnknownEnv is an error returned when the environment is unknown.
var ErrUnknownEnv = errors.New("unknown environment mode")

// EnvType is a type of environment (production or development).
type EnvType int8

// String returns the string representation of the environment type variable.
func (e EnvType) String() string {
	switch e {
	case EnvDevelopment:
		return envDevelopmentString
	case EnvProduction:
		return envProductionString
	default:
		return envUnknownString
	}
}
