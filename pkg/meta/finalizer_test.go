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

package meta

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Meta Package Finalizer", func() {
	const fin = "myF"

	DescribeTable("at AddFinalizer function call", func(existing, expected []string) {
		meta := &metav1.ObjectMeta{
			Finalizers: existing,
		}
		AddFinalizer(meta, fin)
		Expect(meta.Finalizers).To(Equal(expected))
	},
		Entry("add if not present", []string{"f1", "f2"}, []string{"f1", "f2", fin}),
		Entry("no add if present", []string{"f1", fin, "f2"}, []string{"f1", fin, "f2"}),
	)

	DescribeTable("at HasFinalizer function call", func(existing []string, expected bool) {
		meta := &metav1.ObjectMeta{
			Finalizers: existing,
		}
		Expect(HasFinalizer(meta, fin)).To(Equal(expected))
	},
		Entry("returns false if not present", []string{"f1", "f2"}, false),
		Entry("returns true if present", []string{"f1", fin, "f2"}, true),
	)

	DescribeTable("at RemoveFinalizer function call", func(existing, expected []string) {
		meta := &metav1.ObjectMeta{
			Finalizers: existing,
		}
		RemoveFinalizer(meta, fin)
		Expect(meta.Finalizers).To(Equal(expected))
	},
		Entry("no remove if not present", []string{"f1", "f2"}, []string{"f1", "f2"}),
		Entry("remove from middle", []string{"f1", fin, "f2"}, []string{"f1", "f2"}),
		Entry("remove from begin", []string{fin, "f1", "f2"}, []string{"f1", "f2"}),
		Entry("remove from end", []string{"f1", "f2", fin}, []string{"f1", "f2"}),
	)
})
