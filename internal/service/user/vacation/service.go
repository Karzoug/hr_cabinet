package vacation

type service struct {
	dbRepository dbRepository
}

func NewService(dbRepository dbRepository) *service {
	return &service{
		dbRepository: dbRepository,
	}
}
