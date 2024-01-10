package recovery

type service struct {
	recoveryRepository    recoveryRepository
	keyRepository         keyRepository
	notificationDeliverer notificationDeliverer
	passwordVerificator   passwordVerificator
	Domain                string
}

func NewService(rr recoveryRepository,
	kr keyRepository,
	nd notificationDeliverer,
	pv passwordVerificator,
	domain string) *service {
	return &service{
		recoveryRepository:    rr,
		keyRepository:         kr,
		notificationDeliverer: nd,
		passwordVerificator:   pv,
		Domain:                domain,
	}
}
