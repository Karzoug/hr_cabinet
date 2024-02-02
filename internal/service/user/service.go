package user

import (
	"context"

	"github.com/Employee-s-file-cabinet/backend/internal/repo/s3"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/contract"
	contrRepo "github.com/Employee-s-file-cabinet/backend/internal/service/user/contract/repo/postgresql"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/education"
	edRepo "github.com/Employee-s-file-cabinet/backend/internal/service/user/education/repo/postgresql"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/passport"
	pspRepo "github.com/Employee-s-file-cabinet/backend/internal/service/user/passport/repo/postgresql"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/photo"
	fileRepo "github.com/Employee-s-file-cabinet/backend/internal/service/user/repo/s3"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/scan"
	scanRepo "github.com/Employee-s-file-cabinet/backend/internal/service/user/scan/repo/postgresql"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/training"
	trRepo "github.com/Employee-s-file-cabinet/backend/internal/service/user/training/repo/postgresql"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/user"
	userRepo "github.com/Employee-s-file-cabinet/backend/internal/service/user/user/repo/postgresql"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/vacation"
	vacRepo "github.com/Employee-s-file-cabinet/backend/internal/service/user/vacation/repo/postgresql"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/visa"
	visaRepo "github.com/Employee-s-file-cabinet/backend/internal/service/user/visa/repo/postgresql"
	pq "github.com/Employee-s-file-cabinet/backend/pkg/postgresql"
)

type Service struct {
	Contract  contract.Subservice
	Education education.Subservice
	Passport  passport.Subservice
	Photo     photo.Subservice
	Scan      scan.Subservice
	Training  training.Subservice
	User      user.Subservice
	Vacation  vacation.Subservice
	Visa      visa.Subservice
}

func New(ctx context.Context,
	db pq.DB,
	client s3.Client) (Service, error) {
	fileRepo, err := fileRepo.New(ctx, client)
	if err != nil {
		return Service{}, err
	}
	return Service{
		Contract:  contract.New(contrRepo.New(db)),
		Education: education.New(edRepo.New(db)),
		Passport:  passport.New(pspRepo.New(db)),
		Photo:     photo.New(userRepo.New(db), fileRepo),
		Scan:      scan.New(userRepo.New(db), scanRepo.New(db), fileRepo),
		Training:  training.New(trRepo.New(db)),
		User:      user.New(userRepo.New(db)),
		Vacation:  vacation.New(vacRepo.New(db)),
		Visa:      visa.New(visaRepo.New(db)),
	}, nil
}
