/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package log

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	once   sync.Once
	logger *Logger
)

type StdLogger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

type Logger struct {
	*logrus.Logger
}

func Initialize(level string) {
	once.Do(func() {
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			lvl = logrus.DebugLevel
		}

		if logger == nil {
			l := logrus.New()
			logger = &Logger{
				Logger: l,
			}
			l.SetLevel(lvl)
		}
	})
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Debug(args ...interface{}) {
	logger.Debugln(args...)
}

func Info(args ...interface{}) {
	logger.Infoln(args...)
}

func Warn(args ...interface{}) {
	logger.Warnln(args...)
}

func Error(args ...interface{}) {
	logger.Errorln(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatalln(args...)
}

func Panic(args ...interface{}) {
	logger.Panicln(args...)
}
