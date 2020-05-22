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
	"testing"

	"github.com/go-logr/zapr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/go-logr/zapr"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
)

func TestLogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Test log configuration Suite", []Reporter{printer.NewlineReporter{}})
}

var _ = Describe("Logging tests", func() {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
	}

	It("should output summary for k8s object and have stackdrive configured", func() {
		var logOut []byte
		logOutBuffer := bytes.NewBuffer(logOut)
		zapLogger := RawStackdriveZapLoggerTo(logOutBuffer, false)
		logger := zapr.NewLogger(zapLogger)
		logger = logger.WithValues("valueKey", secret)

		logger.Info("test log", "key", secret)

		Expect(zapLogger.Sync()).To(Succeed())

		// unmarshal logs and assert on them
		var data map[string]interface{}
		Expect(json.Unmarshal(logOutBuffer.Bytes(), &data)).To(Succeed())

		// check that is used the stackdriver logger
		Expect(data).To(HaveKey("severity"))

		// assert key field encoded with KubeAwareEncoder
		Expect(data).To(HaveKey("key"))
		Expect(data["key"]).To(HaveKeyWithValue("name", "test"))
		Expect(data["key"]).To(HaveKeyWithValue("namespace", "default"))

		// assert valueKey field
		Expect(data).To(HaveKey("valueKey"))
		// TODO: objects logged with logger.WithValues are not serialized using KubeAwareEncoder encoder
		//Expect(data["valueKey"]).To(HaveKeyWithValue("name", "test"))
		//Expect(data["valueKey"]).To(HaveKeyWithValue("namespace", "default"))
	})
})
