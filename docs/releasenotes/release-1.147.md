** This version is not yet released; this document is gathering release notes for the future release **

* Special shout-outs to acpana, anhdle-sso, app/dependabot, barney-s, cheftako, codebot-robot, dhavalbhensdadiya-crest, fkc1e100, gemmahou, justinsb, katrielt, maqiuyujoyce, xiaoweim for their contributions to this release.

## New Alpha Resources (Direct Reconciler)
* `CloudBuildTrigger`
    * Added Go types for CloudBuildTrigger.
* `DataCatalogPolicyTag`
    * Added basic implementation for DataCatalogPolicyTag.

## New Fields
* `CloudDeployTarget`
    * Updated TargetIDs field to Ref.

## Reconciliation Improvements
* Added structured reporting diff to `WorkflowsWorkflow`, `PrivilegedAccessManagerEntitlement`, `DataformRepository`, `CloudIdentityGroup`, `BigQueryDataset`, `CertificateManagerDNSAuthorization`.

## New features
* `CloudDeployCustomTargetType`
    * Renamed DeployCustomTargetType to CloudDeployCustomTargetType.

## Bug Fixes
* [Issue 6878](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6878)
  Fix typo in StorageDefaultObjectAccessControl documentation.
* [Issue 6879](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6879)
  Add missing reference doc and samples.
