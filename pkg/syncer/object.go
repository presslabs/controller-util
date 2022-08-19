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

package syncer

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-test/deep"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// ObjectSyncer is a syncer.Interface for syncing kubernetes.Objects only by
// passing a SyncFn.
type ObjectSyncer struct {
	Owner          client.Object
	Obj            client.Object
	SyncFn         controllerutil.MutateFn
	Name           string
	Client         client.Client
	previousObject runtime.Object
}

// objectType returns the type of a runtime.Object.
func (s *ObjectSyncer) objectType(obj runtime.Object) string {
	if obj != nil {
		gvk, err := apiutil.GVKForObject(obj, s.Client.Scheme())
		if err != nil {
			return fmt.Sprintf("%T", obj)
		}

		return gvk.String()
	}

	return "nil"
}

// Object returns the ObjectSyncer subject.
func (s *ObjectSyncer) Object() interface{} {
	return s.Obj
}

// ObjectOwner returns the ObjectSyncer owner.
func (s *ObjectSyncer) ObjectOwner() runtime.Object {
	return s.Owner
}

// Sync does the actual syncing and implements the syncer.Inteface Sync method.
func (s *ObjectSyncer) Sync(ctx context.Context) (SyncResult, error) {
	var err error

	result := SyncResult{}
	log := logf.FromContext(ctx, "syncer", s.Name)
	key := client.ObjectKeyFromObject(s.Obj)

	result.Operation, err = controllerutil.CreateOrUpdate(ctx, s.Client, s.Obj, s.mutateFn())

	// check deep diff
	diff := deep.Equal(redact(s.previousObject), redact(s.Obj))

	// don't pass to user error for owner deletion, just don't create the object
	//nolint: gocritic
	if errors.Is(err, ErrOwnerDeleted) {
		log.Info(string(result.Operation), "key", key, "kind", s.objectType(s.Obj), "error", err)
		err = nil
	} else if errors.Is(err, ErrIgnore) {
		log.V(1).Info("syncer skipped", "key", key, "kind", s.objectType(s.Obj), "error", err)
		err = nil
	} else if err != nil {
		result.SetEventData(eventWarning, basicEventReason(s.Name, err),
			fmt.Sprintf("%s %s failed syncing: %s", s.objectType(s.Obj), key, err))
		log.Error(err, string(result.Operation), "key", key, "kind", s.objectType(s.Obj), "diff", diff)
	} else {
		result.SetEventData(eventNormal, basicEventReason(s.Name, err),
			fmt.Sprintf("%s %s %s successfully", s.objectType(s.Obj), key, result.Operation))
		log.V(1).Info(string(result.Operation), "key", key, "kind", s.objectType(s.Obj), "diff", diff)
	}

	return result, err
}

// Given an ObjectSyncer, returns a controllerutil.MutateFn which also sets the
// owner reference if the subject has one.
func (s *ObjectSyncer) mutateFn() controllerutil.MutateFn {
	return func() error {
		s.previousObject = s.Obj.DeepCopyObject()

		err := s.SyncFn()
		if err != nil {
			return err
		}

		if s.Owner == nil {
			return nil
		}

		// set owner reference only if owner resource is not being deleted, otherwise the owner
		// reference will be reset in case of deleting with cascade=false.
		if s.Owner.GetDeletionTimestamp().IsZero() {
			if err := controllerutil.SetControllerReference(s.Owner, s.Obj, s.Client.Scheme()); err != nil {
				return err
			}
		} else if ctime := s.Obj.GetCreationTimestamp(); ctime.IsZero() {
			// the owner is deleted, don't recreate the resource if does not exist, because gc
			// will not delete it again because has no owner reference set
			return ErrOwnerDeleted
		}

		return nil
	}
}

// NewObjectSyncer creates a new kubernetes.Object syncer for a given object
// with an owner and persists data using controller-runtime's CreateOrUpdate.
// The name is used for logging and event emitting purposes and should be an
// valid go identifier in upper camel case. (eg. MysqlStatefulSet).
func NewObjectSyncer(name string, owner, obj client.Object, c client.Client, syncFn controllerutil.MutateFn) Interface {
	return &ObjectSyncer{
		Owner:  owner,
		Obj:    obj,
		SyncFn: syncFn,
		Name:   name,
		Client: c,
	}
}
