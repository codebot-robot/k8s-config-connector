import json
import subprocess
import sys
import re

# 1. Get candidates
with open("candidates.txt", "w") as f:
    subprocess.run(["python3", "find_candidates_dcl2crd.py"], stdout=f)

with open("candidates.txt", "r") as f:
    candidates = [line.strip() for line in f if line.strip()]

# 2. Get all issues
result = subprocess.run(["gh", "issue", "list", "--state", "all", "--limit", "500", "--json", "number,title,labels,state"], capture_output=True, text=True)
issues = json.loads(result.stdout)

# Find pending issues for this task
pending_issues = [i for i in issues if "Create generate.sh and types.go files for" in i["title"] and i["state"] == "OPEN"]
pending_count = len(pending_issues)

# 3. Process candidates
for cand in candidates:
    group, kind = cand.split()
    
    # Check if issue exists
    # We look for an issue whose title ends with "Group Kind" (case insensitive match for group/kind)
    matching_issue = None
    for i in issues:
        title = i["title"].lower()
        if f"for {group.lower()} {kind.lower()}" in title:
            matching_issue = i
            break
            
    if matching_issue:
        # Check labels
        existing_labels = [l["name"] for l in matching_issue["labels"]]
        missing_labels = []
        for required in ["overseer", "area/direct", "priority/medium"]:
            if required not in existing_labels:
                missing_labels.append(required)
                
        if missing_labels:
            print(f"Injecting labels {missing_labels} for issue #{matching_issue['number']} ({matching_issue['title']})")
            subprocess.run(["gh", "issue", "edit", str(matching_issue['number']), "--add-label", ",".join(missing_labels)])
        continue
    else:
        # Need to create
        if pending_count >= 10:
            print(f"There are already {pending_count} pending issues for this task. Skipping creation for {group} {kind} to avoid overwhelming the team.")
            continue
        else:
            print(f"Creating issue for {group} {kind}")
            # we would create the issue here, but we're past the limit
            pending_count += 1

