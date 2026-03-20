import os
import glob
import yaml

files = glob.glob("config/crds/resources/*.yaml")
candidates = []

for f in files:
    with open(f, 'r') as fh:
        try:
            content = fh.read()
            if 'cnrm.cloud.google.com/dcl2crd: "true"' in content:
                docs = yaml.safe_load_all(content)
                for doc in docs:
                    if doc and doc.get("kind") == "CustomResourceDefinition":
                        group = doc["spec"]["group"]
                        kind = doc["spec"]["names"]["kind"]
                        short_group = group.split(".")[0]
                        versions = [v["name"] for v in doc["spec"].get("versions", [])]
                        
                        if "v1beta1" in versions:
                            # check if types.go exists
                            types_path_pattern = f"apis/{short_group}/v1beta1/*_types.go"
                            if not glob.glob(types_path_pattern):
                                candidates.append((short_group, kind))
        except Exception as e:
            pass

for c in sorted(candidates):
    print(f"{c[0]},{c[1]}")
