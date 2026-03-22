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

# Description: Cleans up all created GCP and Kubernetes resources.

# --- Parameters ---
PROJECT_ID=${1:?Usage: $0 PROJECT_ID BUCKET_NAME [REGION] [C1_NAME] [C2_NAME] [KCC_GSA_NAME] [MCL_GSA_NAME]}
BUCKET_NAME=${2:?Bucket name required}
REGION=${3:-us-central1}
C1_NAME=${4:-kcc-ha-cluster-1}
C2_NAME=${5:-kcc-ha-cluster-2}
KCC_GSA_NAME=${6:-kcc-ha-sa}
MCL_GSA_NAME=${7:-kcc-mcl-sa}

# --- Derived Variables ---
ZONE="${REGION}-c"
KCC_GSA_EMAIL="${KCC_GSA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"
MCL_GSA_EMAIL="${MCL_GSA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"
TENANT_NS="tenant-a" # Hardcoded for cleanup simplicity

echo "--- Configuration for Cleanup ---"
# ... (echo parameters)
echo "---------------------"

gcloud config set project ${PROJECT_ID}

read -p "ARE YOU SURE you want to delete clusters ${C1_NAME}, ${C2_NAME}, bucket ${BUCKET_NAME}, and associated resources in project ${PROJECT_ID}? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    exit 1
fi

echo "--- 7. Cleaning up resources ---"

echo "Deleting GKE Clusters..."
gcloud container clusters delete ${C1_NAME} --zone ${ZONE} --project=${PROJECT_ID} --quiet || echo "Failed to delete ${C1_NAME}"
gcloud container clusters delete ${C2_NAME} --zone ${ZONE} --project=${PROJECT_ID} --quiet || echo "Failed to delete ${C2_NAME}"

echo "Deleting GCS Bucket..."
gcloud storage rm -r gs://${BUCKET_NAME} || echo "Failed to delete bucket ${BUCKET_NAME}"

echo "Deleting Artifact Registry Repository..."
gcloud artifacts repositories delete demo-repo --location=${REGION} --project=${PROJECT_ID} --quiet || echo "Failed to delete Artifact Registry repository"

echo "Removing IAM policy bindings..."
gcloud projects remove-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:${KCC_GSA_EMAIL}" \
    --role="roles/editor" --all || echo "Failed to remove KCC editor binding"

KCC_KSA="cnrm-controller-manager-${TENANT_NS}"
gcloud iam service-accounts remove-iam-policy-binding ${KCC_GSA_EMAIL} \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:${PROJECT_ID}.svc.id.goog[cnrm-system/${KCC_KSA}]" \
    --project ${PROJECT_ID} --all || echo "Failed to remove KCC WI binding for ${KCC_KSA}"

gcloud iam service-accounts remove-iam-policy-binding ${MCL_GSA_EMAIL} \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:${PROJECT_ID}.svc.id.goog[multiclusterlease-system/default]" \
    --project ${PROJECT_ID} --all || echo "Failed to remove MCL WI binding"

echo "Deleting Google Service Accounts..."
gcloud iam service-accounts delete ${KCC_GSA_EMAIL} --project=${PROJECT_ID} --quiet || echo "Failed to delete ${KCC_GSA_EMAIL}"
gcloud iam service-accounts delete ${MCL_GSA_EMAIL} --project=${PROJECT_ID} --quiet || echo "Failed to delete ${MCL_GSA_EMAIL}"

echo "Cleanup complete."
