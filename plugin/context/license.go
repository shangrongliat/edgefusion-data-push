package context

import (
	"edgefusion-data-push/plugin/config"
)

type license struct {
}

func init() {
	RegisterFactory("defaultlicense", New)
}

func New() (Plugin, error) {
	return &license{}, nil
}

var _ config.License = &license{}

func (l *license) ProtectCode() error {
	return nil
}

func (l *license) CheckLicense() error {
	return nil
}

func (l *license) Close() error {
	return nil
}
