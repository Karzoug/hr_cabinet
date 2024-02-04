package scan

import (
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type scan struct {
	ID          uint64    `db:"id"`
	UserID      uint64    `db:"user_id"`
	Type        scanType  `db:"type"`
	DocumentID  uint64    `db:"document_id"`
	Description string    `db:"description"`
	UploadedAt  time.Time `db:"created_at"`
}

type scanType string

const (
	scanTypePassport   scanType = "Паспорт"
	scanTypeTaxpayer   scanType = "ИНН"
	scanTypeInsurance  scanType = "СНИЛС"
	scanTypeContract   scanType = "Трудовой договор"
	scanTypePDP        scanType = "Согласие на обработку данных"
	scanTypeMilitary   scanType = "Военный билет"
	scanTypeEducation  scanType = "Документ об образовании"
	scanTypeTraining   scanType = "Сертификат"
	scanTypeBriefing   scanType = "Инструктаж"
	scanTypeWorkPermit scanType = "Разрешение на работу"
	scanTypeMarriage   scanType = "Свидетельство о браке"
	scanTypeBabyBirth  scanType = "Свидетельство о рождении"
	scanTypeOther      scanType = "Другое"
)

func convertFromDAO(s scan) model.Scan {
	var st model.ScanType
	switch s.Type {
	case scanTypePassport:
		st = model.ScanTypePassport
	case scanTypeTaxpayer:
		st = model.ScanTypeTaxpayer
	case scanTypeInsurance:
		st = model.ScanTypeInsurance
	case scanTypeContract:
		st = model.ScanTypeContract
	case scanTypePDP:
		st = model.ScanTypePDP
	case scanTypeMilitary:
		st = model.ScanTypeMilitary
	case scanTypeEducation:
		st = model.ScanTypeEducation
	case scanTypeTraining:
		st = model.ScanTypeTraining
	case scanTypeBriefing:
		st = model.ScanTypeBriefing
	case scanTypeWorkPermit:
		st = model.ScanTypeWorkPermit
	case scanTypeMarriage:
		st = model.ScanTypeMarriage
	case scanTypeBabyBirth:
		st = model.ScanTypeBabyBirth
	case scanTypeOther:
		st = model.ScanTypeOther
	}

	return model.Scan{
		ID:          s.ID,
		Type:        st,
		DocumentID:  s.DocumentID,
		Description: s.Description,
		UploadedAt:  s.UploadedAt,
	}
}

func convertToDAO(ms model.Scan) scan {
	var t scanType
	switch ms.Type {
	case model.ScanTypePassport:
		t = scanTypePassport
	case model.ScanTypeTaxpayer:
		t = scanTypeTaxpayer
	case model.ScanTypeInsurance:
		t = scanTypeInsurance
	case model.ScanTypeContract:
		t = scanTypeContract
	case model.ScanTypePDP:
		t = scanTypePDP
	case model.ScanTypeMilitary:
		t = scanTypeMilitary
	case model.ScanTypeEducation:
		t = scanTypeEducation
	case model.ScanTypeTraining:
		t = scanTypeTraining
	case model.ScanTypeBriefing:
		t = scanTypeBriefing
	case model.ScanTypeWorkPermit:
		t = scanTypeWorkPermit
	case model.ScanTypeMarriage:
		t = scanTypeMarriage
	case model.ScanTypeBabyBirth:
		t = scanTypeBabyBirth
	case model.ScanTypeOther:
		t = scanTypeOther
	}

	return scan{
		ID:          ms.ID,
		DocumentID:  ms.DocumentID,
		Type:        t,
		Description: ms.Description,
	}
}
