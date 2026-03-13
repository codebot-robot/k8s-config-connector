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

package mockcompute

import (
	"context"

	pbv1beta "github.com/GoogleCloudPlatform/k8s-config-connector/mockgcp/generated/mockgcp/cloud/compute/v1beta"
)

type RegionalOperationsV1Beta struct {
	*MockService
	pbv1beta.UnimplementedRegionOperationsServer
}

func (s *RegionalOperationsV1Beta) Get(ctx context.Context, req *pbv1beta.GetRegionOperationRequest) (*pbv1beta.Operation, error) {
	fqn := s.regionalOperationFQN(req.Project, req.Region, req.Operation)

	return s.getOperationV1Beta(ctx, fqn)
}
