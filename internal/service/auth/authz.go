package auth

import (
	"github.com/casbin/casbin/v2"
)

func (s *service) PolicyEnforcer() (*casbin.Enforcer, error) {
	return casbin.NewEnforcer("policy_models/rest.conf", s.authRepository.PolicyAdapter())
}
