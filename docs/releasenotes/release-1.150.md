** This version is not yet released; this document is gathering release notes
for the future release **

*   Special shout-outs to anhdle-sso, app/dependabot, barney-s, cheftako, codebot-robot, gurusai-voleti, justinsb, katrielt, maqiuyujoyce, xiaoweim for their contributions to this release.

## New Beta Resources (Direct Reconciler):

*   [`AlloyDBUser`](https://cloud.google.com/config-connector/docs/reference/resource-docs/alloydb/alloydbuser)

    *   Manage AlloyDB users.

*   [`AccessContextManagerServicePerimeter`](https://cloud.google.com/config-connector/docs/reference/resource-docs/accesscontextmanager/accesscontextmanagerserviceperimeter)

    *   Manage Access Context Manager Service Perimeters.

## New features:

*   Added a `--skip-name-validation` flag to Config Connector controllers to bypass duplicate controller name checks during registration.

*   Added detailed documentation on how to enable VerticalPodAutoscaler (VPA) for Config Connector pods using ControllerResource and NamespacedControllerResource.

*   Support alternative controller comparison in preview mode.

*   Add structured reporting diff to BigQueryAnalyticsHubDataExchange.

## Bug Fixes:

*   [Issue 6943](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6943) Handle ALREADY_EXISTS error during TagKey and TagValue resource creation.
*   [Issue 7106](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7106) Fix typo in container.yaml service mapping file.
*   [Issue 7115](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7115) Restore missing descriptions in CloudBuildTrigger CRD.
