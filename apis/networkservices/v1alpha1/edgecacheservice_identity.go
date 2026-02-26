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

	"github.com/GoogleCloudPlatform/k8s-config-connector/apis/common"
	"github.com/GoogleCloudPlatform/k8s-config-connector/apis/common/identity"
	refsv1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/gcpurls"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// EdgeCacheServiceIdentityURL is the format for the externalRef of a EdgeCacheService.
	EdgeCacheServiceIdentityURL = "projects/{project}/locations/global/edgeCacheServices/{edgeCacheServiceID}"
)

var _ identity.Identity = &EdgeCacheServiceIdentity{}

var edgeCacheServiceURL = gcpurls.Template[EdgeCacheServiceIdentity](
	"networkservices.googleapis.com",
	EdgeCacheServiceIdentityURL,
)

// EdgeCacheServiceIdentity represents the identity of a EdgeCacheService.
// +k8s:deepcopy-gen=false
type EdgeCacheServiceIdentity struct {
	Project           string
	EdgeCacheServiceID string
}

func (i *EdgeCacheServiceIdentity) String() string {
	return edgeCacheServiceURL.ToString(*i)
}

func (i *EdgeCacheServiceIdentity) ID() string {
	return i.EdgeCacheServiceID
}

func (i *EdgeCacheServiceIdentity) FromExternal(ref string) error {
	out, match, err := edgeCacheServiceURL.Parse(ref)
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("format of EdgeCacheService external=%q was not known (use %s)", ref, edgeCacheServiceURL.CanonicalForm())
	}
	*i = *out
	return nil
}

func (i *EdgeCacheServiceIdentity) Parent() *EdgeCacheServiceParent {
	return &EdgeCacheServiceParent{
		ProjectID: i.Project,
	}
}

type EdgeCacheServiceParent struct {
	ProjectID string
}

func (p *EdgeCacheServiceParent) String() string {
	return fmt.Sprintf("projects/%s/locations/global", p.ProjectID)
}

// NewEdgeCacheServiceIdentity builds a EdgeCacheServiceIdentity from the Config Connector EdgeCacheService object.
func NewEdgeCacheServiceIdentity(ctx context.Context, reader client.Reader, obj *NetworkServicesEdgeCacheService) (*EdgeCacheServiceIdentity, error) {

	// Get Parent
	projectRef, err := refsv1beta1.ResolveProject(ctx, reader, obj.GetNamespace(), obj.Spec.ProjectRef)
	if err != nil {
		return nil, err
	}
	projectID := projectRef.ProjectID
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

	specIdentity := &EdgeCacheServiceIdentity{
		Project:           projectID,
		EdgeCacheServiceID: resourceID,
	}

	// Validate against the ID stored in status.externalRef
	var statusIdentity *EdgeCacheServiceIdentity
	externalRef := common.ValueOf(obj.Status.ExternalRef)
	if externalRef != "" {
		statusIdentity = &EdgeCacheServiceIdentity{}
		if err := statusIdentity.FromExternal(externalRef); err != nil {
			return nil, fmt.Errorf("cannot parse existing externalRef=%q: %w", externalRef, err)
		}
	}

	if statusIdentity != nil && statusIdentity.String() != specIdentity.String() {
		return nil, fmt.Errorf("existing externalRef=%q does not match the identity resolved from spec: %q", externalRef, specIdentity.String())
	}

	return specIdentity, nil
}
