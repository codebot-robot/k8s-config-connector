import yaml
import os
import glob

files = glob.glob('config/crds/resources/*.yaml')
for f in files:
    try:
        with open(f, 'r') as fp:
            data = yaml.safe_load(fp)
            labels = data.get('metadata', {}).get('labels', {})
            if labels.get('cnrm.cloud.google.com/dcl2crd') == 'true':
                versions = data.get('spec', {}).get('versions', [])
                has_beta = any(v.get('name') == 'v1beta1' for v in versions)
                if has_beta:
                    group = data.get('spec', {}).get('group', '').split('.')[0]
                    kind = data.get('spec', {}).get('names', {}).get('kind', '')
                    
                    types_glob = f'apis/{group}/v1beta1/*_types.go'
                    if not glob.glob(types_glob):
                        print(f"{group} {kind}")
    except Exception as e:
        pass
