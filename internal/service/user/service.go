package user

type service struct {
	userRepository userRepository
}

func NewService(userRepository userRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}
