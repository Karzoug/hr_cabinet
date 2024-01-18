package model

import "time"

type Scan struct {
	ID          uint64
	DocumentID  uint64
	Type        ScanType
	Description string
	URL         string
	UploadedAt  time.Time
}

type ScanType string

const (
	ScanTypePassport   ScanType = "passport"
	ScanTypeTaxpayer   ScanType = "taxpayer"
	ScanTypeInsurance  ScanType = "insurance"
	ScanTypeContract   ScanType = "contract"
	ScanTypePDP        ScanType = "personal_data_processing"
	ScanTypeMilitary   ScanType = "military"
	ScanTypeEducation  ScanType = "education"
	ScanTypeTraining   ScanType = "training"
	ScanTypeBriefing   ScanType = "briefing"
	ScanTypeWorkPermit ScanType = "work_permit"
	ScanTypeMarriage   ScanType = "marriage"
	ScanTypeBabyBirth  ScanType = "baby_birth"
	ScanTypeOther      ScanType = "other"
)
