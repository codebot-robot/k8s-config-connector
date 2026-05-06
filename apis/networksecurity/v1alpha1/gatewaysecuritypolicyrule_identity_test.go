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
	"testing"
)

func TestNetworkSecurityGatewaySecurityPolicyRuleIdentity_FromExternal(t *testing.T) {
	cases := []struct {
		name     string
		external string
		expected NetworkSecurityGatewaySecurityPolicyRuleIdentity
		hasError bool
	}{
		{
			name:     "valid external",
			external: "projects/my-project/locations/us-central1/gatewaySecurityPolicies/my-policy/rules/my-rule",
			expected: NetworkSecurityGatewaySecurityPolicyRuleIdentity{
				Project:               "my-project",
				Location:              "us-central1",
				GatewaySecurityPolicy: "my-policy",
				Rule:                  "my-rule",
			},
			hasError: false,
		},
		{
			name:     "invalid external (missing prefix)",
			external: "my-policy/rules/my-rule",
			expected: NetworkSecurityGatewaySecurityPolicyRuleIdentity{},
			hasError: true,
		},
		{
			name:     "invalid external (wrong format)",
			external: "projects/my-project/locations/us-central1/gatewaySecurityPolicies/my-policy/wrongType/my-rule",
			expected: NetworkSecurityGatewaySecurityPolicyRuleIdentity{},
			hasError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var id NetworkSecurityGatewaySecurityPolicyRuleIdentity
			err := id.FromExternal(tc.external)
			if tc.hasError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if id != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, id)
				}
			}
		})
	}
}
