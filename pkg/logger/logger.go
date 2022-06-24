package logger

import "github.com/sirupsen/logrus"

func Error(args ...interface{}) error {
	logrus.Error(args...)

	// TODO: return actual errors
	return nil
}

func Errorf(format string, args ...interface{}) error {
	logrus.Errorf(format, args...)

	// TODO: return actual errors
	return nil
}
