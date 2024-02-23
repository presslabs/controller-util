/*
Copyright 2024 Pressinfra SRL.

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

package predicate

import (
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Class Predicate", func() {
	var (
		p      *FilterByClassPredicate
		r      int32
		class  string
		annKey string
	)

	BeforeEach(func() {
		r = rand.Int31() //nolint: gosec

		class = fmt.Sprintf("class-%d", r)
		annKey = fmt.Sprintf("ann.key/%d", r)

		p = NewFilterByClassPredicate(class, annKey)
	})

	It("doesn't match class if annotations are empty", func() {
		Expect(p.matchesClass(&metav1.ObjectMeta{
			Annotations: map[string]string{},
		})).To(BeFalse())
	})

	It("doesn't match class if annotations doenst have the class key", func() {
		Expect(p.matchesClass(&metav1.ObjectMeta{
			Annotations: map[string]string{
				"not.class.key": "not.class.value",
			},
		})).To(BeFalse())
	})

	It("matches class if class annotation is empty, but the default class and predicate class are the same", func() {
		p.WithDefaultClass(p.class)
		Expect(p.matchesClass(&metav1.ObjectMeta{
			Annotations: map[string]string{
				p.annKey: "",
			},
		})).To(BeTrue())
	})

	It("doenst match class if class annotation is empty, but the default class and predicate class are different", func() {
		p.WithDefaultClass("another-class")
		Expect(p.matchesClass(&metav1.ObjectMeta{
			Annotations: map[string]string{
				p.annKey: "",
			},
		})).To(BeFalse())
	})

	It("matches class if class annotation and predicate class are same", func() {
		Expect(p.matchesClass(&metav1.ObjectMeta{
			Annotations: map[string]string{
				p.annKey: p.class,
			},
		})).To(BeTrue())
	})
})
