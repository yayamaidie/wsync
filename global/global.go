package global

import (
	_ "embed"
)

//go:embed wsync.crt
var CertBytes []byte

//go:embed wsync.key
var KeyBytes []byte
