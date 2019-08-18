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
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
	zaplog "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	// Log is the base logger used by kubebuilder.  It delegates
	// to another logr.Logger.  You *must* call SetLogger to
	// get any actual logging.
	// Deprecated: use sigs.k8s.io/controller-runtime/pkg/log.Log instead
	Log = log.Log

	// KBLog is a base parent logger.
	// Deprecated: create your own logger. This will be removed
	KBLog logr.Logger

	// SetLogger sets a concrete logging implementation for all deferred Loggers.
	// Deprecated: use sigs.k8s.io/controller-runtime/pkg/log.SetLogger instead
	SetLogger = log.SetLogger

	// ZapLogger is a Logger implementation.
	// If development is true, a Zap development config will be used
	// (stacktraces on warnings, no sampling), otherwise a Zap production
	// config will be used (stacktraces on errors, sampling).
	// Deprecated: use sigs.k8s.io/controller-runtime/pkg/log/zap.Logger instead
	ZapLogger = zaplog.Logger

	// ZapLoggerTo returns a new Logger implementation using Zap which logs
	// to the given destination, instead of stderr.  It otherwise behaves like
	// ZapLogger.
	// Deprecated: use sigs.k8s.io/controller-runtime/pkg/log/zap.LoggerTo instead
	ZapLoggerTo = zaplog.LoggerTo

	// RawZapLoggerTo returns a new zap.Logger configured with KubeAwareEncoder
	// which logs to a given destination
	// Deprecated: use sigs.k8s.io/controller-runtime/pkg/log/zap.RawLoggerTo instead
	RawZapLoggerTo = zaplog.RawLoggerTo
)

func init() { // nolint: gochecknoinits
	KBLog = Log.WithName("kubebuilder")
}
