#!/bin/bash
# Copyright 2026 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"

# Generate v1beta1 mappers
echo "Generating v1beta1..."
${REPO_ROOT}/apis/bigtable/v1beta1/generate.sh
mv ${REPO_ROOT}/pkg/controller/direct/bigtable/mapper.generated.go ${REPO_ROOT}/pkg/controller/direct/bigtable/mapper_v1beta1.generated.go

# Generate v1alpha1 mappers
echo "Generating v1alpha1..."
${REPO_ROOT}/apis/bigtable/v1alpha1/generate.sh
mv ${REPO_ROOT}/pkg/controller/direct/bigtable/mapper.generated.go ${REPO_ROOT}/pkg/controller/direct/bigtable/mapper_v1alpha1.generated.go

echo "Done. Generated mapper_v1beta1.generated.go and mapper_v1alpha1.generated.go"
