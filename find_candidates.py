import os
import glob
import yaml

crd_dir = "config/crds/resources"
for file_path in glob.glob(os.path.join(crd_dir, "*.yaml")):
    try:
        with open(file_path, 'r') as f:
            data = yaml.safe_load(f)
            
        labels = data.get('metadata', {}).get('labels', {})
        if labels.get('cnrm.cloud.google.com/dcl2crd') != 'true':
            continue
            
        spec = data.get('spec', {})
        versions = spec.get('versions', [])
        has_beta = any(v.get('name') == 'v1beta1' for v in versions)
        if not has_beta:
            continue
            
        group = spec.get('group', '').split('.')[0]
        kind = spec.get('names', {}).get('kind', '')
        
        # Check if types.go exists
        types_pattern = os.path.join('apis', group, 'v1beta1', '*_types.go')
        if not glob.glob(types_pattern):
            print(f"Candidate: {group} {kind}")
    except Exception as e:
        pass
