package visa

import (
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type visa struct {
	ID          uint64    `db:"id"`
	Number      string    `db:"number"`
	Type        string    `db:"type"`
	IssuedState *string   `db:"issued_state"`
	ValidTo     time.Time `db:"valid_to"`
	ValidFrom   time.Time `db:"valid_from"`
}

func convertFromDAO(v visa) model.Visa {
	return model.Visa{
		ID:          v.ID,
		Number:      v.Number,
		Type:        model.VisaType(v.Type),
		IssuedState: v.IssuedState,
		ValidTo:     v.ValidTo,
		ValidFrom:   v.ValidFrom,
	}
}

func convertToDAO(mv model.Visa) visa {
	return visa{
		ID:          mv.ID,
		Number:      mv.Number,
		Type:        string(mv.Type),
		IssuedState: mv.IssuedState,
		ValidTo:     mv.ValidTo,
		ValidFrom:   mv.ValidFrom,
	}
}
