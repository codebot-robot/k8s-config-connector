*   Special shout-outs to acpana, anhdle-sso, app/dependabot, barney-s, cheftako, codebot-robot, dhavalbhensdadiya-crest, fkc1e100, gemmahou, justinsb, katrielt, maqiuyujoyce, and xiaoweim for their contributions to this release.

## New Fields

*   `ContainerCluster`
    *   Marked `bootDiskKMSKeyRef` as mutable.

## New features:

*   Structured reporting diff added for the following resources:
    *   `DataformRepository`
    *   `CloudIdentityGroup`
    *   `BigQueryDataset`
    *   `CertificateManagerDNSAuthorization`
    *   `VMwareEngineExternalAccessRule`
*   `Preview Recorder`
    *   Added CRD filtering support.

## Bug Fixes:

*   Fix stability label for v1beta1 resources.
*   Update `CloudDeployTarget` to rename `multiTarget.targets` to `targetRefs`.
*   Fix typo in `StorageDefaultObjectAccessControl` documentation.
