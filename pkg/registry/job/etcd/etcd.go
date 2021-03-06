/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package etcd

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/experimental"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/registry/generic"
	etcdgeneric "k8s.io/kubernetes/pkg/registry/generic/etcd"
	"k8s.io/kubernetes/pkg/registry/job"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/storage"
)

// rest implements a RESTStorage for jobs against etcd
type REST struct {
	*etcdgeneric.Etcd
}

// jobPrefix is the location for jobs in etcd, only exposed
// for testing
var jobPrefix = "/jobs"

// NewREST returns a RESTStorage object that will work against Jobs.
func NewREST(s storage.Interface) *REST {
	store := &etcdgeneric.Etcd{
		NewFunc: func() runtime.Object { return &experimental.Job{} },

		// NewListFunc returns an object capable of storing results of an etcd list.
		NewListFunc: func() runtime.Object { return &experimental.JobList{} },
		// Produces a path that etcd understands, to the root of the resource
		// by combining the namespace in the context with the given prefix
		KeyRootFunc: func(ctx api.Context) string {
			return etcdgeneric.NamespaceKeyRootFunc(ctx, jobPrefix)
		},
		// Produces a path that etcd understands, to the resource by combining
		// the namespace in the context with the given prefix
		KeyFunc: func(ctx api.Context, name string) (string, error) {
			return etcdgeneric.NamespaceKeyFunc(ctx, jobPrefix, name)
		},
		// Retrieve the name field of a job
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*experimental.Job).Name, nil
		},
		// Used to match objects based on labels/fields for list and watch
		PredicateFunc: func(label labels.Selector, field fields.Selector) generic.Matcher {
			return job.MatchJob(label, field)
		},
		EndpointName: "jobs",

		// Used to validate job creation
		CreateStrategy: job.Strategy,

		// Used to validate job updates
		UpdateStrategy: job.Strategy,

		Storage: s,
	}

	return &REST{store}
}
