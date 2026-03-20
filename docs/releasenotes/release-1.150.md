*   Special shout-outs to @anhdle-sso, @app/dependabot, @barney-s, @cheftako, @codebot-robot, @gurusai-voleti, @justinsb, @katrielt, @maqiuyujoyce, and @xiaoweim for their contributions to this release.

## New Features

*   Added documentation for enabling VerticalPodAutoscaler in Config Connector.

## Reconciliation Improvements

We have added support for direct reconciliation to more resources, with opt-in behaviour. The API is unchanged. To use the direct reconciler, add the `alpha.cnrm.cloud.google.com/reconciler: direct` annotation to the corresponding Config Connector object. The following resources now have direct reconciliation support:

*   `ComputeResourcePolicy`
*   `AlloyDBUser`
*   `AccessContextManagerServicePerimeter`

## Bug Fixes

*   [Issue 7118](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7118) Improve Release Pipeline Robustness and GitHub Workflow.
*   [Issue 7115](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7115) Restore missing descriptions in CloudBuildTrigger CRD.
*   [Issue 7106](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7106) Fix typo in container.yaml service mapping file.
*   [Issue 6943](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6943) Fix handle ALREADY_EXISTS in TagKey and TagValue controllers.
*   [Issue 7082](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7082) Fix TestCloudBuildTriggerFuzz failure.
*   [Issue 6881](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6881) Align mock responses with real GCP behavior for CertificateIssuanceConfig.
*   [Issue 6774](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6774) Add structured reporting diff to BigQueryAnalyticsHubDataExchange.
*   [Issue 7012](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7012) Allow specific integer type changes in CRD equivalence checks.