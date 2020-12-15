/*
Copyright 2019 Pressinfra SRL

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

import (
	"io"

	"github.com/blendle/zapdriver"
	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log"
	zaplog "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	// Log is the base logger used by kubebuilder.  It delegates
	// to another logr.Logger.  You *must* call SetLogger to
	// get any actual logging.
	Log = log.Log

	// KBLog is a base parent logger.
	KBLog logr.Logger

	// SetLogger sets a concrete logging implementation for all deferred Loggers.
	SetLogger = log.SetLogger
)

func init() { // nolint: gochecknoinits
	KBLog = Log.WithName("kubebuilder")
}

// RawStackdriverZapLoggerTo returns a new zap.Logger configured with KubeAwareEncoder and StackDriverEncoder.
func RawStackdriverZapLoggerTo(destWriter io.Writer, development bool, opts ...zap.Option) *zap.Logger {
	return zaplog.NewRaw(zaplog.UseDevMode(development), zaplog.WriteTo(destWriter), withStackDriverEncoder(), zaplog.RawZapOpts(opts...))
}

func withStackDriverEncoder() zaplog.Opts {
	return func(o *zaplog.Options) {
		var enc zapcore.Encoder

		if o.Development {
			encCfg := zapdriver.NewDevelopmentEncoderConfig()
			enc = zapcore.NewConsoleEncoder(encCfg)
		} else {
			encCfg := zapdriver.NewProductionEncoderConfig()
			enc = zapcore.NewJSONEncoder(encCfg)
		}

		o.Encoder = enc
	}
}
