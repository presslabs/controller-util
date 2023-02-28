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

package syncer

import (
	"context"
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// RemoveResourceSyncer is a syncer.Interface for deleting kubernetes.Objects.
type RemoveResourceSyncer struct {
	Owner  client.Object
	Obj    client.Object
	Name   string
	Client client.Client
}

// Object returns the ObjectSyncer subject.
func (s *RemoveResourceSyncer) Object() interface{} {
	return s.Obj
}

// ObjectOwner returns the ObjectSyncer owner.
func (s *RemoveResourceSyncer) ObjectOwner() runtime.Object {
	return s.Owner
}

// Sync does the actual syncing and implements the syncer.Inteface Sync method.
func (s *RemoveResourceSyncer) Sync(ctx context.Context) (SyncResult, error) {
	result := SyncResult{}
	log := logf.FromContext(ctx, "syncer", s.Name)
	key := client.ObjectKeyFromObject(s.Obj)

	result.Operation = controllerutil.OperationResultNone

	// fetch the resource
	if err := s.Client.Get(ctx, key, s.Obj); err != nil {
		if k8serrors.IsNotFound(err) {
			return result, nil
		}

		log.Error(err, string(result.Operation), "key", key, "kind", objectType(s.Obj, s.Client))

		return result, fmt.Errorf("error when fetching resource: %w", err)
	}

	// delete the resource
	if err := s.Client.Delete(ctx, s.Obj); err != nil {
		log.Error(err, string(result.Operation), "key", key, "kind", objectType(s.Obj, s.Client))

		return result, fmt.Errorf("error when deleting resource: %w", err)
	}

	result.Operation = controllerutil.OperationResult("deleted")
	result.SetEventData(eventNormal, basicEventReason(s.Name, nil), fmt.Sprintf("%s %s successfully deleted", objectType(s.Obj, s.Client), key))

	log.V(1).Info(string(result.Operation), "key", key, "kind", objectType(s.Obj, s.Client))

	return result, nil
}

// NewRemoveResourceSyncer creates a new kubernetes.Object syncer for a given object
// with an owner and persists data using controller-runtime's Delete.
// The name is used for logging and event emitting purposes and should be an
// valid go identifier in upper camel case. (eg. MysqlStatefulSet).
func NewRemoveResourceSyncer(name string, owner, obj client.Object, c client.Client) Interface {
	return &RemoveResourceSyncer{
		Owner:  owner,
		Obj:    obj,
		Name:   name,
		Client: c,
	}
}
