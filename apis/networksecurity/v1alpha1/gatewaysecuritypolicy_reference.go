package v1alpha1

import (
	"context"

	"github.com/GoogleCloudPlatform/k8s-config-connector/apis/common/identity"
	refs "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ refs.Ref = &NetworkSecurityGatewaySecurityPolicyRef{}

// NetworkSecurityGatewaySecurityPolicyRef represents a reference to a NetworkSecurityGatewaySecurityPolicy.
type NetworkSecurityGatewaySecurityPolicyRef struct {
	// A reference to an externally managed NetworkSecurityGatewaySecurityPolicy resource.
	// Should be in the format "projects/{projectID}/locations/{location}/gatewaySecurityPolicies/{gateway_security_policy}".
	External string `json:"external,omitempty"`

	// The name of a NetworkSecurityGatewaySecurityPolicy resource.
	Name string `json:"name,omitempty"`

	// The namespace of a NetworkSecurityGatewaySecurityPolicy resource.
	Namespace string `json:"namespace,omitempty"`
}

func init() {
	refs.Register(&NetworkSecurityGatewaySecurityPolicyRef{})
}

func (r *NetworkSecurityGatewaySecurityPolicyRef) GetGVK() schema.GroupVersionKind {
	return NetworkSecurityGatewaySecurityPolicyGVK
}

func (r *NetworkSecurityGatewaySecurityPolicyRef) GetNamespacedName() client.ObjectKey {
	return client.ObjectKey{Name: r.Name, Namespace: r.Namespace}
}

func (r *NetworkSecurityGatewaySecurityPolicyRef) GetExternal() string {
	return r.External
}

func (r *NetworkSecurityGatewaySecurityPolicyRef) SetExternal(external string) {
	r.External = external
}

func (r *NetworkSecurityGatewaySecurityPolicyRef) ValidateExternal(external string) error {
	id := &NetworkSecurityGatewaySecurityPolicyIdentity{}
	return id.FromExternal(external)
}

func (r *NetworkSecurityGatewaySecurityPolicyRef) ParseExternalToIdentity() (identity.Identity, error) {
	id := &NetworkSecurityGatewaySecurityPolicyIdentity{}
	err := id.FromExternal(r.External)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (r *NetworkSecurityGatewaySecurityPolicyRef) Normalize(ctx context.Context, reader client.Reader, namespace string) error {
	fallback := func(u *unstructured.Unstructured) string {
		id, err := getIdentityFromNetworkSecurityGatewaySecurityPolicySpec(ctx, reader, u)
		if err != nil {
			return ""
		}
		return id.String()
	}
	return refs.NormalizeWithFallback(ctx, reader, r, namespace, fallback)
}
