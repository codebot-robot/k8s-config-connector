// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"context"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-config-connector/apis/common"
	"github.com/GoogleCloudPlatform/k8s-config-connector/apis/common/identity"
	refsv1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/gcpurls"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	_ identity.IdentityV2 = &EdgeCacheServiceIdentity{}
	_ identity.Resource   = &NetworkServicesEdgeCacheService{}
)

var networkservicesEdgeCacheServiceGCPURL = gcpurls.Template[EdgeCacheServiceIdentity]("networkservices.googleapis.com", "projects/{project}/locations/global/edgeCacheServices/{resourceID}")
var networkservicesEdgeCacheServiceParentGCPURL = gcpurls.Template[EdgeCacheServiceParent]("networkservices.googleapis.com", "projects/{project}/locations/global")

// EdgeCacheServiceIdentity represents the identity of a NetworkServicesEdgeCacheService.
// +k8s:deepcopy-gen=false
type EdgeCacheServiceIdentity struct {
	Project    string
	ResourceID string
}

func (i *EdgeCacheServiceIdentity) String() string {
	return networkservicesEdgeCacheServiceGCPURL.ToString(*i)
}

func (i *EdgeCacheServiceIdentity) FromExternal(ref string) error {
	ref = strings.TrimPrefix(ref, "/")
	parsed, match, err := networkservicesEdgeCacheServiceGCPURL.Parse(ref)
	if err != nil {
		return fmt.Errorf("format of NetworkServicesEdgeCacheService external=%q was not known (use %s): %w", ref, networkservicesEdgeCacheServiceGCPURL.CanonicalForm(), err)
	}
	if !match {
		return fmt.Errorf("format of NetworkServicesEdgeCacheService external=%q was not known (use %s)", ref, networkservicesEdgeCacheServiceGCPURL.CanonicalForm())
	}
	*i = *parsed
	return nil
}

func (i *EdgeCacheServiceIdentity) Host() string {
	return networkservicesEdgeCacheServiceGCPURL.Host()
}

func (i *EdgeCacheServiceIdentity) Parent() *EdgeCacheServiceParent {
	return &EdgeCacheServiceParent{
		Project: i.Project,
	}
}

func (i *EdgeCacheServiceIdentity) ID() string {
	return i.ResourceID
}

type EdgeCacheServiceParent struct {
	Project string
}

func (p *EdgeCacheServiceParent) String() string {
	return networkservicesEdgeCacheServiceParentGCPURL.ToString(*p)
}

// GetIdentity builds an EdgeCacheServiceIdentity from the Config Connector EdgeCacheService object.
func (obj *NetworkServicesEdgeCacheService) GetIdentity(ctx context.Context, reader client.Reader) (identity.Identity, error) {
	// Get Parent
	projectRef := &refsv1beta1.ProjectRef{
		External:  obj.Spec.ProjectRef.External,
		Name:      obj.Spec.ProjectRef.Name,
		Namespace: obj.Spec.ProjectRef.Namespace,
	}
	projectIDRef, err := refsv1beta1.ResolveProject(ctx, reader, obj.GetNamespace(), projectRef)
	if err != nil {
		return nil, err
	}
	projectID := projectIDRef.ProjectID
	if projectID == "" {
		return nil, fmt.Errorf("cannot resolve project")
	}

	// Get desired ID
	resourceID := common.ValueOf(obj.Spec.ResourceID)
	if resourceID == "" {
		resourceID = obj.GetName()
	}
	if resourceID == "" {
		return nil, fmt.Errorf("cannot resolve resource ID")
	}

	return &EdgeCacheServiceIdentity{
		Project:    projectID,
		ResourceID: resourceID,
	}, nil
}

// NewEdgeCacheServiceIdentity is a helper to get the identity.
// It matches the older pattern used in direct controllers.
func NewEdgeCacheServiceIdentity(ctx context.Context, reader client.Reader, obj *NetworkServicesEdgeCacheService) (*EdgeCacheServiceIdentity, error) {
	id, err := obj.GetIdentity(ctx, reader)
	if err != nil {
		return nil, err
	}
	return id.(*EdgeCacheServiceIdentity), nil
}
