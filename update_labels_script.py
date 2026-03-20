import json
import subprocess
import glob
import yaml

# get candidates
candidates = set()
for file in glob.glob("config/crds/resources/*.yaml"):
    with open(file, "r") as f:
        content = f.read()
        if 'cnrm.cloud.google.com/dcl2crd: "true"' in content:
            docs = yaml.safe_load_all(content)
            for doc in docs:
                if not doc: continue
                group = doc["spec"]["group"].split(".")[0]
                kind = doc["spec"]["names"]["kind"]
                versions = [v["name"] for v in doc["spec"]["versions"]]
                if "v1beta1" in versions:
                    types_files = glob.glob(f"apis/{group}/v1beta1/*_types.go")
                    if not types_files:
                        candidates.add(f"{group} {kind}")

with open("all_issues_for_task.json", "r") as f:
    issues = json.load(f)

required_labels = {"overseer", "area/direct", "priority/medium"}

for issue in issues:
    title = issue["title"]
    labels = {l["name"] for l in issue.get("labels", [])}
    
    # Check if title matches any candidate
    for cand in candidates:
        if cand.lower() in title.lower():
            missing_labels = required_labels - labels
            if missing_labels:
                issue_number = issue["number"]
                print(f"Issue #{issue_number} for {cand} missing labels: {missing_labels}")
                subprocess.run(["gh", "issue", "edit", str(issue_number), "--add-label", ",".join(missing_labels)])
            break
