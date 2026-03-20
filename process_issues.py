import json
import subprocess
import os

def run_cmd(cmd):
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    if result.returncode != 0:
        print(f"Command failed: {cmd}\n{result.stderr}")
    return result.stdout.strip()

# Load candidates from previous step
candidates_str = """
billingbudgets BillingBudgetsBudget
binaryauthorization BinaryAuthorizationAttestor
binaryauthorization BinaryAuthorizationPolicy
cloudfunctions CloudFunctionsFunction
cloudscheduler CloudSchedulerJob
configcontroller ConfigControllerInstance
containeranalysis ContainerAnalysisNote
datafusion DataFusionInstance
dataproc DataprocAutoscalingPolicy
dataproc DataprocCluster
dataproc DataprocWorkflowTemplate
dlp DLPDeidentifyTemplate
dlp DLPInspectTemplate
dlp DLPJobTrigger
dlp DLPStoredInfoType
eventarc EventarcTrigger
filestore FilestoreBackup
filestore FilestoreInstance
identityplatform IdentityPlatformConfig
identityplatform IdentityPlatformOAuthIDPConfig
identityplatform IdentityPlatformTenant
identityplatform IdentityPlatformTenantOAuthIDPConfig
networkconnectivity NetworkConnectivityHub
networkconnectivity NetworkConnectivitySpoke
networkservices NetworkServicesEndpointPolicy
networkservices NetworkServicesGRPCRoute
networkservices NetworkServicesGateway
networkservices NetworkServicesHTTPRoute
networkservices NetworkServicesMesh
networkservices NetworkServicesTCPRoute
networkservices NetworkServicesTLSRoute
osconfig OSConfigGuestPolicy
osconfig OSConfigOSPolicyAssignment
recaptchaenterprise RecaptchaEnterpriseKey
"""

candidates = [line.strip().split() for line in candidates_str.strip().split('\n')]
required_labels = {"overseer", "area/direct", "priority/medium"}

# Get all issues
issues_json = run_cmd('gh issue list --search "Create generate.sh and types.go files for" --state all --json number,title,labels,state --limit 1000')
issues = json.loads(issues_json) if issues_json else []

# Check how many are open
open_issues = [i for i in issues if i['state'].lower() == 'open']

# Track which candidates have issues
candidate_issue_map = {}
for i in issues:
    title = i['title']
    for g, k in candidates:
        # Title format check
        expected_title_lower = f"Create generate.sh and types.go files for {g} {k}".lower()
        if expected_title_lower in title.lower():
            candidate_issue_map[(g, k)] = i
            break

created_count = 0
for g, k in candidates:
    if (g, k) in candidate_issue_map:
        issue = candidate_issue_map[(g, k)]
        issue_labels = {l['name'] for l in issue['labels']}
        missing = required_labels - issue_labels
        if missing:
            print(f"Injecting labels {missing} into issue {issue['number']} for {g} {k}")
            run_cmd(f'gh issue edit {issue["number"]} --add-label "{",".join(missing)}"')
        continue

    # If no issue exists
    if len(open_issues) > 10:
        print(f"There are already {len(open_issues)} pending issues. Skipping creation of new issues.")
        break
    
    # Create issue
    print(f"Creating issue for {g} {k}...")
    
    # Need to form the body and title
    body = f"""As part of moving resources from terraform controllers to direct controllers (Epic #5954), we need to create the Go types for `{k}`.

Currently, `{k}` is managed by the Terraform controller (marked with `tf2crd=true`). The goal is to create the Go types in `apis/{g}/v1beta1/` so that we can eventually migrate the controller implementation to the "direct" approach.

### Instructions

1.  **Create a generate.sh**:
    Create `apis/{g}/v1beta1/generate.sh` which includes `{k}`.
    It likely maps to something like `google.cloud.{g}.v1`.
    Example:
    ```bash
    go run . generate-types \\
      --service google.cloud.{g}.v1 \\
      --api-version {g}.cnrm.cloud.google.com/v1beta1 \\
      --resource {k}:PolicyTag \\
      --include-skipped-output

    go run . generate-mapper \\
      --service google.cloud.{g}.v1 \\
      --api-version {g}.cnrm.cloud.google.com/v1beta1 \\
      --include-skipped-output
    ```

2.  Set the write permission on the new `apis/{g}/v1beta1/generate.sh` file. You should do this by running both `chmod +x apis/{g}/v1beta1/generate.sh` and `git add --chmod=+x apis/{g}/v1beta1/generate.sh`.

3.  **Generate Scaffolding**:
    Run `apis/{g}/v1beta1/generate.sh`. This should create `apis/{g}/v1beta1/{k.lower()}_types.go`.

4.  **Iterate on Types**:
    Compare the generated CRD with the existing one using `dev/tasks/diff-crds`.
    Modify `apis/{g}/v1beta1/{k.lower()}_types.go` until the CRD matches the existing one at `config/crds/resources/apiextensions.k8s.io_v1_customresourcedefinition_{k.lower()}s.{g}.cnrm.cloud.google.com.yaml`.

    **Acceptance Criteria:**
    - Running `dev/tasks/diff-crds` should not show differences (or minimal acceptable ones like descriptions).
    - Ensure that running the check_crd_equivalence MCP on the CRD should return EQUIVALENT.
    - Changes to the schema (fields added/removed) are NOT acceptable.

5.  **Copyright Headers**:
    Ensure that new files have the correct copyright header:
    ```go
    // Copyright 2026 Google LLC
    ```
    Please do not change the copyright on existing files.

6.  **Labels**:
    Ensure the controller-runtime annotations match the existing CRD labels, including:
    ```go
    // +kubebuilder:metadata:labels="cnrm.cloud.google.com/managed-by-kcc=true"
    // +kubebuilder:metadata:labels="cnrm.cloud.google.com/system=true"
    // +kubebuilder:metadata:labels="cnrm.cloud.google.com/stability-level=stable"
    // +kubebuilder:metadata:labels="cnrm.cloud.google.com/tf2crd=true"
    ```
    The goal is to maintain these annotations, not add an annotation if it is missing.

7.  **Status**:
    `status.observedGeneration` should be an `*int64`.

8. **Generate Mappers**:
   - Running `dev/tasks/generate-types-and-mappers` will generate the mapper code once the `apis/{g}/v1beta1/{k.lower()}_types.go` file is generating an equivalent CRD.
   - Run `make all-binary` to ensure the generated mapper code compiles. Please fix any issue discovered by this compilation.

This issue is part of Epic #5954.
"""
    title = f"Create generate.sh and types.go files for {g} {k}"
    with open('body.md', 'w') as f:
        f.write(body)
    
    cmd = f'gh issue create --title "{title}" --body-file body.md --label "overseer,area/direct,priority/medium" --milestone ""'
    print(f"Running command: {cmd}")
    # Wait we also need to make it a subtask of epic 5954 if possible? 
    # github doesn't have native subtasks, except by putting "This issue is part of Epic #5954." or commenting on the epic.
    # We will just create it.
    out = run_cmd(cmd)
    print(f"Created issue: {out}")
    open_issues.append({'title': title, 'state': 'open'}) # simulate the increase
    created_count += 1
    if created_count >= 1: # at most 1 issue
        break

if created_count == 0 and len(open_issues) <= 10:
    print("No new issues created (all candidates already have issues).")

