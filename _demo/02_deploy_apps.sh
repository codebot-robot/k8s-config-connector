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

# Description: Deploys KCC, MCL, Syncer and configures Workload Identity.

# Stop on any error
set -e

# --- Parameters ---
PROJECT_ID=${1:?Usage: $0 PROJECT_ID BUCKET_NAME [REGION] [C1_NAME] [C2_NAME] [KCC_GSA_NAME] [MCL_GSA_NAME] [TENANT_NS] [LEASE_NAME]}
BUCKET_NAME=${2:?Bucket name required}
REGION=${3:-us-central1}
C1_NAME=${4:-kcc-ha-cluster-1}
C2_NAME=${5:-kcc-ha-cluster-2}
KCC_GSA_NAME=${6:-kcc-ha-sa}
MCL_GSA_NAME=${7:-kcc-mcl-sa}
TENANT_NS=${8:-tenant-a}
LEASE_NAME=${9:-kcc-managers}

# --- Derived Variables ---
ZONE="${REGION}-c"
KCC_GSA_EMAIL="${KCC_GSA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"
MCL_GSA_EMAIL="${MCL_GSA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"
C1_CTX="gke_${PROJECT_ID}_${ZONE}_${C1_NAME}"
C2_CTX="gke_${PROJECT_ID}_${ZONE}_${C2_NAME}"
C1_ID="cluster-1"
C2_ID="cluster-2"

echo "--- Configuration ---"
# ... (echo parameters)
echo "---------------------"

gcloud config set project ${PROJECT_ID}
gcloud container clusters get-credentials ${C1_NAME} --zone ${ZONE} --project=${PROJECT_ID}
gcloud container clusters get-credentials ${C2_NAME} --zone ${ZONE} --project=${PROJECT_ID}

echo "--- 0.5. Preparing Repositories ---"
TMP_DIR=$(mktemp -d)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
echo "Working in ${TMP_DIR}..."
cd ${TMP_DIR}
echo "Copying local KCC repository from ${REPO_ROOT}..."
rsync -a --exclude=node_modules --exclude=.build --exclude=bin --exclude=artifactz --exclude="*.pb" ${REPO_ROOT}/ connector/
git clone https://github.com/gke-labs/multicluster-leader-election.git
git clone https://github.com/gke-labs/kube-etl.git

echo "--- 0.6. Creating Essential Namespaces ---"
for CTX in $C1_CTX $C2_CTX; do
  kubectl --context=$CTX create namespace cnrm-system || true
  kubectl --context=$CTX create namespace multiclusterlease-system || true
  kubectl --context=$CTX create namespace syncer-system || true
done

echo "--- 1. Building and Installing KCC Operator ---"
KCC_VERSION="1.146.0"
cd connector
# Disable BuildKit as it fails in some local environments
sed -i 's/DOCKER_BUILDKIT=1/DOCKER_BUILDKIT=0/g' Makefile

# Build and push the manager/controller image with the version tag expected by the operator
echo "Building local KCC controller image v${KCC_VERSION}..."
make docker-build-manager PROJECT_ID=${PROJECT_ID}
LOCAL_CONTROLLER_IMG=$(docker images --format "{{.Repository}}:{{.Tag}}" | grep "cnrm/controller" | head -n 1)
docker tag ${LOCAL_CONTROLLER_IMG} ${REGION}-docker.pkg.dev/${PROJECT_ID}/demo-repo/controller:${KCC_VERSION}
docker push ${REGION}-docker.pkg.dev/${PROJECT_ID}/demo-repo/controller:${KCC_VERSION}

# Build and push the operator itself
make config-connector-manifests-standard push-operator-manifest PROJECT_ID=${PROJECT_ID}
cd ..

for CTX in $C1_CTX $C2_CTX; do
  kubectl --context=$CTX apply -f connector/config/installbundle/release-manifests/standard/manifests.yaml
done
echo "Waiting for KCC operators to be ready..."
sleep 5
for CTX in $C1_CTX $C2_CTX; do
   kubectl --context=$CTX wait -n configconnector-operator-system --for=condition=Ready pod/configconnector-operator-0 --timeout=180s
done

echo "--- 1.1. Installing CRDs ---"
kubectl --context=$C1_CTX apply -f multicluster-leader-election/config/crd/bases/
kubectl --context=$C2_CTX apply -f multicluster-leader-election/config/crd/bases/
kubectl --context=$C1_CTX apply -f kube-etl/syncer/config/crd/
kubectl --context=$C2_CTX apply -f kube-etl/syncer/config/crd/

