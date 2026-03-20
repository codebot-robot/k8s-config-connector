import json
import os
import subprocess
import yaml

with open("all_issues.json") as f:
    issues = json.load(f)

issue_map = {}
for issue in issues:
    title = issue["title"]
    if "Create generate.sh and types.go files for " in title:
        parts = title.replace("Create generate.sh and types.go files for ", "").strip().split(" ")
        if len(parts) >= 2:
            group = parts[0]
            kind = parts[1]
            issue_map[f"{group} {kind}"] = issue

crd_dir = "config/crds/resources"
candidates = []

for filename in os.listdir(crd_dir):
    if not filename.endswith(".yaml"): continue
    filepath = os.path.join(crd_dir, filename)
    with open(filepath, "r") as f:
        try:
            doc = yaml.safe_load(f)
            labels = doc.get("metadata", {}).get("labels", {})
            if str(labels.get("cnrm.cloud.google.com/dcl2crd")).lower() == "true":
                group = doc["spec"]["group"].split(".")[0]
                kind = doc["spec"]["names"]["kind"]
                versions = doc.get("spec", {}).get("versions", [])
                for v in versions:
                    if "beta" in v["name"]:
                        version = v["name"]
                        api_dir = f"apis/{group}/{version}"
                        types_exist = False
                        if os.path.exists(api_dir):
                            for f_name in os.listdir(api_dir):
                                if f_name.endswith("_types.go"):
                                    types_exist = True
                                    break
                        if not types_exist:
                            candidates.append((group, kind, version))
                        break
        except Exception as e:
            pass

required_labels = ["overseer", "area/direct", "priority/medium"]
new_issues_needed = []

for group, kind, version in candidates:
    key = f"{group} {kind}"
    if key in issue_map:
        issue = issue_map[key]
        labels = [l["name"] for l in issue["labels"]]
        needs_update = False
        for rl in required_labels:
            if rl not in labels:
                needs_update = True
        if needs_update:
            num = issue["number"]
            print(f"Injecting labels for existing issue #{num} for {group} {kind}")
            subprocess.run(["gh", "issue", "edit", str(num), "--add-label", ",".join(required_labels)])
    else:
        new_issues_needed.append((group, kind, version))

print(f"Total new issues that would be created: {len(new_issues_needed)}")
with open("new_issues_needed.json", "w") as f:
    json.dump(new_issues_needed, f)
