package callExpr

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
)

func Got() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Error("Error hahaha")
	logrus.Info("Info")
	logrus.Info("Info1", "Info2")
}

func Want() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Error().Msg("Error hahaha")
	log.Info().Msg("Info")
	log.Info().Msgf("%v %v", "Info1", "Info2")
}
