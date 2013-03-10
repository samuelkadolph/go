package raw

// #include "logging.h"
import "C"

import (
	"fmt"
)

type LogLevel int

const (
	Critical = C.PHIDGET_LOG_CRITICAL
	Error    = C.PHIDGET_LOG_ERROR
	Warning  = C.PHIDGET_LOG_WARNING
	Debug    = C.PHIDGET_LOG_DEBUG
	Info     = C.PHIDGET_LOG_INFO
	Verbose  = C.PHIDGET_LOG_VERBOSE
)

func DisableLogging() error {
	return result(C.CPhidget_disableLogging())
}

func EnableLogging(level LogLevel, path string) error {
	return result(C.CPhidget_enableLogging(C.CPhidgetLog_level(level), convertString(path)))
}

func Log(level LogLevel, id, format string, args ...interface{}) error {
	return result(C._log(C.CPhidgetLog_level(level), convertString(id), convertString(fmt.Sprintf(format, args...))))
}
