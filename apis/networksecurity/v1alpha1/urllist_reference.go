package v1alpha1

import (
	"context"

	"github.com/GoogleCloudPlatform/k8s-config-connector/apis/common/identity"
	refs "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ refs.Ref = &NetworkSecurityUrlListRef{}

// NetworkSecurityUrlListRef represents a reference to a NetworkSecurityUrlList.
type NetworkSecurityUrlListRef struct {
	// A reference to an externally managed NetworkSecurityUrlList resource.
	// Should be in the format "projects/{projectID}/locations/{location}/urlLists/{url_list}".
	External string `json:"external,omitempty"`

	// The name of a NetworkSecurityUrlList resource.
	Name string `json:"name,omitempty"`

	// The namespace of a NetworkSecurityUrlList resource.
	Namespace string `json:"namespace,omitempty"`
}

func init() {
	refs.Register(&NetworkSecurityUrlListRef{})
}

func (r *NetworkSecurityUrlListRef) GetGVK() schema.GroupVersionKind {
	return NetworkSecurityUrlListGVK
}

func (r *NetworkSecurityUrlListRef) GetNamespacedName() client.ObjectKey {
	return client.ObjectKey{Name: r.Name, Namespace: r.Namespace}
}

func (r *NetworkSecurityUrlListRef) GetExternal() string {
	return r.External
}

func (r *NetworkSecurityUrlListRef) SetExternal(external string) {
	r.External = external
}

func (r *NetworkSecurityUrlListRef) ValidateExternal(external string) error {
	id := &NetworkSecurityUrlListIdentity{}
	return id.FromExternal(external)
}

func (r *NetworkSecurityUrlListRef) ParseExternalToIdentity() (identity.Identity, error) {
	id := &NetworkSecurityUrlListIdentity{}
	err := id.FromExternal(r.External)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (r *NetworkSecurityUrlListRef) Normalize(ctx context.Context, reader client.Reader, namespace string) error {
	fallback := func(u *unstructured.Unstructured) string {
		id, err := getIdentityFromNetworkSecurityUrlListSpec(ctx, reader, u)
		if err != nil {
			return ""
		}
		return id.String()
	}
	return refs.NormalizeWithFallback(ctx, reader, r, namespace, fallback)
}
