package convert

import (
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

func ToAPIScan(ms *model.Scan) api.Scan {
	var docID *uint64
	if ms.DocumentID > 0 {
		id := ms.DocumentID
		docID = &id
	}
	return api.Scan{
		ID:          ms.ID,
		Type:        api.ScanType(ms.Type),
		DocumentID:  docID,
		Description: &ms.Description,
		Url:         ms.URL,
		UploadAt:    ms.UploadedAt,
	}
}

func ToAPIScans(scans []model.Scan) []api.Scan {
	res := make([]api.Scan, len(scans))
	for i := 0; i < len(scans); i++ {
		res[i] = ToAPIScan(&scans[i])
	}
	return res
}
