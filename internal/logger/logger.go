package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type MyLogger interface {
	LogFatalLvl(err error)
	LogErrorLvl(
		msg string,
		targetPackage string,
		function string,
		err error,
		data interface{},
	)
	LogWarningLvl(
		msg string,
		targetPackage string,
		function string,
		err error,
		data interface{},
	)
	LogTraceLvl(
		msg string,
		targetPackage string,
		function string,
		err error,
		data interface{},
	)
}

type myLogger struct {
	logger *logrus.Logger
}

// Create a My Logger
func NewMyLooger(logger *logrus.Logger) MyLogger {
	return &myLogger{
		logger: logger,
	}
}

// log Fatal level error
func (l myLogger) LogFatalLvl(err error) {
	l.logger.Fatal(err)
	os.Exit(1)
}

// log Error level error
func (l myLogger) LogErrorLvl(
	msg string,
	targetPackage string,
	function string,
	err error,
	data interface{},
) {
	l.logger.WithFields(logrus.Fields{
		"package":  targetPackage,
		"function": function,
		"error":    err,
		"data":     data,
	}).Error(msg)
}

// log Warning level error
func (l myLogger) LogWarningLvl(
	msg string,
	targetPackage string,
	function string,
	err error,
	data interface{},
) {
	l.logger.WithFields(logrus.Fields{
		"package":  targetPackage,
		"function": function,
		"error":    err,
		"data":     data,
	}).Warning(msg)
}

// log Trace level error
func (l myLogger) LogTraceLvl(
	msg string,
	targetPackage string,
	function string,
	err error,
	data interface{},
) {
	l.logger.WithFields(logrus.Fields{
		"package":  targetPackage,
		"function": function,
		"error":    err,
		"data":     data,
	}).Trace(msg)
}
