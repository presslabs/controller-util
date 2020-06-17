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
	"k8s.io/apimachinery/pkg/runtime"
)

// GetObject returns the ObjectSyncer subject
// Deprecated: use github.com/presslabs/controller-util/syncer.Object() instead.
func (s *ObjectSyncer) GetObject() interface{} {
	return s.Object()
}

// GetOwner returns the ObjectSyncer owner
// Deprecated: use github.com/presslabs/controller-util/syncer.ObjectOwner() instead.
func (s *ObjectSyncer) GetOwner() runtime.Object {
	return s.ObjectOwner()
}
