** This version is not yet released; this document is gathering release notes
for the future release **

*   Special shout-outs to acpana, anhdle-sso, app/dependabot, barney-s, cheftako, codebot-robot, dhavalbhensdadiya-crest, fkc1e100, gemmahou, justinsb, katrielt, maqiuyujoyce, xiaoweim for their contributions to this release.

## Announcement

### Simplified and More Reliable Resource Development

*   We launched a major improvement to the Config Connector resource
    development! Our new approach significantly enhances reliability and
    provides a more native Kubernetes experience. Learn more in our
    [guide](https://github.com/GoogleCloudPlatform/k8s-config-connector/tree/master/docs/develop-resources)

## New Alpha Resources (Direct Reconciler):

*   `CloudBuildTrigger`
    *   Added Go types and Identity/Reference files for CloudBuildTrigger.
*   `DataCatalogPolicyTag`
    *   Added Go types and generate script.

## New Fields:

*   [`CloudDeployTarget`](https://cloud.google.com/config-connector/docs/reference/resource-docs/clouddeploy/clouddeploytarget)
    *   Update CloudDeployTarget `TargetIDs` field to `Ref`.

## Reconciliation Improvements

We have added support for structured reporting diff to the controller in the following resources:
*   `BigQueryDataset`
*   `BigQueryReservationAssignment`
*   `CertificateManagerDNSAuthorization`
*   `CloudIdentityGroup`
*   `DataformRepository`
*   `PrivilegedAccessManagerEntitlement`
*   `WorkflowsWorkflow`

## New features:

*   Clarify Stability of Resources Served as v1beta1 but Labeled Alpha
    *   Promoted to Stable: `TagsTagKey`, `TagsTagValue`, `TagsTagBinding`, `TagsLocationTagBinding`, `BigQueryRoutine`, `BigQueryAnalyticsHubDataExchange`, `ConfigControllerInstance`, `DataCatalogTaxonomy`, `DataformRepository`, `VertexAIMetadataStore`.
*   Add CRD filtering for preview recorder
    *   Added CRD filtering for the preview recorder to skip non-CNRM objects and resources listed in IgnoredCRDList.
*   Rename `DeployCustomTargetType` to `CloudDeployCustomTargetType`
    *   Renamed the resource to maintain naming consistency with other Cloud Deploy resources.
*   Implement MockGCP for NetworkServices `LBRouteExtension`

## Bug Fixes:

*   [Issue 6878](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6878)
    Fix typo in `StorageDefaultObjectAccessControl` documentation.
*   [Issue 6983](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6983)
    Fix use of 'latest' in `generate.sh`.