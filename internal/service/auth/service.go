package auth

type service struct {
	authRepository      authRepository
	passwordVerificator passwordVerificator
	tokenManager        tokenManager
}

func NewService(ar authRepository,
	pv passwordVerificator,
	tm tokenManager) *service {
	return &service{
		authRepository:      ar,
		passwordVerificator: pv,
		tokenManager:        tm,
	}
}
