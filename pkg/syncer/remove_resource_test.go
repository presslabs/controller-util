/*
Copyright 2023 Pressinfra SRL.

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
	"context"
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"

	"github.com/presslabs/controller-util/pkg/syncer"
)

var _ = Describe("ObjectSyncer", func() {
	var (
		removeResSyncer *syncer.RemoveResourceSyncer
		deployment      *appsv1.Deployment
		recorder        *record.FakeRecorder
		owner           *corev1.ConfigMap
		key             types.NamespacedName
	)

	BeforeEach(func() {
		r := rand.Int31() //nolint: gosec

		key = types.NamespacedName{
			Name:      fmt.Sprintf("example-%d", r),
			Namespace: fmt.Sprintf("default-%d", r),
		}
		deplLabels := map[string]string{"test": "test"}
		deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: deplLabels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: deplLabels,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "test",
								Image: "test",
							},
						},
					},
				},
			},
		}
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
		recorder = record.NewFakeRecorder(100)

		Expect(c.Create(context.TODO(), ns)).To(Succeed())
		Expect(c.Create(context.TODO(), deployment)).To(Succeed())
		Expect(c.Create(context.TODO(), owner)).To(Succeed())
	})

	When("syncing", func() {
		It("successfully deletes existing resource", func() {
			var convOk bool

			removeResSyncer, convOk = syncer.NewRemoveResourceSyncer("test-remove-resource-syncer", owner, deployment, c).(*syncer.RemoveResourceSyncer)
			Expect(convOk).To(BeTrue())
			Expect(syncer.Sync(context.TODO(), removeResSyncer, recorder)).To(Succeed())

			Expect(k8serrors.IsNotFound(c.Get(context.TODO(), key, deployment))).To(BeTrue())

			Expect(<-recorder.Events).To(Equal(
				fmt.Sprintf("Normal TestRemoveResourceSyncerSyncSuccessfull apps/v1, Kind=Deployment %s/%s successfully deleted", key.Namespace, key.Name),
			))

			// no more events
			Consistently(recorder.Events).ShouldNot(Receive())
		})

		It("skip deleting when the resource is lready deleted", func() {
			var convOk bool

			Expect(c.Delete(context.TODO(), deployment)).To(Succeed())

			removeResSyncer, convOk = syncer.NewRemoveResourceSyncer("test-remove-resource-syncer", owner, deployment, c).(*syncer.RemoveResourceSyncer)
			Expect(convOk).To(BeTrue())
			Expect(syncer.Sync(context.TODO(), removeResSyncer, recorder)).To(Succeed())

			Expect(k8serrors.IsNotFound(c.Get(context.TODO(), key, deployment))).To(BeTrue())

			// since this is an ownerless object, no event is emitted
			Consistently(recorder.Events).ShouldNot(Receive())
		})
	})
})
