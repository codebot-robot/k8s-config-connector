** This version is not yet released; this document is gathering release notes
for the future release **

*   Special shout-outs to anhdle-sso, app/dependabot, barney-s, cheftako, codebot-robot, gurusai-voleti, justinsb, katrielt, maqiuyujoyce, and xiaoweim for their contributions to this release.

## New Alpha Resources (Direct Reconciler):

*   `AccessContextManagerServicePerimeter`
*   `AlloyDBUser`
*   `ComputeResourcePolicy`

## Reconciliation Improvements

*   [`BigQueryAnalyticsHubDataExchange`](https://cloud.google.com/config-connector/docs/reference/resource-docs/bigqueryanalyticshub/bigqueryanalyticshubdataexchange)
    *   Add structured reporting diff.

## New features:

*   Allow specific integer type changes in CRD equivalence checks.
*   Add documentation for enabling VerticalPodAutoscaler in Config Connector.
*   Support alternative controller comparison in preview mode.
*   Introduce `skip-name-validation` flag and consolidate tests.

## Bug Fixes:

*   [Issue 6943](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6943) Handle `ALREADY_EXISTS` in `TagKey` and `TagValue` controllers.
*   [Issue 7115](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7115) Restore missing descriptions in `CloudBuildTrigger` CRD.
*   [Issue 7106](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7106) Fix typo in container.yaml service mapping file.
*   [Issue 7082](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7082) Fix `TestCloudBuildTriggerFuzz` failure.
