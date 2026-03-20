import yaml
import os
import glob

files = glob.glob('config/crds/resources/*.yaml')
candidates = []

for f in files:
    with open(f, 'r') as file:
        try:
            doc = yaml.safe_load(file)
            labels = doc.get('metadata', {}).get('labels', {})
            if labels.get('cnrm.cloud.google.com/dcl2crd') == 'true':
                versions = doc.get('spec', {}).get('versions', [])
                is_beta = any(v.get('name') == 'v1beta1' for v in versions)
                group = doc.get('spec', {}).get('group')
                kind = doc.get('spec', {}).get('names', {}).get('kind')
                if is_beta and group and kind:
                    group_short = group.split('.')[0]
                    # Check if types.go exists
                    types_file_pattern = f"apis/{group_short}/v1beta1/*_types.go"
                    types_files = glob.glob(types_file_pattern)
                    if not types_files:
                        candidates.append((group_short, kind))
        except Exception as e:
            pass

print("\n".join(f"{g} {k}" for g, k in sorted(candidates)))
