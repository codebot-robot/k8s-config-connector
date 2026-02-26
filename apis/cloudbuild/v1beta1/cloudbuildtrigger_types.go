// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

import (
	pubsubv1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/pubsub/v1beta1"
	refsv1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	storagev1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/storage/v1beta1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/k8s/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var CloudBuildTriggerGVK = GroupVersion.WithKind("CloudBuildTrigger")

// CloudBuildTriggerSpec defines the desired state of CloudBuildTrigger
// +kcc:spec:proto=google.devtools.cloudbuild.v1.BuildTrigger
type CloudBuildTriggerSpec struct {
	// The location of this resource.
	Location string `json:"location,omitempty"`

	// Human-readable description of the trigger.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.description
	Description *string `json:"description,omitempty"`

	// Tags for annotation of a `BuildTrigger`
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.tags
	Tags []string `json:"tags,omitempty"`

	// Template describing the types of source changes to trigger a build.
	//
	//  Branch and tag names in trigger templates are interpreted as regular
	//  expressions. Any branch or tag change that matches that regular
	//  expression will trigger a build.
	//
	//  Mutually exclusive with `github`.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.trigger_template
	TriggerTemplate *TriggerTemplate `json:"triggerTemplate,omitempty"`

	// GitHubEventsConfig describes the configuration of a trigger that creates
	//  a build whenever a GitHub event is received.
	//
	//  Mutually exclusive with `trigger_template`.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.github
	Github *GitHubEventsConfig `json:"github,omitempty"`

	// PubsubConfig describes the configuration of a trigger that
	//  creates a build whenever a Pub/Sub message is published.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.pubsub_config
	PubsubConfig *TriggerPubsubConfig `json:"pubsubConfig,omitempty"`

	// WebhookConfig describes the configuration of a trigger that
	//  creates a build whenever a webhook is sent to a trigger's webhook URL.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.webhook_config
	WebhookConfig *TriggerWebhookConfig `json:"webhookConfig,omitempty"`

	// Contents of the build template.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.build
	Build *TriggerBuild `json:"build,omitempty"`

	// Path, from the source root, to a file whose contents is used for the
	//  template.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.filename
	Filename *string `json:"filename,omitempty"`

	// Substitutions for Build resource.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.substitutions
	Substitutions map[string]string `json:"substitutions,omitempty"`

	// ignored_files and included_files are file glob matches using
	//  https://golang.org/pkg/path/filepath/#Match extended with support for "**".
	//
	//  If ignored_files and changed files are both empty, then they are not
	//  used to determine whether or not to trigger a build.
	//
	//  If ignored_files is not empty, then we ignore any files that match any
	//  of the ignored_file globs. If the change has no files that are outside
	//  of the ignored_files globs, then we do not trigger a build.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.ignored_files
	IgnoredFiles []string `json:"ignoredFiles,omitempty"`

	// If any of the files altered in the commit pass the ignored_files filter
	//  and included_files is empty, then as far as this filter is concerned, we
	//  should trigger the build.
	//
	//  If any of the files altered in the commit pass the ignored_files filter
	//  and included_files is not empty, then we make sure that at least one of
	//  those files matches a included_files glob. If not, then we do not trigger
	//  a build.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.included_files
	IncludedFiles []string `json:"includedFiles,omitempty"`

	// A Common Expression Language string.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.filter
	Filter *string `json:"filter,omitempty"`

	// IAM service account whose credentials will be used at build runtime.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.service_account
	ServiceAccountRef *refsv1beta1.IAMServiceAccountRef `json:"serviceAccountRef,omitempty"`

	// The file source describing the local or remote Build template.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.git_file_source
	GitFileSource *TriggerGitFileSource `json:"gitFileSource,omitempty"`

	// The repo and ref of the repository from which to build.
	//  This field is used only for those triggers that do not respond to SCM events.
	//  Triggers that respond to such events build source at whatever commit caused
	//  the event.
	//  This field is currently only used by Webhook, Pub/Sub, Manual, and Cron
	//  triggers.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.source_to_build
	SourceToBuild *TriggerGitRepoSource `json:"sourceToBuild,omitempty"`

	// The configuration of a trigger that creates a build whenever an event from
	//  Repo API is received.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.repository_event_config
	RepositoryEventConfig *RepositoryEventConfig `json:"repositoryEventConfig,omitempty"`

	// BitbucketServerTriggerConfig describes the configuration of a trigger
	//  that creates a build whenever a Bitbucket Server event is received.
	BitbucketServerTriggerConfig *BitbucketServerTriggerConfig `json:"bitbucketServerTriggerConfig,omitempty"`

	// Configuration for manual approval to start a build invocation of this
	//  BuildTrigger.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.approval_config
	ApprovalConfig *ApprovalConfig `json:"approvalConfig,omitempty"`

	// Build logs will be sent back to GitHub as part of the checkrun
	//  result.  Values can be INCLUDE_BUILD_LOGS_UNSPECIFIED or
	//  INCLUDE_BUILD_LOGS_WITH_STATUS
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.include_build_logs_with_status
	IncludeBuildLogs *string `json:"includeBuildLogs,omitempty"`

	// Whether the trigger is disabled or not. If true, the trigger will never
	//  result in a build.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.BuildTrigger.disabled
	Disabled *bool `json:"disabled,omitempty"`
}

// CloudBuildTriggerStatus defines the config connector machine state of CloudBuildTrigger
type CloudBuildTriggerStatus struct {
	/* Conditions represent the latest available observations of the
	   object's current state. */
	Conditions []v1alpha1.Condition `json:"conditions,omitempty"`

	// ObservedGeneration is the generation of the resource that was most recently observed by the Config Connector controller. If this is equal to metadata.generation, then that means that the current reported status reflects the most recent desired state of the resource.
	ObservedGeneration *int64 `json:"observedGeneration,omitempty"`

	// A unique specifier for the CloudBuildTrigger resource in GCP.
	ExternalRef *string `json:"externalRef,omitempty"`

	// ObservedState is the state of the resource as most recently observed in GCP.
	ObservedState *CloudBuildTriggerObservedState `json:"observedState,omitempty"`

	// The unique identifier for the trigger.
	TriggerID *string `json:"triggerId,omitempty"`

	// Time when the trigger was created.
	CreateTime *string `json:"createTime,omitempty"`
}

// CloudBuildTriggerObservedState is the state of the CloudBuildTrigger resource as most recently observed in GCP.
// +kcc:observedstate:proto=google.devtools.cloudbuild.v1.BuildTrigger
type CloudBuildTriggerObservedState struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:categories=gcp,shortName=gcpcloudbuildtrigger;gcpcloudbuildtriggers
// +kubebuilder:subresource:status
// +kubebuilder:metadata:labels="cnrm.cloud.google.com/managed-by-kcc=true"
// +kubebuilder:metadata:labels="cnrm.cloud.google.com/system=true"
// +kubebuilder:metadata:labels="cnrm.cloud.google.com/stability-level=stable"
// +kubebuilder:metadata:labels="cnrm.cloud.google.com/tf2crd=true"
// +kubebuilder:printcolumn:name="Age",JSONPath=".metadata.creationTimestamp",type="date"
// +kubebuilder:printcolumn:name="Ready",JSONPath=".status.conditions[?(@.type=='Ready')].status",type="string",description="When 'True', the most recent reconcile of the resource succeeded"
// +kubebuilder:printcolumn:name="Status",JSONPath=".status.conditions[?(@.type=='Ready')].reason",type="string",description="The reason for the value in 'Ready'"
// +kubebuilder:printcolumn:name="Status Age",JSONPath=".status.conditions[?(@.type=='Ready')].lastTransitionTime",type="date",description="The last transition time for the value in 'Status'"

// CloudBuildTrigger is the Schema for the CloudBuildTrigger API
// +k8s:openapi-gen=true
type CloudBuildTrigger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +required
	Spec   CloudBuildTriggerSpec   `json:"spec,omitempty"`
	Status CloudBuildTriggerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// CloudBuildTriggerList contains a list of CloudBuildTrigger
type CloudBuildTriggerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CloudBuildTrigger `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CloudBuildTrigger{}, &CloudBuildTriggerList{})
}

// +kcc:proto=google.devtools.cloudbuild.v1.RepoSource
type TriggerTemplate struct {
	// The Cloud Source Repository to build. If omitted, the repo with
	//  name "default" is assumed.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.RepoSource.repo_name
	RepoRef *SourceRepoRepositoryRef `json:"repoRef,omitempty"`

	// Regex matching branches to build.
	//
	//  The syntax of the regular expressions accepted is the syntax accepted by
	//  RE2 and described at https://github.com/google/re2/wiki/Syntax
	// +kcc:proto:field=google.devtools.cloudbuild.v1.RepoSource.branch_name
	BranchName *string `json:"branchName,omitempty"`

	// Regex matching tags to build.
	//
	//  The syntax of the regular expressions accepted is the syntax accepted by
	//  RE2 and described at https://github.com/google/re2/wiki/Syntax
	// +kcc:proto:field=google.devtools.cloudbuild.v1.RepoSource.tag_name
	TagName *string `json:"tagName,omitempty"`

	// Explicit commit SHA to build.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.RepoSource.commit_sha
	CommitSha *string `json:"commitSha,omitempty"`

	// Directory, relative to the source root, in which to run the build.
	//
	//  This must be a relative path. If a step's `dir` is specified and is an
	//  absolute path, this value is ignored for that step's execution.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.RepoSource.dir
	Dir *string `json:"dir,omitempty"`

	// Only trigger a build if the revision regex does NOT match the revision
	//  regex.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.RepoSource.invert_regex
	InvertRegex *bool `json:"invertRegex,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.PubsubConfig
type TriggerPubsubConfig struct {

	// The name of the topic from which this subscription is receiving messages.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.PubsubConfig.topic
	TopicRef *pubsubv1beta1.PubSubTopicRef `json:"topicRef"`

	// Service account that will make the push request.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.PubsubConfig.service_account_email
	ServiceAccountRef *refsv1beta1.IAMServiceAccountRef `json:"serviceAccountRef,omitempty"`

	// Potential issues with the underlying Pub/Sub subscription configuration.
	//  Only populated on get requests.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.PubsubConfig.state
	State *string `json:"state,omitempty"`

	// Output only. Name of the subscription.
	Subscription *string `json:"subscription,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.WebhookConfig
type TriggerWebhookConfig struct {
	// Required. Resource name for the secret required as a URL parameter.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.WebhookConfig.secret
	SecretRef *refsv1beta1.SecretManagerSecretVersionRef `json:"secretRef"`

	// Potential issues with the underlying Pub/Sub subscription configuration.
	//  Only populated on get requests.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.WebhookConfig.state
	State *string `json:"state,omitempty"`
}

type BitbucketServerConfigRef struct {
	// A reference to an externally managed CloudBuildBitbucketServerConfig resource.
	// Should be in the format "projects/{{projectID}}/locations/{{location}}/bitbucketServerConfigs/{{bitbucketServerConfigID}}".
	External string `json:"external,omitempty"`

	// The name of a CloudBuildBitbucketServerConfig resource.
	Name string `json:"name,omitempty"`

	// The namespace of a CloudBuildBitbucketServerConfig resource.
	Namespace string `json:"namespace,omitempty"`
}

// BitbucketServerTriggerConfig describes the configuration of a trigger
//
//	that creates a build whenever a Bitbucket Server event is received.
type BitbucketServerTriggerConfig struct {
	// The full resource name of the bitbucket server config. Format:
	// projects/{project}/locations/{location}/bitbucketServerConfigs/{id}.
	BitbucketServerConfigResourceRef *BitbucketServerConfigRef `json:"bitbucketServerConfigResourceRef"`

	// Key of the project that the repo is in. For example:
	// The key for https://mybitbucket.server/projects/TEST/repos/test-repo
	// is "TEST".
	ProjectKey string `json:"projectKey"`

	// Slug of the repository. A repository slug is a URL-friendly version of a repository name, automatically generated by Bitbucket for use in the URL.
	// For example, if the repository name is 'test repo', in the URL it would become 'test-repo' as in https://mybitbucket.server/projects/TEST/repos/test-repo.
	RepoSlug string `json:"repoSlug"`

	// Filter to match changes in pull requests.
	PullRequest *PullRequestFilter `json:"pullRequest,omitempty"`

	// Filter to match changes in refs like branches, tags.
	Push *PushFilter `json:"push,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.GitFileSource
type TriggerGitFileSource struct {
	// The path of the file, with the repo root as the root of the path.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.GitFileSource.path
	Path *string `json:"path"`

	// The URI of the repo.
	//  Either uri or repository can be specified.
	//  If unspecified, the repo from which the trigger invocation originated is
	//  assumed to be the repo from which to read the specified path.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.GitFileSource.uri
	URI *string `json:"uri,omitempty"`

	// The fully qualified resource name of the Repos API repository.
	//  Either URI or repository can be specified.
	//  If unspecified, the repo from which the trigger invocation originated is
	//  assumed to be the repo from which to read the specified path.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.GitFileSource.repository
	RepositoryRef *CloudBuildV2RepositoryRef `json:"repositoryRef,omitempty"`

	// See RepoType above.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.GitFileSource.repo_type
	RepoType *string `json:"repoType"`

	// The branch, tag, arbitrary ref, or SHA version of the repo to use when
	//  resolving the filename (optional).
	//  This field respects the same syntax/resolution as described here:
	//  https://git-scm.com/docs/gitrevisions
	//  If unspecified, the revision from which the trigger invocation originated
	//  is assumed to be the revision from which to read the specified path.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.GitFileSource.revision
	Revision *string `json:"revision,omitempty"`

	// The full resource name of the github enterprise config.
	//  Format:
	//  `projects/{project}/locations/{location}/githubEnterpriseConfigs/{id}`.
	//  `projects/{project}/githubEnterpriseConfigs/{id}`.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.GitFileSource.github_enterprise_config
	GithubEnterpriseConfigRef *CloudBuildGithubEnterpriseConfigRef `json:"githubEnterpriseConfigRef,omitempty"`

	// Bitbucket Server Config resource name.
	BitbucketServerConfigRef *BitbucketServerConfigRef `json:"bitbucketServerConfigRef,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.GitRepoSource
type TriggerGitRepoSource struct {
	// The URI of the repo (e.g. https://github.com/user/repo.git).
	//  Either `uri` or `repository` can be specified and is required.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.GitRepoSource.uri
	URI *string `json:"uri,omitempty"`

	// The connected repository resource name, in the format
	//  `projects/*/locations/*/connections/*/repositories/*`. Either `uri` or
	//  `repository` can be specified and is required.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.GitRepoSource.repository
	RepositoryRef *CloudBuildV2RepositoryRef `json:"repositoryRef,omitempty"`

	// The branch or tag to use. Must start with "refs/" (required).
	// +kcc:proto:field=google.devtools.cloudbuild.v1.GitRepoSource.ref
	Ref *string `json:"ref"`

	// See RepoType below.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.GitRepoSource.repo_type
	RepoType *string `json:"repoType"`

	// The full resource name of the github enterprise config.
	//  Format:
	//  `projects/{project}/locations/{location}/githubEnterpriseConfigs/{id}`.
	//  `projects/{project}/githubEnterpriseConfigs/{id}`.
	// +kcc:proto:field=google.devtools.cloudbuild.v1.GitRepoSource.github_enterprise_config
	GithubEnterpriseConfigRef *CloudBuildGithubEnterpriseConfigRef `json:"githubEnterpriseConfigRef,omitempty"`

	// Bitbucket Server Config resource name.
	BitbucketServerConfigRef *BitbucketServerConfigRef `json:"bitbucketServerConfigRef,omitempty"`
}

type SourceRepoRepositoryRef struct {
	/* The name of the SourceRepoRepository resource. */
	Name string `json:"name,omitempty"`

	/* The namespace of the SourceRepoRepository resource. */
	Namespace string `json:"namespace,omitempty"`

	/* The external name of the SourceRepoRepository resource. */
	External string `json:"external,omitempty"`
}

type CloudBuildV2RepositoryRef struct {
	/* The name of the CloudBuildV2Repository resource. */
	Name string `json:"name,omitempty"`

	/* The namespace of the CloudBuildV2Repository resource. */
	Namespace string `json:"namespace,omitempty"`

	/* The external name of the CloudBuildV2Repository resource. */
	External string `json:"external,omitempty"`
}

type CloudBuildGithubEnterpriseConfigRef struct {
	/* The name of the CloudBuildGithubEnterpriseConfig resource. */
	Name string `json:"name,omitempty"`

	/* The namespace of the CloudBuildGithubEnterpriseConfig resource. */
	Namespace string `json:"namespace,omitempty"`

	/* The external name of the CloudBuildGithubEnterpriseConfig resource. */
	External string `json:"external,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.Build
type TriggerBuild struct {
	Artifacts        *TriggerArtifacts                `json:"artifacts,omitempty"`
	AvailableSecrets *TriggerAvailableSecrets         `json:"availableSecrets,omitempty"`
	Images           []string                         `json:"images,omitempty"`
	LogsBucketRef    *storagev1beta1.StorageBucketRef `json:"logsBucketRef,omitempty"`
	Options          *TriggerBuildOptions             `json:"options,omitempty"`
	QueueTTL         *string                          `json:"queueTtl,omitempty"`
	Secret           []TriggerSecret                  `json:"secret,omitempty"`
	Source           *TriggerSource                   `json:"source,omitempty"`
	Step             []TriggerBuildStep               `json:"step"`
	Substitutions    map[string]string                `json:"substitutions,omitempty"`
	Tags             []string                         `json:"tags,omitempty"`
	Timeout          *string                          `json:"timeout,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.Artifacts
type TriggerArtifacts struct {
	Images  []string                `json:"images,omitempty"`
	Objects *TriggerArtifactObjects `json:"objects,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.Artifacts.ArtifactObjects
type TriggerArtifactObjects struct {
	Location *string  `json:"location,omitempty"`
	Paths    []string `json:"paths,omitempty"`
	// Timing is Output only in proto, but in Spec in existing CRD?
	// CRD: timing: description: Output only. ...
	// If it is output only, why is it in Spec?
	// Maybe Terraform allows setting it? Unlikely.
	// I will add it to match CRD.
	Timing []TriggerTimeSpan `json:"timing,omitempty"`
}

type TriggerTimeSpan struct {
	StartTime *string `json:"startTime,omitempty"`
	EndTime   *string `json:"endTime,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.Secrets
type TriggerAvailableSecrets struct {
	SecretManager []TriggerSecretManagerSecret `json:"secretManager,omitempty"`
}

type TriggerSecretManagerSecret struct {
	Env        *string                                    `json:"env"`
	VersionRef *refsv1beta1.SecretManagerSecretVersionRef `json:"versionRef"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.BuildOptions
type TriggerBuildOptions struct {
	DiskSizeGb            *int64          `json:"diskSizeGb,omitempty"`
	DynamicSubstitutions  *bool           `json:"dynamicSubstitutions,omitempty"`
	Env                   []string        `json:"env,omitempty"`
	LogStreamingOption    *string         `json:"logStreamingOption,omitempty"`
	Logging               *string         `json:"logging,omitempty"`
	MachineType           *string         `json:"machineType,omitempty"`
	RequestedVerifyOption *string         `json:"requestedVerifyOption,omitempty"`
	SecretEnv             []string        `json:"secretEnv,omitempty"`
	SourceProvenanceHash  []string        `json:"sourceProvenanceHash,omitempty"`
	SubstitutionOption    *string         `json:"substitutionOption,omitempty"`
	Volumes               []TriggerVolume `json:"volumes,omitempty"`
	WorkerPool            *string         `json:"workerPool,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.Volume
type TriggerVolume struct {
	Name *string `json:"name,omitempty"`
	Path *string `json:"path,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.Secret
type TriggerSecret struct {
	KmsKeyRef *refsv1beta1.KMSCryptoKeyRef `json:"kmsKeyRef"`
	SecretEnv map[string]string            `json:"secretEnv,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.Source
type TriggerSource struct {
	RepoSource    *TriggerRepoSource    `json:"repoSource,omitempty"`
	StorageSource *TriggerStorageSource `json:"storageSource,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.RepoSource
type TriggerRepoSource struct {
	BranchName    *string                  `json:"branchName,omitempty"`
	CommitSha     *string                  `json:"commitSha,omitempty"`
	Dir           *string                  `json:"dir,omitempty"`
	InvertRegex   *bool                    `json:"invertRegex,omitempty"`
	ProjectId     *string                  `json:"projectId,omitempty"`
	RepoRef       *SourceRepoRepositoryRef `json:"repoRef"`
	Substitutions map[string]string        `json:"substitutions,omitempty"`
	TagName       *string                  `json:"tagName,omitempty"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.StorageSource
type TriggerStorageSource struct {
	BucketRef  *storagev1beta1.StorageBucketRef `json:"bucketRef"`
	Generation *string                          `json:"generation,omitempty"` // String in CRD, int64 in Proto. Existing CRD says string.
	Object     *string                          `json:"object"`
}

// +kcc:proto=google.devtools.cloudbuild.v1.BuildStep
type TriggerBuildStep struct {
	AllowExitCodes []int32         `json:"allowExitCodes,omitempty"` // integer in CRD (int32 is fine)
	AllowFailure   *bool           `json:"allowFailure,omitempty"`
	Args           []string        `json:"args,omitempty"`
	Dir            *string         `json:"dir,omitempty"`
	Entrypoint     *string         `json:"entrypoint,omitempty"`
	Env            []string        `json:"env,omitempty"`
	ID             *string         `json:"id,omitempty"`
	Name           *string         `json:"name"`
	Script         *string         `json:"script,omitempty"`
	SecretEnv      []string        `json:"secretEnv,omitempty"`
	Timeout        *string         `json:"timeout,omitempty"`
	Timing         *string         `json:"timing,omitempty"` // Output only
	Volumes        []TriggerVolume `json:"volumes,omitempty"`
	WaitFor        []string        `json:"waitFor,omitempty"`
}