echo "--- 2. Configuring ConfigConnector (Namespaced Mode) ---"
cat <<EOF | kubectl --context=$C1_CTX apply -f -
apiVersion: core.cnrm.cloud.google.com/v1beta1
kind: ConfigConnector
metadata:
  name: configconnector.core.cnrm.cloud.google.com
spec:
  mode: namespaced
  experiments:
    multiClusterLease:
      leaseName: $LEASE_NAME
      namespace: cnrm-system
      clusterCandidateIdentity: $C1_ID
EOF
cat <<EOF | kubectl --context=$C2_CTX apply -f -
apiVersion: core.cnrm.cloud.google.com/v1beta1
kind: ConfigConnector
metadata:
  name: configconnector.core.cnrm.cloud.google.com
spec:
  mode: namespaced
  experiments:
    multiClusterLease:
      leaseName: $LEASE_NAME
      namespace: cnrm-system
      clusterCandidateIdentity: $C2_ID
EOF
for CTX in $C1_CTX $C2_CTX; do
  kubectl --context=$CTX create namespace $TENANT_NS || true
  cat <<EOF | kubectl --context=$CTX apply -f -
apiVersion: core.cnrm.cloud.google.com/v1beta1
kind: ConfigConnectorContext
metadata:
  name: configconnectorcontext.core.cnrm.cloud.google.com
  namespace: $TENANT_NS
spec:
  googleServiceAccount: $KCC_GSA_EMAIL
EOF
done

echo "--- 2.1. Binding Workload Identity for KCC ---"
echo "Waiting for KCC to create per-namespace service accounts and StatefulSets..."
sleep 20 # Wait for operator to process the context
for CTX in $C1_CTX $C2_CTX; do
  KCC_KSA="cnrm-controller-manager-${TENANT_NS}"
  if kubectl --context=$CTX get serviceaccount ${KCC_KSA} -n cnrm-system &> /dev/null; then
    gcloud iam service-accounts add-iam-policy-binding ${KCC_GSA_EMAIL} \
        --role roles/iam.workloadIdentityUser \
        --member "serviceAccount:${PROJECT_ID}.svc.id.goog[cnrm-system/${KCC_KSA}]" \
        --project ${PROJECT_ID} --condition=None
  else
    echo "Warning: KCC KSA ${KCC_KSA} not found in cnrm-system on ${CTX} yet."
  fi

  echo "Patching KCC Manager to use locally built custom image..."
  # Scale KCC operator to 0 to prevent it from overwriting our image patch
  kubectl --context=$CTX scale statefulset configconnector-operator -n configconnector-operator-system --replicas=0
  
  # Wait for StatefulSet to exist
  while ! kubectl --context=$CTX get statefulset -n cnrm-system -l cnrm.cloud.google.com/component=cnrm-controller-manager > /dev/null 2>&1; do sleep 2; done
  
  MGR_STS=$(kubectl --context=$CTX get statefulset -n cnrm-system -l cnrm.cloud.google.com/component=cnrm-controller-manager -o jsonpath='{.items[0].metadata.name}')
  
  # Patch the manager container image
  kubectl --context=$CTX patch statefulset $MGR_STS -n cnrm-system --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/image", "value": "'${REGION}-docker.pkg.dev/${PROJECT_ID}/demo-repo/controller:${KCC_VERSION}'"}]'
  
  # Delete any existing crashing pods to force recreation with new image
  kubectl --context=$CTX delete pod -n cnrm-system -l cnrm.cloud.google.com/component=cnrm-controller-manager --grace-period=0 --force || true
done

