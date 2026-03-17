import json
import sys
import subprocess

with open('existing_issues.json') as f:
    existing_issues = json.load(f)

# Collect existing titles
existing_titles = {issue['title']: issue for issue in existing_issues}

with open('find_candidates.py', 'r') as f:
    pass # we already ran it and have the output, but let's just parse the output from previous step. Wait, it's easier to just run it inside python

import glob
import yaml

candidates = []
for fpath in glob.glob("config/crds/resources/*.yaml"):
    try:
        with open(fpath, 'r') as f:
            docs = yaml.safe_load_all(f)
            for doc in docs:
                if not doc: continue
                labels = doc.get("metadata", {}).get("labels", {})
                if labels.get("cnrm.cloud.google.com/dcl2crd") == "true":
                    group = doc.get("spec", {}).get("group", "")
                    kind = doc.get("spec", {}).get("names", {}).get("kind", "")
                    
                    versions = doc.get("spec", {}).get("versions", [])
                    has_beta = any(v.get("name") == "v1beta1" for v in versions)
                    
                    if has_beta:
                        short_group = group.split(".")[0]
                        types_files = glob.glob(f"apis/{short_group}/v1beta1/*_types.go")
                        if not types_files:
                            candidates.append((short_group, kind))
    except Exception as e:
        pass

open_issue_count = sum(1 for issue in existing_issues)

for short_group, kind in candidates:
    expected_title = f"Create generate.sh and types.go files for {short_group} {kind}"
    if expected_title in existing_titles:
        # Check if labels are correct: overseer, area/direct, priority/medium
        issue = existing_titles[expected_title]
        existing_labels = {lbl['name'] for lbl in issue['labels']}
        required_labels = {"overseer", "area/direct", "priority/medium"}
        missing_labels = required_labels - existing_labels
        if missing_labels:
            print(f"Adding labels {missing_labels} to issue {issue['number']}")
            lbl_args = ",".join(missing_labels)
            subprocess.run(["gh", "issue", "edit", str(issue['number']), "--add-label", lbl_args])
            print(f"Added labels to {issue['number']}")
            sys.exit(0)
    else:
        # Needs to create an issue
        if open_issue_count >= 10:
            print(f"There are already 10 pending issues ({open_issue_count} total). Skipping creation for {short_group} {kind}.")
            sys.exit(0)
        else:
            print(f"Should create issue for {short_group} {kind}")
            # wait, if open_issue_count < 10, I can create one. But it's >= 10.
            sys.exit(0)
