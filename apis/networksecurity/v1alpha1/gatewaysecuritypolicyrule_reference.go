package v1alpha1

import (
	"context"
	"encoding/json"

	"github.com/GoogleCloudPlatform/k8s-config-connector/apis/common/identity"
	refs "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ refs.Ref = &NetworkSecurityGatewaySecurityPolicyRuleRef{}

// NetworkSecurityGatewaySecurityPolicyRuleRef represents a reference to a NetworkSecurityGatewaySecurityPolicyRule.
type NetworkSecurityGatewaySecurityPolicyRuleRef struct {
	// A reference to an externally managed NetworkSecurityGatewaySecurityPolicyRule resource.
	// Should be in the format "projects/{projectID}/locations/{location}/gatewaySecurityPolicies/{gateway_security_policy}/rules/{rule}".
	External string `json:"external,omitempty"`

	// The name of a NetworkSecurityGatewaySecurityPolicyRule resource.
	Name string `json:"name,omitempty"`

	// The namespace of a NetworkSecurityGatewaySecurityPolicyRule resource.
	Namespace string `json:"namespace,omitempty"`
}

func init() {
	refs.Register(&NetworkSecurityGatewaySecurityPolicyRuleRef{})
}

func (r *NetworkSecurityGatewaySecurityPolicyRuleRef) GetGVK() schema.GroupVersionKind {
	return NetworkSecurityGatewaySecurityPolicyRuleGVK
}

func (r *NetworkSecurityGatewaySecurityPolicyRuleRef) GetNamespacedName() client.ObjectKey {
	return client.ObjectKey{Name: r.Name, Namespace: r.Namespace}
}

func (r *NetworkSecurityGatewaySecurityPolicyRuleRef) GetExternal() string {
	return r.External
}

func (r *NetworkSecurityGatewaySecurityPolicyRuleRef) SetExternal(external string) {
	r.External = external
}

func (r *NetworkSecurityGatewaySecurityPolicyRuleRef) ValidateExternal(external string) error {
	id := &NetworkSecurityGatewaySecurityPolicyRuleIdentity{}
	return id.FromExternal(external)
}

func (r *NetworkSecurityGatewaySecurityPolicyRuleRef) ParseExternalToIdentity() (identity.Identity, error) {
	id := &NetworkSecurityGatewaySecurityPolicyRuleIdentity{}
	err := id.FromExternal(r.External)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (r *NetworkSecurityGatewaySecurityPolicyRuleRef) Normalize(ctx context.Context, reader client.Reader, namespace string) error {
	fallback := func(u *unstructured.Unstructured) string {
		spec := NetworkSecurityGatewaySecurityPolicyRuleSpec{}
		if specMap, found, _ := unstructured.NestedMap(u.Object, "spec"); found {
			specBytes, _ := json.Marshal(specMap)
			json.Unmarshal(specBytes, &spec)
		}
		id, err := getIdentityFromNetworkSecurityGatewaySecurityPolicyRuleSpec(ctx, reader, u, spec)
		if err != nil {
			return ""
		}
		return id.String()
	}
	return refs.NormalizeWithFallback(ctx, reader, r, namespace, fallback)
}
