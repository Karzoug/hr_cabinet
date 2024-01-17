package recovery

type service struct {
	recoveryRepository    recoveryRepository
	keyRepository         keyRepository
	notificationDeliverer notificationDeliverer
	passwordVerificator   passwordVerificator
	Config                Config
}

func NewService(rr recoveryRepository,
	kr keyRepository,
	nd notificationDeliverer,
	pv passwordVerificator,
	cfg Config) *service {
	return &service{
		recoveryRepository:    rr,
		keyRepository:         kr,
		notificationDeliverer: nd,
		passwordVerificator:   pv,
		Config:                cfg,
	}
}
