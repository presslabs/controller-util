/*
Copyright 2018 Pressinfra SRL.

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
package syncer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"golang.org/x/net/context"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"

	. "github.com/presslabs/controller-util/syncer"
)

var _ = Describe("ObjectSyncer", func() {
	var syncer *ObjectSyncer
	var deployment *appsv1.Deployment
	var recorder *record.FakeRecorder
	key := types.NamespacedName{Name: "example", Namespace: "default"}

	BeforeEach(func() {
		deployment = &appsv1.Deployment{}
		recorder = record.NewFakeRecorder(100)
	})

	AfterEach(func() {
		// nolint: errcheck
		c.Delete(context.TODO(), deployment)
	})

	When("syncing", func() {
		It("successfully creates an ownerless object when owner is nil", func() {
			syncer = NewDeploymentSyncer(nil).(*ObjectSyncer)
			Expect(Sync(context.TODO(), syncer, recorder)).To(Succeed())

			Expect(c.Get(context.TODO(), key, deployment)).To(Succeed())

			Expect(deployment.ObjectMeta.OwnerReferences).To(HaveLen(0))

			Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(deployment.Spec.Template.Spec.Containers[0].Name).To(Equal("busybox"))
			Expect(deployment.Spec.Template.Spec.Containers[0].Image).To(Equal("busybox"))

			// since this is an ownerless object, no event is emitted
			Consistently(recorder.Events).ShouldNot(Receive())
		})

		It("successfully creates an object and set owner references", func() {
			owner := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
			}
			Expect(c.Create(context.TODO(), owner)).To(Succeed())
			syncer = NewDeploymentSyncer(owner).(*ObjectSyncer)
			Expect(Sync(context.TODO(), syncer, recorder)).To(Succeed())

			Expect(c.Get(context.TODO(), key, deployment)).To(Succeed())

			Expect(deployment.ObjectMeta.OwnerReferences).To(HaveLen(1))
			Expect(deployment.ObjectMeta.OwnerReferences[0].Name).To(Equal(owner.ObjectMeta.Name))
			Expect(*deployment.ObjectMeta.OwnerReferences[0].Controller).To(BeTrue())

			var event string
			Expect(recorder.Events).To(Receive(&event))
			Expect(event).To(ContainSubstring("ExampleDeploymentSyncSuccessfull"))
			Expect(event).To(ContainSubstring("*v1.Deployment default/example created successfully"))
		})
	})
})
