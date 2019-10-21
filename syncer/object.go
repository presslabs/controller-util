package syncer

import (
	"context"
	"fmt"

	"github.com/go-test/deep"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var (
	errOwnerDeleted = fmt.Errorf("owner is deleted")

	// ErrIgnore when returned the syncer ignores it and returns nil
	ErrIgnore = fmt.Errorf("ignored error")
)

// ObjectSyncer is a syncer.Interface for syncing kubernetes.Objects only by
// passing a SyncFn
type ObjectSyncer struct {
	Owner          runtime.Object
	Obj            runtime.Object
	SyncFn         controllerutil.MutateFn
	Name           string
	Client         client.Client
	Scheme         *runtime.Scheme
	previousObject runtime.Object
}

// objectWithoutSecretData returns the object without secretData
func objectWithoutSecretData(obj runtime.Object) runtime.Object {
	// if obj is secret, don't print secret data
	s, ok := obj.(*corev1.Secret)
	if ok {
		s.Data = nil
		s.StringData = nil

		return s
	}

	return obj
}

// GetObject returns the ObjectSyncer subject
func (s *ObjectSyncer) GetObject() interface{} {
	return objectWithoutSecretData(s.Obj)
}

// GetPreviousObject returns the ObjectSyncer previous subject
func (s *ObjectSyncer) GetPreviousObject() interface{} {
	return objectWithoutSecretData(s.previousObject)
}

// GetObjectType returns the type of the ObjectSyncer subject
func (s *ObjectSyncer) GetObjectType() string {
	return fmt.Sprintf("%T", s.Obj)
}

// GetOwner returns the ObjectSyncer owner
func (s *ObjectSyncer) GetOwner() runtime.Object {
	return s.Owner
}

// GetOwnerType returns the type of the ObjectSyncer owner
func (s *ObjectSyncer) GetOwnerType() string {
	return fmt.Sprintf("%T", s.Owner)
}

// Sync does the actual syncing and implements the syncer.Inteface Sync method
func (s *ObjectSyncer) Sync(ctx context.Context) (SyncResult, error) {
	result := SyncResult{}

	key, err := getKey(s.Obj)
	if err != nil {
		return result, err
	}

	result.Operation, err = controllerutil.CreateOrUpdate(ctx, s.Client, s.Obj, s.mutateFn())

	// check deep diff
	diff := deep.Equal(s.GetPreviousObject(), s.GetObject())

	// don't pass to user error for owner deletion, just don't create the object
	// nolint: gocritic
	if err == errOwnerDeleted {
		log.Info(string(result.Operation), "key", key, "kind", s.GetObjectType(), "error", err)
		err = nil
	} else if err == ErrIgnore {
		log.V(1).Info("syncer skipped", "key", key, "kind", s.GetObjectType())
		err = nil
	} else if err != nil {
		result.SetEventData(eventWarning, basicEventReason(s.Name, err),
			fmt.Sprintf("%s %s failed syncing: %s", s.GetObjectType(), key, err))
		log.Error(err, string(result.Operation), "key", key, "kind", s.GetObjectType(), "diff", diff)
	} else {
		result.SetEventData(eventNormal, basicEventReason(s.Name, err),
			fmt.Sprintf("%s %s %s successfully", s.GetObjectType(), key, result.Operation))
		log.V(1).Info(string(result.Operation), "key", key, "kind", s.GetObjectType(), "diff", diff)
	}

	return result, err
}

// Given an ObjectSyncer, returns a controllerutil.MutateFn which also sets the
// owner reference if the subject has one
func (s *ObjectSyncer) mutateFn() controllerutil.MutateFn {
	return func() error {
		s.previousObject = s.Obj.DeepCopyObject()
		err := s.SyncFn()
		if err != nil {
			return err
		}
		if s.Owner != nil {
			existingMeta, ok := s.Obj.(metav1.Object)
			if !ok {
				return fmt.Errorf("%s is not a metav1.Object", s.GetObjectType())
			}
			ownerMeta, ok := s.Owner.(metav1.Object)
			if !ok {
				return fmt.Errorf("%s is not a metav1.Object", s.GetOwnerType())
			}

			// set owner reference only if owner resource is not being deleted, otherwise the owner
			// reference will be reset in case of deleting with cascade=false.
			if ownerMeta.GetDeletionTimestamp().IsZero() {
				err := controllerutil.SetControllerReference(ownerMeta, existingMeta, s.Scheme)
				if err != nil {
					return err
				}
			} else if ctime := existingMeta.GetCreationTimestamp(); ctime.IsZero() {
				// the owner is deleted, don't recreate the resource if does not exist, because gc
				// will not delete it again because has no owner reference set
				return errOwnerDeleted
			}
		}
		return nil
	}
}

// NewObjectSyncer creates a new kubernetes.Object syncer for a given object
// with an owner and persists data using controller-runtime's CreateOrUpdate.
// The name is used for logging and event emitting purposes and should be an
// valid go identifier in upper camel case. (eg. MysqlStatefulSet)
func NewObjectSyncer(name string, owner, obj runtime.Object, c client.Client, scheme *runtime.Scheme,
	syncFn controllerutil.MutateFn) Interface {
	return &ObjectSyncer{
		Owner:  owner,
		Obj:    obj,
		SyncFn: syncFn,
		Name:   name,
		Client: c,
		Scheme: scheme,
	}
}
