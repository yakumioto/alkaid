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

func Initialize(level string) {
	once.Do(func() {
		l, err := logrus.ParseLevel(level)
		if err != nil {
			l = logrus.DebugLevel
		}
		if logger == nil {
			log := logrus.New()
			log.SetLevel(l)
			logger = &Logger{StdLogger: log}
		}
	})
}

func NewLogger() *Logger {
	return &Logger{
		StdLogger: logrus.New(),
	}
}

type Logger struct {
	StdLogger
}

func (l *Logger) WithField(key string, value interface{}) StdLogger {
	return &Logger{
		StdLogger: l.WithField(key, value),
	}
}

func (l *Logger) WithFields(fields Fields) StdLogger {
	return &Logger{
		StdLogger: l.WithFields(fields),
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Warningf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Errorf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatalf(format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Panicf(format, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Warn(args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.Warning(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Error(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Fatal(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.Panic(args...)
}
