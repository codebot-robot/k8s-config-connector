*   Special shout-outs to acpana, anhdle-sso, app/dependabot, barney-s, cheftako, codebot-robot, dhavalbhensdadiya-crest, fkc1e100, gemmahou, justinsb, katrielt, maqiuyujoyce, and xiaoweim for their contributions to this release.

## New Fields

*   `CloudDeployTarget`
    *   Updated `TargetIDs` field to `Ref`.

## New Features

*   Added structured reporting diff to MetastoreService, BigQueryTable, BigQueryReservationAssignment, WorkflowsWorkflow, PrivilegedAccessManagerEntitlement, DataformRepository, CloudIdentityGroup, BigQueryDataset, and CertificateManagerDNSAuthorization.
*   Added a KCC release agent.
*   Added CRD filtering for preview recorder.
*   Renamed DeployCustomTargetType to CloudDeployCustomTargetType.
*   Added Go types, Fuzzer, Mapper, Identity, and Reference files for CloudBuild Trigger.
*   Created generate.sh and types.go files for DataCatalog PolicyTag.
*   Implemented MockGCP for NetworkServices LBRouteExtension.

## Bug Fixes

*   [Issue 6878](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6878)
    Fix typo in StorageDefaultObjectAccessControl documentation.
*   [Issue 6983](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6983)
    Fix use of 'latest' in generate.sh.
*   [Issue 6844](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6844)
    Fix port contention.
