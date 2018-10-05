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

// Package transformers provide mergo transformers for Kubernetes objects
package transformers

import (
	"fmt"
	"reflect"

	"github.com/imdario/mergo"

	corev1 "k8s.io/api/core/v1"
)

type transformerMap map[reflect.Type]func(dst, src reflect.Value) error

// PodSpec mergo transformers for corev1.PodSpec
var PodSpec transformerMap

func init() {
	PodSpec = transformerMap{
		reflect.TypeOf([]corev1.Container{}):            PodSpec.mergeListByKey("Name"),
		reflect.TypeOf([]corev1.ContainerPort{}):        PodSpec.mergeListByKey("ContainerPort"),
		reflect.TypeOf([]corev1.EnvVar{}):               PodSpec.mergeListByKey("Name", mergo.WithOverride),
		reflect.TypeOf(corev1.EnvVar{}):                 PodSpec.overrideFields("Value", "ValueFrom"),
		reflect.TypeOf([]corev1.Toleration{}):           PodSpec.mergeListByKey("Key"),
		reflect.TypeOf([]corev1.Volume{}):               PodSpec.mergeListByKey("Name"),
		reflect.TypeOf([]corev1.LocalObjectReference{}): PodSpec.mergeListByKey("Name"),
		reflect.TypeOf([]corev1.HostAlias{}):            PodSpec.mergeListByKey("IP"),
		reflect.TypeOf([]corev1.VolumeMount{}):          PodSpec.mergeListByKey("MountPath"),
	}
}

func (s transformerMap) Transformer(t reflect.Type) func(dst, src reflect.Value) error {
	if fn, ok := s[t]; ok {
		return fn
	}
	return nil
}

func (s *transformerMap) mergeByKey(key string, dst, elem reflect.Value, opts ...func(*mergo.Config)) error {
	elemKey := elem.FieldByName(key)
	for i := 0; i < dst.Len(); i++ {
		dstKey := dst.Index(i).FieldByName(key)
		if elemKey.Kind() != dstKey.Kind() {
			return fmt.Errorf("cannot merge when key type differs")
		}
		eq := false
		switch elemKey.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			eq = elemKey.Int() == dstKey.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			eq = elemKey.Uint() == dstKey.Uint()
		case reflect.String:
			eq = elemKey.String() == dstKey.String()
		case reflect.Float32, reflect.Float64:
			eq = elemKey.Float() == dstKey.Float()
		}
		if eq {
			opts = append(opts, mergo.WithTransformers(s))
			return mergo.Merge(dst.Index(i).Addr().Interface(), elem.Interface(), opts...)
		}
	}
	dst.Set(reflect.Append(dst, elem))
	return nil
}

func (s *transformerMap) mergeListByKey(key string, opts ...func(*mergo.Config)) func(_, _ reflect.Value) error {
	return func(dst, src reflect.Value) error {
		for i := 0; i < src.Len(); i++ {
			elem := src.Index(i)
			err := s.mergeByKey(key, dst, elem, opts...)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (s *transformerMap) overrideFields(fields ...string) func(_, _ reflect.Value) error {
	return func(dst, src reflect.Value) error {
		for _, field := range fields {
			srcValue := src.FieldByName(field)
			dst.FieldByName(field).Set(srcValue)
		}
		return nil
	}
}
