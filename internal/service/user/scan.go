package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Employee-s-file-cabinet/backend/internal/repo/s3"
	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

const (
	MaxScanSize = 20 << 20 // bytes
)

type scanUserRepo interface {
	Exist(ctx context.Context, userID uint64) (bool, error)
}

type scanRepo interface {
	Get(ctx context.Context, userID, scanID uint64) (*model.Scan, error)
	List(ctx context.Context, userID uint64) ([]model.Scan, error)
	Add(ctx context.Context, userID uint64, ms model.Scan) (uint64, error)
}

type scanFileRepo interface {
	Upload(context.Context, s3.File) error
	PresignedURL(ctx context.Context, prefix, name string) (string, error)
}

type ScanUseCase struct {
	userRepo scanUserRepo
	scanRepo scanRepo
	fileRepo scanFileRepo
}

func NewScanUseCase(userRepository scanUserRepo,
	scanRepository scanRepo,
	fileRepository scanFileRepo) ScanUseCase {
	return ScanUseCase{
		userRepo: userRepository,
		scanRepo: scanRepository,
		fileRepo: fileRepository,
	}
}

func (s ScanUseCase) Get(ctx context.Context, userID, scanID uint64) (*model.Scan, error) {
	const op = "user service: get scan"

	sc, err := s.scanRepo.Get(ctx, userID, scanID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "user or scan file not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// objectName in s3: {user_id}_{scan_type}-{document_id}
	// Ex.: 3_passport-2, 7_pdp-1, 21_contract-4
	if sc.URL, err = s.fileRepo.PresignedURL(ctx,
		strconv.Itoa(int(userID)),
		fmt.Sprintf("%s-%d", sc.Type, sc.DocumentID)); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sc, nil
}

func (s ScanUseCase) List(ctx context.Context, userID uint64) ([]model.Scan, error) {
	const op = "user service: list scans"

	scans, err := s.scanRepo.List(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return nil, serr.NewError(serr.NotFound, "user not found")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return scans, nil
}

func (s ScanUseCase) Upload(ctx context.Context, userID uint64, ms model.Scan, f model.File) (uint64, error) {
	const op = "user service: upload scan"

	if f.Size > MaxScanSize {
		return 0, serr.NewError(serr.ContentTooLarge, "scan file size too large")
	}

	if exist, err := s.userRepo.Exist(ctx, userID); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	} else if !exist {
		return 0, serr.NewError(serr.NotFound, "user not found")
	}

	// Prefix={user_id}; Name={scan_type}-{document_id}
	// objectName in s3: Prefix_Name - {user_id}_{scan_type}-{document_id}
	// Ex.: 3_passport-2, 7_pdp-1, 21_contract-4
	if err := s.fileRepo.Upload(ctx, s3.File{
		Prefix:      strconv.FormatUint(userID, 10),
		Name:        fmt.Sprintf("%s-%d", ms.Type, ms.DocumentID),
		Reader:      f.Reader,
		Size:        f.Size,
		ContentType: f.ContentType,
	}); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.scanRepo.Add(ctx, userID, ms)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
