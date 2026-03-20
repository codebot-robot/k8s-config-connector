#!/bin/bash
for file in $(grep -rl 'cnrm.cloud.google.com/dcl2crd: "true"' config/crds/resources/ | sort); do
    group=$(yq '.spec.group' $file | awk -F. '{print $1}')
    kind=$(yq '.spec.names.kind' $file)
    version=$(yq '.spec.versions[] | select(.name == "v1beta1") | .name' $file)
    
    if [[ "$version" == "v1beta1" ]]; then
        # Check if types.go exists
        has_types=0
        if ls apis/$group/v1beta1/*_types.go 1> /dev/null 2>&1; then
            has_types=1
        fi
        
        if [[ $has_types -eq 0 ]]; then
            echo "$group $kind"
            exit 0
        fi
    fi
done
