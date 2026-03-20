import os
import glob
import re
import json
import subprocess
import sys

def run_cmd(cmd):
    return subprocess.check_output(cmd, shell=True).decode('utf-8')

def find_candidates():
    crd_dir = 'config/crds/resources'
    crd_files = glob.glob(os.path.join(crd_dir, '*.yaml'))
    candidates = []

    for f in crd_files:
        with open(f, 'r') as fp:
            content = fp.read()
        
        # Check for dcl2crd
        if 'cnrm.cloud.google.com/dcl2crd: "true"' not in content:
            continue
            
        group_match = re.search(r'\bgroup:\s*([^\n]+)', content)
        
        # Kind is in spec.names.kind
        kind_match = re.search(r'names:\n(?:.*\n)*?\s*kind:\s*([^\n]+)', content)
        if not kind_match:
            kind_match = re.search(r'kind:\s*([^\n]+)', content[content.find('names:'):])
        
        if not group_match or not kind_match:
            continue
        
        group = group_match.group(1).strip().split('.')[0]
        kind = kind_match.group(1).strip()
        
        if 'name: v1beta1' not in content and '- name: v1beta1' not in content:
            continue
            
        type_files = glob.glob(f'apis/{group}/v1beta1/*_types.go')
        if type_files:
            continue
            
        candidates.append((group, kind))
    return candidates

def main():
    candidates = find_candidates()
    
    # Get all issues with the matching title prefix
    issues_json = run_cmd('gh issue list --search "in:title Create generate.sh and types.go files for" --state all --json number,title,state,labels -L 100')
    issues = json.loads(issues_json)

    pending_count = sum(1 for issue in issues if issue['state'] == 'OPEN')

    for group, kind in candidates:
        target_title_lower = f"Create generate.sh and types.go files for {group} {kind}".lower()
        
        existing_issue = None
        for issue in issues:
            if issue['title'].lower() == target_title_lower:
                existing_issue = issue
                break
                
        if existing_issue:
            # Check labels
            current_labels = [l['name'] for l in existing_issue['labels']]
            required_labels = ["overseer", "area/direct", "priority/medium"]
            missing_labels = [l for l in required_labels if l not in current_labels]
            if missing_labels:
                print(f"Injecting missing labels {missing_labels} for issue #{existing_issue['number']}")
                run_cmd(f'gh issue edit {existing_issue["number"]} --add-label {",".join(missing_labels)}')
            continue
        else:
            if pending_count >= 10:
                print("There are already more than 10 pending issues for this task. Skipping creating new ones until some of the existing issues are resolved.")
                sys.exit(0)
            
            # We would create an issue here if pending_count < 10
            # Wait, the instruction says to create AT MOST ONE.
            
            issue_title = f"Create generate.sh and types.go files for {group} {kind}"
            issue_body = f"""As part of moving resources from terraform controllers to direct controllers (Epic #5954), we need to create the Go types for `{kind}`.

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
    Run `apis/{group}/v1beta1/generate.sh`. This should create `apis/{group}/v1beta1/{kind.lower()}_types.go`.

4.  **Iterate on Types**:
    Compare the generated CRD with the existing one using `dev/tasks/diff-crds`.
    Modify `apis/{group}/v1beta1/{kind.lower()}_types.go` until the CRD matches the existing one at `config/crds/resources/apiextensions.k8s.io_v1_customresourcedefinition_{kind.lower()}s.{group}.cnrm.cloud.google.com.yaml`.

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
            
            with open("issue_body.txt", "w") as f:
                f.write(issue_body)
                
            cmd = f'gh issue create --title "{issue_title}" --body-file issue_body.txt --label "overseer,area/direct,priority/medium"'
            # Ensure it is a subtask by adding it to the project or mentioning epic
            # We already have "This issue is part of Epic #5954." in the body, but let's also pass milestone/project if needed.
            print(f"Creating issue: {issue_title}")
            run_cmd(cmd)
            sys.exit(0)

if __name__ == "__main__":
    main()
