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

// nolint: errcheck
package syncer_test

import (
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"

	. "github.com/presslabs/controller-util/pkg/syncer"
)

var _ = Describe("ObjectSyncer", func() {
	var (
		syncer     *ObjectSyncer
		deployment *appsv1.Deployment
		recorder   *record.FakeRecorder
		owner      *corev1.ConfigMap
		key        types.NamespacedName
	)

	BeforeEach(func() {
		r := rand.Int31()

		key = types.NamespacedName{
			Name:      fmt.Sprintf("example-%d", r),
			Namespace: fmt.Sprintf("default-%d", r),
		}

		deployment = &appsv1.Deployment{}
		recorder = record.NewFakeRecorder(100)
		owner = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
		}
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: key.Namespace,
			},
		}
		Expect(c.Create(context.TODO(), ns)).To(Succeed())
		Expect(c.Create(context.TODO(), owner)).To(Succeed())
	})

	AfterEach(func() {
		c.Delete(context.TODO(), deployment)
		c.Delete(context.TODO(), owner)
	})

	When("syncing", func() {
		It("successfully creates an ownerless object when owner is nil", func() {
			syncer = NewDeploymentSyncer(nil, key).(*ObjectSyncer)
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
			syncer = NewDeploymentSyncer(owner, key).(*ObjectSyncer)
			Expect(Sync(context.TODO(), syncer, recorder)).To(Succeed())

			Expect(c.Get(context.TODO(), key, deployment)).To(Succeed())

			Expect(deployment.ObjectMeta.OwnerReferences).To(HaveLen(1))
			Expect(deployment.ObjectMeta.OwnerReferences[0].Name).To(Equal(owner.ObjectMeta.Name))
			Expect(*deployment.ObjectMeta.OwnerReferences[0].Controller).To(BeTrue())

			var event string
			Expect(recorder.Events).To(Receive(&event))
			Expect(event).To(ContainSubstring("ExampleDeploymentSyncSuccessfull"))
			Expect(event).To(ContainSubstring(
				fmt.Sprintf("apps/v1, Kind=Deployment %s/%s created successfully", key.Namespace, key.Name),
			))
		})

		It("should ignore ErrIgnore", func() {
			obj := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example",
					Namespace: "default",
				},
			}

			syn := NewObjectSyncer("xxx", nil, obj, c, func() error {
				return ErrIgnore
			})

			Expect(Sync(context.TODO(), syn, recorder)).To(Succeed())
		})

		When("owner is deleted", func() {
			BeforeEach(func() {
				// set deletion timestamp on owner resource
				now := metav1.Now()
				owner.ObjectMeta.DeletionTimestamp = &now
			})
			It("should not create the resource if not exists", func() {
				syncer = NewDeploymentSyncer(owner, key).(*ObjectSyncer)
				Expect(Sync(context.TODO(), syncer, recorder)).To(Succeed())

				// check deployment is not created
				Expect(c.Get(context.TODO(), key, deployment)).ToNot(Succeed())
			})

			It("should not set owner reference", func() {
				// create the deployment
				syncer = NewDeploymentSyncer(nil, key).(*ObjectSyncer)
				Expect(Sync(context.TODO(), syncer, recorder)).To(Succeed())

				// try to set owner reference
				syncer = NewDeploymentSyncer(owner, key).(*ObjectSyncer)
				Expect(Sync(context.TODO(), syncer, recorder)).To(Succeed())

				// check deployment does not have owner reference set
				Expect(c.Get(context.TODO(), key, deployment)).To(Succeed())
				Expect(deployment.ObjectMeta.OwnerReferences).To(HaveLen(0))
			})
		})

	})
})
