/*
Copyright 2019 Pressinfra SRL.

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

package syncer

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("ObjectSyncer", func() {
	Describe("stripSecrets function", func() {
		It("returns initial object when it doesn't contain secret data", func() {
			obj := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "awesome-pod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "awesome-container",
						},
					},
				},
			}

			Expect(redact(obj)).To(Equal(obj))
		})

		It("returns the object without secret data when the object is a secret", func() {
			obj := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "awesome-secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"awesome-secret-key": []byte("awesome-secret-data"),
				},
				StringData: map[string]string{
					"another-awesome-secret-key": "another-awesome-secret-data",
				},
			}

			expectedObj := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "awesome-secret",
					Namespace: "default",
				},
			}

			Expect(redact(obj)).To(Equal(expectedObj))
			Expect(obj.Data).To(HaveKey("awesome-secret-key"))
		})
	})
})
