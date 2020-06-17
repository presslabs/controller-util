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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/imdario/mergo"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/presslabs/controller-util/mergo/transformers"
)

var (
	ten32   = int32(10)
	five32  = int32(5)
	ten64   = int64(10)
	five64  = int64(5)
	trueVar = true
)

var _ = Describe("PodSpec Transformer", func() {
	var deployment *appsv1.Deployment

	BeforeEach(func() {
		r := rand.Int31()
		runtimeClass := "old-runtime-class-name"
		sharedPN := false
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
								Resources: corev1.ResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceCPU: resource.MustParse("100m"),
									},
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "code",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							},
							{
								Name: "media",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							},
						},
						Affinity: &corev1.Affinity{
							NodeAffinity: &corev1.NodeAffinity{
								RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
									NodeSelectorTerms: []corev1.NodeSelectorTerm{
										{
											MatchExpressions: []corev1.NodeSelectorRequirement{
												{
													Key:      "some-label-key",
													Operator: corev1.NodeSelectorOpIn,
													Values:   []string{"some-label-value"},
												},
											},
										},
									},
								},
							},
						},
						PriorityClassName:             "old-priority-class",
						TerminationGracePeriodSeconds: &ten64,
						Priority:                      &ten32,
						RuntimeClassName:              &runtimeClass,
						HostIPC:                       false,
						ShareProcessNamespace:         &sharedPN,
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

	It("allows container image update", func() {
		newSpec := deployment.Spec.Template.Spec.DeepCopy()
		newSpec.Containers[0].Image = "main-image-v2"
		Expect(mergo.Merge(&deployment.Spec.Template.Spec, newSpec, mergo.WithTransformers(transformers.PodSpec))).To(Succeed())
		Expect(deployment.Spec.Template.Spec.Containers[0].Name).To(Equal("main"))
		Expect(deployment.Spec.Template.Spec.Containers[0].Image).To(Equal("main-image-v2"))
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
		Expect(deployment.Spec.Template.Spec.Containers[0].Env).To(HaveLen(1))
		Expect(deployment.Spec.Template.Spec.Containers[0].Env[0].Name).To(Equal("TEST-2"))
		Expect(deployment.Spec.Template.Spec.Containers[0].Env[0].Value).To(Equal("me-2"))
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

	It("allows prepending volume", func() {
		newSpec := deployment.Spec.Template.Spec.DeepCopy()
		newSpec.Volumes = []corev1.Volume{
			{
				Name: "config",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
			{
				Name: "code",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{},
				},
			},
			{
				Name: "media",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		}
		Expect(mergo.Merge(&deployment.Spec.Template.Spec, newSpec, mergo.WithTransformers(transformers.PodSpec))).To(Succeed())
		Expect(deployment.Spec.Template.Spec.Volumes).To(HaveLen(3))
		Expect(deployment.Spec.Template.Spec.Volumes[0].Name).To(Equal(newSpec.Volumes[0].Name))
		Expect(deployment.Spec.Template.Spec.Volumes[1].Name).To(Equal(newSpec.Volumes[1].Name))
		Expect(deployment.Spec.Template.Spec.Volumes[2].Name).To(Equal(newSpec.Volumes[2].Name))

		Expect(deployment.Spec.Template.Spec.Volumes[1].EmptyDir).To(BeNil())
		Expect(deployment.Spec.Template.Spec.Volumes[1].HostPath).ToNot(BeNil())
	})

	It("allows replacing volume list", func() {
		newSpec := deployment.Spec.Template.Spec.DeepCopy()
		newSpec.Volumes = []corev1.Volume{
			{
				Name: "config",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		}
		Expect(mergo.Merge(&deployment.Spec.Template.Spec, newSpec, mergo.WithTransformers(transformers.PodSpec))).To(Succeed())
		Expect(deployment.Spec.Template.Spec.Volumes).To(HaveLen(1))
		Expect(deployment.Spec.Template.Spec.Volumes[0].Name).To(Equal(newSpec.Volumes[0].Name))
	})

	It("override existing affinity with new one, instead of merging them", func() {
		newSpec := deployment.Spec.Template.Spec.DeepCopy()
		newAffinity := &corev1.Affinity{
			NodeAffinity: &corev1.NodeAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
					{
						Weight: 42,
						Preference: corev1.NodeSelectorTerm{
							MatchExpressions: []corev1.NodeSelectorRequirement{
								{
									Key:      "some-label-key",
									Operator: corev1.NodeSelectorOpExists,
								},
							},
						},
					},
				},
			},
		}

		newSpec.Affinity = newAffinity
		Expect(mergo.Merge(&deployment.Spec.Template.Spec, newSpec, mergo.WithTransformers(transformers.PodSpec))).To(Succeed())
		Expect(deployment.Spec.Template.Spec.Affinity).To(Equal(newAffinity))
	})

	It("should update unknown transformer type like Quantity", func() {
		oldSpec := deployment.Spec.Template.Spec
		newSpec := deployment.Spec.Template.Spec.DeepCopy()

		newCPU := resource.MustParse("3")
		newSpec.Containers[1].Resources.Requests[corev1.ResourceCPU] = newCPU

		Expect(mergo.Merge(&oldSpec, newSpec, mergo.WithTransformers(transformers.PodSpec))).To(Succeed())
		Expect(oldSpec.Containers[1].Resources.Requests[corev1.ResourceCPU]).To(Equal(newCPU))

		// don't update with empty value
		newSpec.Containers[1].Resources.Requests = corev1.ResourceList{}
		Expect(mergo.Merge(&oldSpec, newSpec, mergo.WithTransformers(transformers.PodSpec))).To(Succeed())
		Expect(oldSpec.Containers[1].Resources.Requests[corev1.ResourceCPU]).To(Equal(newCPU))
	})

	It("updates the filds for string, *string, *int32, *int64, bool, *bool", func() {
		oldSpec := deployment.Spec.Template.Spec
		newSpec := deployment.Spec.Template.Spec.DeepCopy()

		// type string
		newSpec.PriorityClassName = "new-priority-class"
		// type *int64
		newSpec.TerminationGracePeriodSeconds = &five64
		// type *int32
		newSpec.Priority = &five32
		// type *string
		rcn := "new-runtime-class"
		newSpec.RuntimeClassName = &rcn
		// type bool
		newSpec.HostIPC = true
		// type *bool
		newSpec.ShareProcessNamespace = &trueVar

		Expect(mergo.Merge(&oldSpec, newSpec, mergo.WithTransformers(transformers.PodSpec))).To(Succeed())

		Expect(oldSpec.PriorityClassName).To(Equal(newSpec.PriorityClassName))
		Expect(oldSpec.TerminationGracePeriodSeconds).To(Equal(&five64))
		Expect(oldSpec.Priority).To(Equal(&five32))
		Expect(oldSpec.RuntimeClassName).To(Equal(&rcn))
		Expect(oldSpec.HostIPC).To(Equal(newSpec.HostIPC))
		Expect(oldSpec.ShareProcessNamespace).To(Equal(newSpec.ShareProcessNamespace))
	})

	It("should not update string with empty value", func() {
		oldSpec := deployment.Spec.Template.Spec
		newSpec := deployment.Spec.Template.Spec.DeepCopy()
		newSpec.PriorityClassName = ""

		Expect(mergo.Merge(&oldSpec, newSpec, mergo.WithTransformers(transformers.PodSpec))).To(Succeed())
		Expect(oldSpec.PriorityClassName).To(Equal("old-priority-class"))
	})
})
