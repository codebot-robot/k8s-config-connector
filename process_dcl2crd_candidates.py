import json
import os
import subprocess
import sys

# Load candidates
candidates = []
with open("find_candidates_dcl2crd_v2.py") as f:
    # Just run it and get output
    out = subprocess.check_output(["python3", "find_candidates_dcl2crd_v2.py"]).decode("utf-8")
    for line in out.strip().split("\n"):
        if line:
            candidates.append(line.split(" "))

# Load issues
with open("all_dcl2crd_issues.json") as f:
    issues = json.load(f)

pending_count = sum(1 for i in issues if i["state"] == "OPEN")
required_labels = {"overseer", "area/direct", "priority/medium"}

# Map issue titles to issues
issue_map = {}
for i in issues:
    issue_map[i["title"].lower()] = i

created_one = False
print(f"Pending issues: {pending_count}")

for group, kind in candidates:
    title = f"Create generate.sh and types.go files for {group} {kind}"
    
    if title.lower() in issue_map:
        issue = issue_map[title.lower()]
        existing_labels = {l["name"] for l in issue["labels"]}
        missing_labels = required_labels - existing_labels
        
        if missing_labels:
            labels_to_add = ",".join(missing_labels)
            print(f"Adding labels {labels_to_add} to issue #{issue['number']} ({title})")
            subprocess.check_call(["gh", "issue", "edit", str(issue["number"]), "--add-label", labels_to_add])
        else:
            # print(f"Issue #{issue['number']} already has all labels.")
            pass
    else:
        # Issue doesn't exist
        if pending_count >= 10:
            print(f"Skipping {group} {kind}: There are already >= 10 pending issues.")
            continue
        
        if not created_one:
            # Create issue
            print(f"Creating issue for {group} {kind}")
            
            body = f"""As part of moving resources from terraform controllers to direct controllers (Epic #5954), we need to create the Go types for `{kind}`.

Currently, `{kind}` is managed by the Terraform controller (marked with `tf2crd=true`). The goal is to create the Go types in `apis/{group.lower()}/v1beta1/` so that we can eventually migrate the controller implementation to the "direct" approach.

### Instructions

1.  **Create a generate.sh**:
    Create `apis/{group.lower()}/v1beta1/generate.sh` which includes `{kind}`.
    It likely maps to something like `google.cloud.{group.lower()}.v1`.
    Example:
    ```bash
    go run . generate-types \\
      --service google.cloud.{group.lower()}.v1 \\
      --api-version {group.lower()}.cnrm.cloud.google.com/v1beta1 \\
      --resource {kind}:{kind} \\
      --include-skipped-output

    go run . generate-mapper \\
      --service google.cloud.{group.lower()}.v1 \\
      --api-version {group.lower()}.cnrm.cloud.google.com/v1beta1 \\
      --include-skipped-output
    ```

2.  Set the write permission on the new `apis/{group.lower()}/v1beta1/generate.sh` file. You should do this by running both `chmod +x apis/{group.lower()}/v1beta1/generate.sh` and `git add --chmod=+x apis/{group.lower()}/v1beta1/generate.sh`.

3.  **Generate Scaffolding**:
    Run `apis/{group.lower()}/v1beta1/generate.sh`. This should create `apis/{group.lower()}/v1beta1/{kind.lower()}_types.go`.

4.  **Iterate on Types**:
    Compare the generated CRD with the existing one using `dev/tasks/diff-crds`.
    Modify `apis/{group.lower()}/v1beta1/{kind.lower()}_types.go` until the CRD matches the existing one at `config/crds/resources/apiextensions.k8s.io_v1_customresourcedefinition_{group.lower()}{kind.lower()}s.{group.lower()}.cnrm.cloud.google.com.yaml`.

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
   - Running `dev/tasks/generate-types-and-mappers` will generate the mapper code once the `apis/{group.lower()}/v1beta1/{kind.lower()}_types.go` file is generating an equivalent CRD.
   - Run `make all-binary` to ensure the generated mapper code compiles. Please fix any issue discovered by this compilation.

This issue is part of Epic #5954."""
            with open("issue_body.md", "w") as f_body:
                f_body.write(body)

            labels_str = ",".join(required_labels)
            cmd = [
                "gh", "issue", "create",
                "--title", title,
                "--body-file", "issue_body.md",
                "--label", labels_str
            ]
            print(f"Running cmd: {' '.join(cmd)}")
            subprocess.check_call(cmd)
            
            created_one = True
            pending_count += 1
