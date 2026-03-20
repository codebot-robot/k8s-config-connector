import yaml
import glob
import os

candidates = []
for file_path in glob.glob('config/crds/resources/*.yaml'):
    with open(file_path, 'r') as f:
        try:
            doc = yaml.safe_load(f)
        except Exception:
            continue
        
        labels = doc.get('metadata', {}).get('labels', {})
        if labels.get('cnrm.cloud.google.com/dcl2crd') == 'true':
            versions = doc.get('spec', {}).get('versions', [])
            for v in versions:
                version = v.get('name')
                if 'beta' in version:
                    group = doc.get('spec', {}).get('group', '').split('.')[0]
                    kind = doc.get('spec', {}).get('names', {}).get('kind', '')
                    
                    types_pattern = f"apis/{group}/{version}/*_types.go"
                    if not glob.glob(types_pattern):
                        print(f"{group} {kind} {version}")
