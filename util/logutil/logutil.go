package logutil

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"strings"
)

func LoglvlFromString(lv string) (log.Lvl, error) {
	if lv == "" {
		// default
		return log.INFO, nil
	}

	switch strings.ToLower(lv) {
	case "debug":
		return log.DEBUG, nil
	case "DEBUG":
		return log.DEBUG, nil
	case "info":
		return log.INFO, nil
	case "INFO":
		return log.INFO, nil
	case "warn":
		return log.WARN, nil
	case "WARN":
		return log.WARN, nil
	case "error":
		return log.ERROR, nil
	case "ERROR":
		return log.ERROR, nil
	}

	return log.INFO, fmt.Errorf("invalid log level %s", lv)
}
