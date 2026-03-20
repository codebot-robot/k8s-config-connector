import os
import re
import sys

crd_dir = "config/crds/resources"
files = [f for f in os.listdir(crd_dir) if f.endswith(".yaml")]

candidates = []

for f in files:
    filepath = os.path.join(crd_dir, f)
    with open(filepath) as fp:
        content = fp.read()
        if 'cnrm.cloud.google.com/dcl2crd: "true"' in content:
            # Extract group
            group_match = re.search(r'\n  group: ([^\s]+)', content)
            group = group_match.group(1) if group_match else None
            
            # Extract kind
            kind_match = re.search(r'\n    kind: ([^\s]+)', content)
            kind = kind_match.group(1) if kind_match else None
            
            # Check for beta
            if 'name: v1beta1' in content and group and kind:
                short_group = group.split('.')[0]
                
                # Check if types.go exists
                # e.g., apis/<GROUP>/v1beta1/*_types.go
                api_dir = os.path.join("apis", short_group, "v1beta1")
                has_types = False
                if os.path.exists(api_dir):
                    for api_file in os.listdir(api_dir):
                        if api_file.endswith("_types.go"):
                            # To be strict, check if the file name matches the kind?
                            # Usually any types.go might mean it's generated, but let's check if there is ANY types.go that contains the Kind
                            with open(os.path.join(api_dir, api_file)) as f_go:
                                go_content = f_go.read()
                                if f"type {kind} struct" in go_content or f"type {kind}List struct" in go_content:
                                    has_types = True
                                    break
                
                if not has_types:
                    candidates.append(f"{short_group} {kind}")

for c in sorted(candidates):
    print(c)
