package server

import (
	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
)

var _ api.ServerInterface = (*server)(nil)

type server struct {
}
