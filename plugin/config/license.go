package config

import "io"

//go:generate mockgen -destination=../mock/plugin/license.go -package=plugin gitlab.arboo.cn/ab-edgefusion/edgefusion-go/v2/plugin License

type License interface {
	ProtectCode() error
	CheckLicense() error
	io.Closer
}
