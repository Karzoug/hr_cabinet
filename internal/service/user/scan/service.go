package scan

type service struct {
	userRepository userRepository
	scanRepository scanRepository
	fileRepository s3FileRepository
}

func NewService(userRepository userRepository,
	scanRepository scanRepository,
	fileRepository s3FileRepository) *service {
	return &service{
		userRepository: userRepository,
		scanRepository: scanRepository,
		fileRepository: fileRepository,
	}
}
