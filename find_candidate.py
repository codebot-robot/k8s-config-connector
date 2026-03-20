import yaml
import glob
import os

with open('candidates.txt', 'w') as out:
    for filepath in glob.glob('config/crds/resources/*.yaml'):
        with open(filepath, 'r') as f:
            try:
                data = yaml.safe_load(f)
            except:
                continue
            
            labels = data.get('metadata', {}).get('labels', {})
            if labels.get('cnrm.cloud.google.com/dcl2crd') != 'true':
                continue
            
            versions = data.get('spec', {}).get('versions', [])
            has_beta = any(v.get('name') == 'v1beta1' for v in versions)
            if not has_beta:
                continue
                
            group = data.get('spec', {}).get('group', '').split('.')[0]
            kind = data.get('spec', {}).get('names', {}).get('kind', '')
            
            # Check if apis/<group>/v1beta1/*_types.go exists
            types_files = glob.glob(f'apis/{group}/v1beta1/*_types.go')
            if not types_files:
                out.write(f'{group} {kind}\n')

