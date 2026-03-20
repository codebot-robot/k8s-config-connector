import json
import subprocess

def run_cmd(cmd):
    result = subprocess.run(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    return result.stdout.strip()

# Get all open issues to check pending count
open_issues_output = run_cmd("gh issue list --state open --search 'in:title \"Create generate.sh and types.go files for\"' --json title,number,labels --limit 100")
open_issues = json.loads(open_issues_output) if open_issues_output else []
pending_count = len(open_issues)
print(f"Pending issues count: {pending_count}")

# Get all issues (open and closed) to check existence
all_issues_output = run_cmd("gh issue list --state all --search 'in:title \"Create generate.sh and types.go files for\"' --json title,number,labels --limit 1000")
all_issues = json.loads(all_issues_output) if all_issues_output else []

# Build a map of "Group Kind" to issue
issue_map = {}
for issue in all_issues:
    title = issue['title']
    # Extract Group and Kind from title "Create generate.sh and types.go files for <Group> <Kind>"
    prefix = "Create generate.sh and types.go files for "
    if title.startswith(prefix):
        suffix = title[len(prefix):].strip()
        issue_map[suffix] = issue

with open("candidates.txt", "w") as f:
    f.write(run_cmd("python3 find_candidates_dcl2crd.py"))

with open("candidates.txt", "r") as f:
    candidates = [line.strip() for line in f if line.strip()]

required_labels = {"overseer", "area/direct", "priority/medium"}

new_issue_created = False

for candidate in candidates:
    if candidate in issue_map:
        issue = issue_map[candidate]
        existing_labels = {lbl["name"] for lbl in issue.get("labels", [])}
        missing_labels = required_labels - existing_labels
        if missing_labels:
            print(f"Injecting missing labels {missing_labels} for issue #{issue['number']} ({candidate})")
            labels_arg = ",".join(missing_labels)
            run_cmd(f"gh issue edit {issue['number']} --add-label {labels_arg}")
        else:
            print(f"Issue #{issue['number']} for {candidate} already has all required labels.")
    else:
        if pending_count > 10:
            print(f"There are already > 10 pending issues ({pending_count}). Skipping creation of new issue for {candidate}.")
            continue
        
        if not new_issue_created:
            # Create the issue
            group, kind = candidate.split()
            lower_kind = kind.lower()
            
            body = f"""As part of moving resources from terraform controllers to direct controllers (Epic #5954), we need to create the Go types for `{kind}`.

Currently, `{kind}` is managed by the Terraform controller (marked with `tf2crd=true`). The goal is to create the Go types in `apis/{group}/v1beta1/` so that we can eventually migrate the controller implementation to the "direct" approach.

### Instructions

1.  **Create a generate.sh**:
    Create `apis/{group}/v1beta1/generate.sh` which includes `{kind}`.
    It likely maps to something like `google.cloud.{group}.v1`.
    Example:
    ```bash
    go run . generate-types \\
      --service google.cloud.{group}.v1 \\
      --api-version {group}.cnrm.cloud.google.com/v1beta1 \\
      --resource {kind}:{kind} \\
      --include-skipped-output

    go run . generate-mapper \\
      --service google.cloud.{group}.v1 \\
      --api-version {group}.cnrm.cloud.google.com/v1beta1 \\
      --include-skipped-output
    ```

2.  Set the write permission on the new `apis/{group}/v1beta1/generate.sh` file. You should do this by running both `chmod +x apis/{group}/v1beta1/generate.sh` and `git add --chmod=+x apis/{group}/v1beta1/generate.sh`.

3.  **Generate Scaffolding**:
    Run `apis/{group}/v1beta1/generate.sh`. This should create `apis/{group}/v1beta1/{lower_kind}_types.go`.

4.  **Iterate on Types**:
    Compare the generated CRD with the existing one using `dev/tasks/diff-crds`.
    Modify `apis/{group}/v1beta1/{lower_kind}_types.go` until the CRD matches the existing one at `config/crds/resources/apiextensions.k8s.io_v1_customresourcedefinition_{lower_kind}s.{group}.cnrm.cloud.google.com.yaml`.

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
   - Running `dev/tasks/generate-types-and-mappers` will generate the mapper code once the `apis/{group}/v1beta1/{lower_kind}_types.go` file is generating an equivalent CRD.
   - Run `make all-binary` to ensure the generated mapper code compiles. Please fix any issue discovered by this compilation.

This issue is part of Epic #5954.
"""
            with open("issue_body.txt", "w") as fbody:
                fbody.write(body)
            
            title = f"Create generate.sh and types.go files for {candidate}"
            print(f"Creating new issue: {title}")
            create_cmd = f'gh issue create --title "{title}" --body-file issue_body.txt --label "overseer,area/direct,priority/medium"'
            out = run_cmd(create_cmd)
            print(f"Created issue: {out}")
            
            # The prompt says: "In a single run create at most one issue to avoid overwhelming the team."
            new_issue_created = True
            pending_count += 1
