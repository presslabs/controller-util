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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type externalSyncer struct {
	name   string
	obj    interface{}
	owner  runtime.Object
	syncFn func(context.Context, interface{}) (controllerutil.OperationResult, error)
}

func (s *externalSyncer) Object() interface{} {
	return s.obj
}

func (s *externalSyncer) ObjectType() string {
	return fmt.Sprintf("%T", s.obj)
}

func (s *externalSyncer) ObjectOwner() runtime.Object {
	return s.owner
}

func (s *externalSyncer) Sync(ctx context.Context) (SyncResult, error) {
	var err error

	result := SyncResult{}
	result.Operation, err = s.syncFn(ctx, s.obj)

	if err != nil {
		result.SetEventData(eventWarning, basicEventReason(s.name, err),
			fmt.Sprintf("%s failed syncing: %s", s.ObjectType(), err))
		log.Error(err, string(result.Operation), "kind", s.ObjectType())
	} else {
		result.SetEventData(eventNormal, basicEventReason(s.name, err),
			fmt.Sprintf("%s successfully %s", s.ObjectType(), result.Operation))
		log.V(1).Info(string(result.Operation), "kind", s.ObjectType())
	}

	return result, err
}

// NewExternalSyncer creates a new syncer which syncs a generic object
// persisting it's state into and external store The name is used for logging
// and event emitting purposes and should be an valid go identifier in upper
// camel case. (eg. GiteaRepo).
func NewExternalSyncer(name string, owner runtime.Object, obj interface{},
	syncFn func(context.Context, interface{}) (controllerutil.OperationResult, error)) Interface {
	return &externalSyncer{
		name:   name,
		obj:    obj,
		owner:  owner,
		syncFn: syncFn,
	}
}
