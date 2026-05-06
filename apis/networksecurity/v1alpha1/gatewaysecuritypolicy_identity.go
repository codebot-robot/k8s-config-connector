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

	"github.com/GoogleCloudPlatform/k8s-config-connector/apis/common/identity"
	refs "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/gcpurls"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	_ identity.Identity   = &NetworkSecurityGatewaySecurityPolicyIdentity{}
	_ identity.IdentityV2 = &NetworkSecurityGatewaySecurityPolicyIdentity{}
	_ identity.Resource   = &NetworkSecurityGatewaySecurityPolicy{}
)

var NetworkSecurityGatewaySecurityPolicyIdentityFormat = gcpurls.Template[NetworkSecurityGatewaySecurityPolicyIdentity]("networksecurity.googleapis.com", "projects/{project}/locations/{location}/gatewaySecurityPolicies/{gatewaySecurityPolicy}")

// NetworkSecurityGatewaySecurityPolicyIdentity represents the identity for a NetworkSecurityGatewaySecurityPolicy.
// +k8s:deepcopy-gen=false
type NetworkSecurityGatewaySecurityPolicyIdentity struct {
	Project               string
	Location              string
	GatewaySecurityPolicy string
}

func (i *NetworkSecurityGatewaySecurityPolicyIdentity) String() string {
	return NetworkSecurityGatewaySecurityPolicyIdentityFormat.ToString(*i)
}

func (i *NetworkSecurityGatewaySecurityPolicyIdentity) FromExternal(ref string) error {
	parsed, match, err := NetworkSecurityGatewaySecurityPolicyIdentityFormat.Parse(ref)
	if err != nil {
		return fmt.Errorf("format of NetworkSecurityGatewaySecurityPolicy external=%q was not known (use %s): %w", ref, NetworkSecurityGatewaySecurityPolicyIdentityFormat.CanonicalForm(), err)
	}
	if !match {
		return fmt.Errorf("format of NetworkSecurityGatewaySecurityPolicy external=%q was not known (use %s)", ref, NetworkSecurityGatewaySecurityPolicyIdentityFormat.CanonicalForm())
	}
	*i = *parsed
	return nil
}

func (i *NetworkSecurityGatewaySecurityPolicyIdentity) Host() string {
	return NetworkSecurityGatewaySecurityPolicyIdentityFormat.Host()
}

// GetIdentity implements identity.Resource
func (r *NetworkSecurityGatewaySecurityPolicy) GetIdentity(ctx context.Context, reader client.Reader) (identity.Identity, error) {
	id, err := getIdentityFromNetworkSecurityGatewaySecurityPolicySpec(ctx, reader, r)
	if err != nil {
		return nil, err
	}
	if r.Status.ExternalRef != nil {
		statusID := &NetworkSecurityGatewaySecurityPolicyIdentity{}
		if err := statusID.FromExternal(*r.Status.ExternalRef); err == nil {
			if statusID.Project != id.Project || statusID.Location != id.Location || statusID.GatewaySecurityPolicy != id.GatewaySecurityPolicy {
				return nil, fmt.Errorf("identity in status.externalRef %q does not match identity from spec: %q", statusID.String(), id.String())
			}
		}
	}
	return id, nil
}

func getIdentityFromNetworkSecurityGatewaySecurityPolicySpec(ctx context.Context, reader client.Reader, obj client.Object) (*NetworkSecurityGatewaySecurityPolicyIdentity, error) {
	resourceID, err := refs.GetResourceID(obj)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve resource ID")
	}

	location, err := refs.GetLocation(obj)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve location")
	}

	projectID, err := refs.ResolveProjectID(ctx, reader, obj)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve project")
	}

	return &NetworkSecurityGatewaySecurityPolicyIdentity{
		Project:               projectID,
		Location:              location,
		GatewaySecurityPolicy: resourceID,
	}, nil
}
