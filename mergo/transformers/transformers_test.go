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
package transformers_test

import (
	"fmt"
	"math/rand"

	"github.com/imdario/mergo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/presslabs/controller-util/mergo/transformers"
)

var _ = Describe("PodSpec Transformer", func() {
	var deployment *appsv1.Deployment

	BeforeEach(func() {
		r := rand.Int31()
		name := fmt.Sprintf("depl-%d", r)
		deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
			},
			Spec: appsv1.DeploymentSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "main",
								Image: "main-image",
								Env: []corev1.EnvVar{
									{
										Name:  "TEST",
										Value: "me",
									},
								},
							},
							{
								Name:  "helper",
								Image: "helper-image",
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: 80,
										Protocol:      corev1.ProtocolTCP,
									},
									{
										Name:          "prometheus",
										ContainerPort: 9125,
										Protocol:      corev1.ProtocolTCP,
									},
								},
							},
						},
					},
				},
			},
		}
	})

	It("removes unused containers", func() {
		newSpec := corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "helper",
					Image: "helper-image",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 80,
							Protocol:      corev1.ProtocolTCP,
						},
						{
							Name:          "prometheus",
							ContainerPort: 9125,
							Protocol:      corev1.ProtocolTCP,
						},
					},
				},
			},
		}

		Expect(mergo.Merge(&deployment.Spec.Template.Spec, newSpec, mergo.WithTransformers(transformers.PodSpec))).To(Succeed())
		Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(1))
		Expect(deployment.Spec.Template.Spec.Containers[0].Name).To(Equal("helper"))
		Expect(deployment.Spec.Template.Spec.Containers[0].Ports).To(HaveLen(2))
	})
	It("allows container rename", func() {
		newSpec := corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "new-helper",
					Image: "helper-image",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 80,
							Protocol:      corev1.ProtocolTCP,
						},
						{
							Name:          "prometheus",
							ContainerPort: 9125,
							Protocol:      corev1.ProtocolTCP,
						},
					},
				},
			},
		}

		Expect(mergo.Merge(&deployment.Spec.Template.Spec, newSpec, mergo.WithTransformers(transformers.PodSpec))).To(Succeed())
		Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(1))
		Expect(deployment.Spec.Template.Spec.Containers[0].Name).To(Equal("new-helper"))
		Expect(deployment.Spec.Template.Spec.Containers[0].Ports).To(HaveLen(2))
	})
	It("merges env vars", func() {
		newSpec := corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "main",
					Image: "main-image",
					Env: []corev1.EnvVar{
						{
							Name:  "TEST-2",
							Value: "me-2",
						},
					},
				},
			},
		}

		Expect(mergo.Merge(&deployment.Spec.Template.Spec, newSpec, mergo.WithTransformers(transformers.PodSpec))).To(Succeed())
		Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(1))
		Expect(deployment.Spec.Template.Spec.Containers[0].Env).To(HaveLen(2))
		Expect(deployment.Spec.Template.Spec.Containers[0].Env[0].Name).To(Equal("TEST"))
		Expect(deployment.Spec.Template.Spec.Containers[0].Env[0].Value).To(Equal("me"))
		Expect(deployment.Spec.Template.Spec.Containers[0].Env[1].Name).To(Equal("TEST-2"))
		Expect(deployment.Spec.Template.Spec.Containers[0].Env[1].Value).To(Equal("me-2"))
	})
	It("merges container ports", func() {
		newSpec := deployment.Spec.Template.Spec.DeepCopy()
		newSpec.Containers[1].Ports = []corev1.ContainerPort{
			{
				Name:          "prometheus",
				ContainerPort: 9125,
				Protocol:      corev1.ProtocolTCP,
			},
		}
		Expect(mergo.Merge(&deployment.Spec.Template.Spec, newSpec, mergo.WithTransformers(transformers.PodSpec))).To(Succeed())
		Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(2))
		Expect(deployment.Spec.Template.Spec.Containers[1].Ports).To(HaveLen(1))
		Expect(deployment.Spec.Template.Spec.Containers[1].Ports[0].ContainerPort).To(Equal(int32(9125)))
	})
})
