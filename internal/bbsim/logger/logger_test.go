/*
 * Copyright 2018-present Open Networking Foundation

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logger_test

import (
	bbsimLogger "github.com/opencord/bbsim/internal/bbsim/logger"
	"github.com/sirupsen/logrus"
	"gotest.tools/assert"
	"testing"
)

func Test_SetLogLevel(t *testing.T) {
	log := logrus.New()

	bbsimLogger.SetLogLevel(log, "trace", false)
	assert.Equal(t, log.Level, logrus.TraceLevel)

	bbsimLogger.SetLogLevel(log, "debug", false)
	assert.Equal(t, log.Level, logrus.DebugLevel)

	bbsimLogger.SetLogLevel(log, "info", false)
	assert.Equal(t, log.Level, logrus.InfoLevel)

	bbsimLogger.SetLogLevel(log, "warn", false)
	assert.Equal(t, log.Level, logrus.WarnLevel)

	bbsimLogger.SetLogLevel(log, "error", false)
	assert.Equal(t, log.Level, logrus.ErrorLevel)

	bbsimLogger.SetLogLevel(log, "foobar", false)
	assert.Equal(t, log.Level, logrus.DebugLevel)
}

func Test_SetLogLevelCaller(t *testing.T) {
	log := logrus.New()

	bbsimLogger.SetLogLevel(log, "debug", true)
	assert.Equal(t, log.ReportCaller, true)

	bbsimLogger.SetLogLevel(log, "debug", false)
	assert.Equal(t, log.ReportCaller, false)
}
