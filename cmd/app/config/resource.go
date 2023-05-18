package config

import (
	_ "embed"
)

//go:embed resource/config.yaml
var CfgBin []byte
var CfgFileName = "config.yaml"

//go:embed resource/data.db
var DataDBBin []byte
var DataDBFileName = "data.db"

//go:embed resource/nym-logo.png
var NymLogoBin []byte
