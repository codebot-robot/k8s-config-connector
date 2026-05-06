#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd ${REPO_ROOT}/dev/tools/controllerbuilder

./generate-proto.sh

go run . generate-types \
  --service google.cloud.networksecurity.v1 \
  --api-version networksecurity.cnrm.cloud.google.com/v1alpha1 \
  --resource NetworkSecurityUrlList:UrlList \
  --resource NetworkSecurityTLSInspectionPolicy:TlsInspectionPolicy \
  --resource NetworkSecurityGatewaySecurityPolicy:GatewaySecurityPolicy \
  --resource NetworkSecurityGatewaySecurityPolicyRule:GatewaySecurityPolicyRule

go run . generate-mapper \
  --service google.cloud.networksecurity.v1 \
  --api-version networksecurity.cnrm.cloud.google.com/v1alpha1

cd ${REPO_ROOT}
dev/tasks/generate-crds
