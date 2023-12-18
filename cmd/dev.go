//go:build dev

package main

import "github.com/Employee-s-file-cabinet/backend/internal/config"

func init() {
	envMode = config.EnvDevelopment
}
