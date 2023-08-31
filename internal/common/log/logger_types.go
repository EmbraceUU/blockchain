package log

import (
	"github.com/sirupsen/logrus"
)

// Logging struct that holds all user configurable options for the logger
type Logging struct {
	Enabled              bool   `json:"enabled,omitempty"`
	File                 string `json:"file"`
	ColourOutput         bool   `json:"colour"`
	ColourOutputOverride bool   `json:"colourOverride,omitempty"`
	Level                string `json:"level"`
	Rotate               bool   `json:"rotate"`
	LogPath              string `json:"logpath"`
	MaxAge               int    `json:"maxage"`
	Rotationtime         int    `json:"rotationtime"`
}

var (
	defaultLogger *logrus.Logger
)
