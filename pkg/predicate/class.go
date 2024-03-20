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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// FilterByClassPredicate allows filtering by the class annotation.
type FilterByClassPredicate struct {
	class        string
	annKey       string
	defaultClass string
}

var _ predicate.Predicate = &FilterByClassPredicate{}

// NewFilterByClassPredicate return a new ClassPredicate.
// class param represents the class of predicate.
// annKey represents the class annotation key.
func NewFilterByClassPredicate(class, annKey string) *FilterByClassPredicate {
	return &FilterByClassPredicate{
		class:        class,
		annKey:       annKey,
		defaultClass: "",
	}
}

// WithDefaultClass represents the default value for the class annotation, used when the annotation value is empty.
func (p *FilterByClassPredicate) WithDefaultClass(defaultClass string) {
	p.defaultClass = defaultClass
}

func (p *FilterByClassPredicate) matchesClass(m metav1.Object) bool {
	annotations := m.GetAnnotations()

	class, exists := annotations[p.annKey]
	if !exists || class == "" {
		class = p.defaultClass
	}

	return p.class == class
}

// Create returns true if the Create event should be processed.
func (p *FilterByClassPredicate) Create(e event.CreateEvent) bool {
	return p.matchesClass(e.Object)
}

// Delete returns true if the Delete event should be processed.
func (p *FilterByClassPredicate) Delete(e event.DeleteEvent) bool {
	return p.matchesClass(e.Object)
}

// Update returns true if the Update event should be processed.
func (p *FilterByClassPredicate) Update(e event.UpdateEvent) bool {
	return p.matchesClass(e.ObjectNew)
}

// Generic returns true if the Generic event should be processed.
func (p *FilterByClassPredicate) Generic(_ event.GenericEvent) bool {
	return true
}
