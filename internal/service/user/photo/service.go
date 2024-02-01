package photo

type service struct {
	userRepository userRepository
	fileRepository s3FileRepository
}

func NewService(userRepository userRepository, fileRepository s3FileRepository) *service {
	return &service{
		userRepository: userRepository,
		fileRepository: fileRepository,
	}
}
