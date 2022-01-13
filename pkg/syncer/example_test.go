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
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/presslabs/controller-util/pkg/syncer"
)

var (
	recorder record.EventRecorder
	owner    client.Object
	log      = logf.Log.WithName("controllerutil-examples")
)

func NewDeploymentSyncer(owner client.Object) syncer.Interface {
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example",
			Namespace: "default",
		},
	}

	// c is client.Client
	return syncer.NewObjectSyncer("ExampleDeployment", owner, deploy, c, func() error {
		// Deployment selector is immutable so we set this value only if
		// a new object is going to be created
		if deploy.ObjectMeta.CreationTimestamp.IsZero() {
			deploy.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: map[string]string{"foo": "bar"},
			}
		}

		// update the Deployment pod template
		deploy.Spec.Template = corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"foo": "bar",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "busybox",
						Image: "busybox",
					},
				},
			},
		}

		return nil
	})
}

func ExampleNewObjectSyncer() {
	// recorder is record.EventRecorder
	// owner is the owner for the syncer subject

	deploymentSyncer := NewDeploymentSyncer(owner)
	err := syncer.Sync(context.TODO(), deploymentSyncer, recorder)
	if err != nil {
		log.Error(err, "unable to sync")
	}
}
