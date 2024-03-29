/*
Copyright 2020 Pressinfra SRL

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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test log configuration Suite")
}

var _ = Describe("Logging tests", func() {
	Context("production stackdrive logger", func() {
		var (
			name, ns     string
			logOutBuffer *bytes.Buffer
			zapLogger    *zap.Logger
			logger       logr.Logger
			secret       *corev1.Secret
		)

		BeforeEach(func() {
			r := rand.Int31() //nolint: gosec
			name = fmt.Sprintf("test-%d", r)
			ns = fmt.Sprintf("default-%d", r)

			secret = &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: ns,
				},
			}

			var logOut []byte
			logOutBuffer = bytes.NewBuffer(logOut)
			zapLogger = RawStackdriverZapLoggerTo(logOutBuffer, false)
			logger = zapr.NewLogger(zapLogger)
		})

		It("should output a summary for k8s object", func() {
			// log new entry and flush it
			logger.Info("test log", "key", secret)
			Expect(zapLogger.Sync()).To(Succeed())

			// unmarshal logs and assert on them
			var data map[string]interface{}
			Expect(json.Unmarshal(logOutBuffer.Bytes(), &data)).To(Succeed())

			// check that is used the stackdriver logger
			Expect(data).To(HaveKey("severity"))

			// assert key field encoded with KubeAwareEncoder
			Expect(data).To(HaveKey("key"))
			Expect(data["key"]).To(HaveKeyWithValue("name", name))
			Expect(data["key"]).To(HaveKeyWithValue("namespace", ns))
		})

		It("should output summary even if uses log.WithValues", func() {
			// NOTE: objects logged with logger.WithValues are not serialized using KubeAwareEncoder encoder
			Skip("bug not fixed")

			// set WithValues a key
			logger = logger.WithValues("withValuesKey", secret)

			// log new entry and flush it
			logger.Info("test log", "key", secret)
			Expect(zapLogger.Sync()).To(Succeed())

			// unmarshal logs and assert on them
			var data map[string]interface{}
			Expect(json.Unmarshal(logOutBuffer.Bytes(), &data)).To(Succeed())

			// assert withValuesKey field
			Expect(data).To(HaveKey("withValuesKey"))
			Expect(data["withValuesKey"]).To(HaveKeyWithValue("name", name))
			Expect(data["withValuesKey"]).To(HaveKeyWithValue("namespace", ns))
		})
	})

	Context("development stackdrive logger", func() {
		var (
			logOutBuffer *bytes.Buffer
			zapLogger    *zap.Logger
			logger       logr.Logger
		)

		BeforeEach(func() {
			var logOut []byte
			logOutBuffer = bytes.NewBuffer(logOut)
			zapLogger = RawStackdriverZapLoggerTo(logOutBuffer, true)
			logger = zapr.NewLogger(zapLogger)
		})

		It("should print stacktrace in development mode", func() {
			logger.Error(errors.New("test error message"), "logging a stacktrace") //nolint: goerr113

			// assert a piece of stacktrace
			Expect(logOutBuffer.String()).To(ContainSubstring("github.com/onsi/ginkgo/v2"))
		})
	})
})
