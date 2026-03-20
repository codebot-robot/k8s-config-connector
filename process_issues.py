import json
import subprocess
import os

def run(cmd):
    return subprocess.run(cmd, shell=True, capture_output=True, text=True)

with open('all_migration_issues.json') as f:
    issues = json.load(f)

# pending issues
open_issues = [i for i in issues if i['state'] == 'OPEN' and 'Create generate.sh and types.go files for' in i['title']]
pending_count = len(open_issues)

# Run get_candidates.py to get candidates
out = run('python3 get_candidates.py')
candidates = []
for line in out.stdout.splitlines():
    if line.strip():
        parts = line.strip().split()
        if len(parts) == 2:
            candidates.append((parts[0], parts[1]))

issue_created = False

for group, kind in candidates:
    title_expected = f"Create generate.sh and types.go files for {group} {kind}"
    # find if issue exists
    existing = [i for i in issues if title_expected.lower() == i['title'].lower()]
    
    if existing:
        # Check labels
        issue = existing[0]
        issue_number = issue['number']
        labels = [l['name'] for l in issue['labels']]
        missing_labels = []
        for l in ['overseer', 'area/direct', 'priority/medium']:
            if l not in labels:
                missing_labels.append(l)
        
        if missing_labels:
            print(f"Injecting labels {missing_labels} for issue #{issue_number} ({group} {kind})")
            run(f"gh issue edit {issue_number} --add-label " + ",".join(missing_labels))
            
    else:
        # Create issue
        if pending_count > 10:
            print(f"There are already {pending_count} pending issues (> 10). Skipping creation for new ones.")
            # According to instructions, we should skip creating new ones, but we still want to inject labels for existing ones, so we just continue
            # Wait, the prompt says "skip creating new ones" so we don't break the loop, just skip creation.
            continue
        
        if not issue_created:
            print(f"Creating issue for {group} {kind}")
            
            crd_file_name_guess = f"apiextensions.k8s.io_v1_customresourcedefinition_{kind.lower()}s.{group}.cnrm.cloud.google.com.yaml"
            # It's better to just use a placeholder or let the developer find it
            
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
      --resource {kind}:___TBD___ \\
      --include-skipped-output

    go run . generate-mapper \\
      --service google.cloud.{group}.v1 \\
      --api-version {group}.cnrm.cloud.google.com/v1beta1 \\
      --include-skipped-output
    ```

2.  Set the write permission on the new `apis/{group}/v1beta1/generate.sh` file. You should do this by running both `chmod +x apis/{group}/v1beta1/generate.sh` and `git add --chmod=+x apis/{group}/v1beta1/generate.sh`.

3.  **Generate Scaffolding**:
    Run `apis/{group}/v1beta1/generate.sh`. This should create `apis/{group}/v1beta1/{kind.lower()}_types.go`.

4.  **Iterate on Types**:
    Compare the generated CRD with the existing one using `dev/tasks/diff-crds`.
    Modify `apis/{group}/v1beta1/{kind.lower()}_types.go` until the CRD matches the existing one at `config/crds/resources/{crd_file_name_guess}`.

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
   - Running `dev/tasks/generate-types-and-mappers` will generate the mapper code once the `apis/{group}/v1beta1/{kind.lower()}_types.go` file is generating an equivalent CRD.
   - Run `make all-binary` to ensure the generated mapper code compiles. Please fix any issue discovered by this compilation.

This issue is part of Epic #5954.
"""
            with open("body.txt", "w") as f_body:
                f_body.write(body)
                
            cmd = f'gh issue create --title "{title_expected}" --body-file body.txt --label overseer,area/direct,priority/medium'
            run_cmd = run(cmd)
            print("Issue creation output:", run_cmd.stdout)
            print("Issue creation errors:", run_cmd.stderr)
            
            issue_created = True
            pending_count += 1