echo "--- 3. Applying Additional RBAC for KCC Manager ---"
cat <<EOF > kcc-mcl-syncer-rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cnrm-syncer-mcl-role
rules:
- apiGroups: ["syncer.gkelabs.io"]
  resources: ["krmsyncers", "krmsyncers/status"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: ["multicluster.gkelabs.io", "multicluster.core.cnrm.cloud.google.com"]
  resources: ["multiclusterleases"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cnrm-syncer-mcl-rolebinding-${TENANT_NS}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cnrm-syncer-mcl-role
subjects:
- kind: ServiceAccount
  name: cnrm-controller-manager-${TENANT_NS}
  namespace: cnrm-system
EOF
kubectl --context=$C1_CTX apply -f kcc-mcl-syncer-rbac.yaml
kubectl --context=$C2_CTX apply -f kcc-mcl-syncer-rbac.yaml

echo "--- 4. Deploy MCL and Syncer Controllers ---"
echo "Building MCL and Syncer Controller images (this may take a few minutes)..."
MCL_IMAGE="${REGION}-docker.pkg.dev/${PROJECT_ID}/demo-repo/multiclusterlease:latest"
SYNCER_IMAGE="${REGION}-docker.pkg.dev/${PROJECT_ID}/demo-repo/syncer:latest"

gcloud auth application-default set-quota-project ${PROJECT_ID}
gcloud auth configure-docker ${REGION}-docker.pkg.dev --quiet

# Build MCL
cd multicluster-leader-election
docker build -t ${MCL_IMAGE} .
docker push ${MCL_IMAGE}
cd ..

# Build Syncer
cd kube-etl/syncer
sed -i 's/FROM golang:1.24/FROM golang:1.26.1/g' Dockerfile
docker build -t ${SYNCER_IMAGE} .
docker push ${SYNCER_IMAGE}
cd ../..

# Deploy MCL
kubectl --context=$C1_CTX apply -k multicluster-leader-election/config/default
kubectl --context=$C2_CTX apply -k multicluster-leader-election/config/default

# Deploy Syncer
cat <<EOF > syncer-deployment.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: syncer-manager
  namespace: syncer-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: syncer-manager-role
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: syncer-manager-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: syncer-manager-role
subjects:
- kind: ServiceAccount
  name: syncer-manager
  namespace: syncer-system
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: syncer-manager
  namespace: syncer-system
spec:
  replicas: 1
  selector:
    matchLabels:
      control-plane: syncer-manager
  template:
    metadata:
      labels:
        control-plane: syncer-manager
    spec:
      serviceAccountName: syncer-manager
      containers:
      - name: manager
        image: ${SYNCER_IMAGE}
        args:
        - "--v=2"
EOF

for CTX in $C1_CTX $C2_CTX; do
  kubectl --context=$CTX apply -f syncer-deployment.yaml
done

echo "Binding Workload Identity GSA to MCL controller's KSA..."
for CTX in $C1_CTX $C2_CTX; do
    gcloud iam service-accounts add-iam-policy-binding ${MCL_GSA_EMAIL} \
        --role roles/iam.workloadIdentityUser \
        --member "serviceAccount:${PROJECT_ID}.svc.id.goog[multiclusterlease-system/default]" \
        --project ${PROJECT_ID} --condition=None
    kubectl --context=$CTX annotate serviceaccount default -n multiclusterlease-system iam.gke.io/gcp-service-account=${MCL_GSA_EMAIL} --overwrite
done

echo "Patching MCL deployments..."
PATCH_STR='[{"op": "replace", "path": "/spec/template/spec/containers/0/args", "value": ["--gcs-bucket='${BUCKET_NAME}'", "--metrics-addr=:8080"]}, {"op": "replace", "path": "/spec/template/spec/containers/0/image", "value": "'${MCL_IMAGE}'"}, {"op": "remove", "path": "/spec/template/spec/containers/0/volumeMounts"}, {"op": "remove", "path": "/spec/template/spec/volumes"}, {"op": "remove", "path": "/spec/template/spec/containers/0/env"}]'
for CTX in $C1_CTX $C2_CTX; do
  kubectl --context=$CTX patch deployment multiclusterlease-controller-manager -n multiclusterlease-system --type='json' -p="$PATCH_STR"
done
echo "Restarting MCL controller pods..."
kubectl --context=$C1_CTX delete pod -n multiclusterlease-system -l control-plane=controller-manager --grace-period=0 --force || true
kubectl --context=$C2_CTX delete pod -n multiclusterlease-system -l control-plane=controller-manager --grace-period=0 --force || true
sleep 15

echo "--- 5. Setup Cross-Cluster Authentication ---"
KUBECONFIG=~/.kube/config kubectl config view --raw --context=$C1_CTX > cluster1-kubeconfig
KUBECONFIG=~/.kube/config kubectl config view --raw --context=$C2_CTX > cluster2-kubeconfig
kubectl --context=$C1_CTX create secret generic $C2_ID --namespace cnrm-system --from-file=kubeconfig=cluster2-kubeconfig --dry-run -o yaml | kubectl --context=$C1_CTX apply -f -
kubectl --context=$C2_CTX create secret generic $C1_ID --namespace cnrm-system --from-file=kubeconfig=cluster1-kubeconfig --dry-run -o yaml | kubectl --context=$C2_CTX apply -f -
rm cluster1-kubeconfig cluster2-kubeconfig
cd ~

echo "--- Deployment Complete ---"
echo "TMP_DIR on this machine is: ${TMP_DIR} - run 'rm -rf ${TMP_DIR}' to clean up later."
