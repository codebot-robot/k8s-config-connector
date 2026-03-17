import os
import glob
import yaml

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
                            print(f"{group} {kind} {short_group}")
    except Exception as e:
        pass