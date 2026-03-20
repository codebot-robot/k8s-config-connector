import json
import subprocess
import sys
import re

# Load all issues
with open("all_issues.json") as f:
    issues = json.load(f)

# Parse candidates
candidates = []
with open("candidates.txt") as f:
    for line in f:
        line = line.strip()
        if line:
            parts = line.split(" ")
            if len(parts) == 2:
                candidates.append((parts[0], parts[1]))

# Create a map of issues
issue_map = {}
open_count = 0
for issue in issues:
    title = issue["title"]
    labels = [lbl["name"] for lbl in issue.get("labels", [])]
    number = issue["number"]
    
    # Check if open
    # Note: the all_issues.json might not have state, wait, I didn't query state in the JSON fields.
    # Actually I used `gh issue list --search "Create generate.sh and types.go files for in:title" --state all --json title,number,labels,state`? No, I forgot `state`.
    
    match = re.search(r'Create generate.sh and types.go files for (\w+) (\w+)', title, re.IGNORECASE)
    if match:
        group = match.group(1).lower()
        kind = match.group(2).lower()
        issue_map[(group, kind)] = {
            "number": number,
            "labels": labels,
            "title": title
        }

print(f"Found {len(candidates)} candidates.")

for group, kind in candidates:
    key = (group.lower(), kind.lower())
    if key in issue_map:
        issue = issue_map[key]
        number = issue["number"]
        labels = issue["labels"]
        missing_labels = []
        for required_label in ["overseer", "area/direct", "priority/medium"]:
            if required_label not in labels:
                missing_labels.append(required_label)
        
        if missing_labels:
            print(f"Adding labels {missing_labels} to issue #{number} for {group} {kind}")
            subprocess.run(["gh", "issue", "edit", str(number), "--add-label", ",".join(missing_labels)], check=True)
        else:
            print(f"Issue #{number} for {group} {kind} already has all labels.")
    else:
        print(f"No issue found for {group} {kind}. Needs creation.")

