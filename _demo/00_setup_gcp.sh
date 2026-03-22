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

# Description: Sets up GCP project, APIs, Service Accounts, and GCS Bucket.

# Stop on any error
set -e

# --- Parameters ---
PROJECT_ID=${1:?Usage: $0 PROJECT_ID [REGION] [BUCKET_SUFFIX] [KCC_GSA_NAME] [MCL_GSA_NAME]}
REGION=${2:-us-central1}
BUCKET_SUFFIX=${3:-$(date +%s)}
KCC_GSA_NAME=${4:-kcc-ha-sa}
MCL_GSA_NAME=${5:-kcc-mcl-sa}

# --- Derived Variables ---
BUCKET_NAME="kcc-mcl-lease-${PROJECT_ID}-${BUCKET_SUFFIX}"
KCC_GSA_EMAIL="${KCC_GSA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"
MCL_GSA_EMAIL="${MCL_GSA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

echo "--- Configuration ---"
echo "PROJECT_ID: ${PROJECT_ID}"
echo "REGION: ${REGION}"
echo "BUCKET_NAME: ${BUCKET_NAME}"
echo "KCC_GSA_EMAIL: ${KCC_GSA_EMAIL}"
echo "MCL_GSA_EMAIL: ${MCL_GSA_EMAIL}"
echo "---------------------"

echo "--- 0.1. Setting up GCP Project and Enabling APIs ---"
gcloud config set project ${PROJECT_ID}

echo "Enabling necessary APIs..."
gcloud services enable container.googleapis.com \
    storage.googleapis.com \
    iam.googleapis.com \
    iamcredentials.googleapis.com \
    cloudresourcemanager.googleapis.com \
    cloudbuild.googleapis.com \
    artifactregistry.googleapis.com \
    --project=${PROJECT_ID}

echo "--- 0.2. Granting Cloud Build Permissions ---"
PROJECT_NUMBER=$(gcloud projects describe ${PROJECT_ID} --format='value(projectNumber)')
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com" \
    --role="roles/artifactregistry.writer" \
    --condition=None


echo "--- 0.3. Creating GCS Bucket ---"
if ! gcloud storage buckets describe gs://${BUCKET_NAME} &> /dev/null; then
  gcloud storage buckets create gs://${BUCKET_NAME} --location ${REGION} --project=${PROJECT_ID}
else
  echo "Bucket gs://${BUCKET_NAME} already exists."
fi

echo "--- 0.4. Creating Google Service Accounts ---"
if ! gcloud iam service-accounts describe ${KCC_GSA_EMAIL} --project=${PROJECT_ID} &> /dev/null; then
  echo "Creating KCC GSA: ${KCC_GSA_EMAIL}"
  gcloud iam service-accounts create ${KCC_GSA_NAME} --display-name="KCC HA Demo SA" --project=${PROJECT_ID}
else
  echo "KCC GSA ${KCC_GSA_EMAIL} already exists."
fi

if ! gcloud iam service-accounts describe ${MCL_GSA_EMAIL} --project=${PROJECT_ID} &> /dev/null; then
  echo "Creating MCL GSA: ${MCL_GSA_EMAIL}"
  gcloud iam service-accounts create ${MCL_GSA_NAME} --display-name="KCC MCL Demo SA" --project=${PROJECT_ID}
else
  echo "MCL GSA ${MCL_GSA_EMAIL} already exists."
fi

echo "Waiting for service accounts to propagate..."
sleep 15

echo "Granting KCC GSA project editor role..."
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:${KCC_GSA_EMAIL}" \
    --role="roles/editor" \
    --condition=None # Avoids issues if a condition exists

echo "Granting MCL GSA storage admin role on the bucket..."
gcloud storage buckets add-iam-policy-binding gs://${BUCKET_NAME} \
    --member="serviceAccount:${MCL_GSA_EMAIL}" \
    --role="roles/storage.admin"

echo "--- 0.5. Creating Artifact Registry Repository ---"
if ! gcloud artifacts repositories describe demo-repo --location=${REGION} --project=${PROJECT_ID} &> /dev/null; then
  echo "Creating Artifact Registry repository: demo-repo"
  gcloud artifacts repositories create demo-repo \
      --repository-format=docker \
      --location=${REGION} \
      --description="Docker repository for HA demo" \
      --project=${PROJECT_ID}
else
  echo "Artifact Registry repository demo-repo already exists."
fi

echo "--- GCP Setup Complete ---"
echo "Bucket Name: ${BUCKET_NAME}"
