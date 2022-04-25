package common

type RuntimeOptions struct {
	// Namespace is the systemNamespace to create HelmCharts and HelmReleases in
	// It's generally expected that this namespace is not widely accessible by all users in your cluster; it's recommended that it is placed
	// in something akin to a System Project that is locked down in terms of permissions since resources like HelmCharts and HelmReleases are deployed there
	Namespace string `usage:"Namespace to create HelmCharts and HelmReleases; if ProjectLabel is not provided, this will also be the namespace to watch ProjectHelmCharts" default:"cattle-helm-system" env:"NAMESPACE"`

	// NodeName is the name of the node running the operator; it adds additional information to events about where they were generated from
	NodeName string `usage:"Name of the node this controller is running on" env:"NODE_NAME"`

	// HelmJobImage is the job image to use to run the HelmChart job (default rancher/klipper-helm:v0.7.0-build20220315)
	// Generally, this HelmJobImage can be left undefined, but may be necessary to be set if you are running with a non-default image
	HelmJobImage string `usage:"Job image to use to perform helm operations on HelmChart creation" env:"HELM_JOB_IMAGE"`

	// ClusterID identifies the cluster that the operator is being operated frmo within; it adds an additional annotation to project registration
	// namespaces that indicates the projectID with the cluster label.
	//
	// Note: primarily used for integration with Rancher Projects
	ClusterID string `usage:"Identifies the cluster this controller is running on. Ignored if --project-label is not provided." env:"CLUSTER_ID"`

	// SystemDefaultRegistry is the prefix to be added to all images deployed by the HelmChart embedded into the Project Operator
	// to point at the right set of images that need to be deployed. This is usually provided in Rancher as global.cattle.systemDefaultRegistry
	SystemDefaultRegistry string `usage:"Default system registry to use for Docker images deployed by underlying Helm Chart. Provided as global.cattle.systemDefaultRegistry in the Helm Chart" env:"SYSTEM_DEFAULT_REGISTRY"`

	// CattleURL is the Rancher URL that this chart has been deployed onto. This is usually provided in Rancher Helm charts as global.cattle.url
	CattleURL string `usage:"Default Rancher URL to provide to the Helm chart under global.cattle.url" env:"CATTLE_URL"`

	// ProjectLabel is the label that identifies projects
	// Note: this field is optional and ensures that ProjectHelmCharts auto-infer their spec.projectNamespaceSelector
	// If provided, any spec.projectNamespaceSelector provided will be ignored
	// example: field.cattle.io/projectId
	ProjectLabel string `usage:"Label on namespaces to create Project Registration Namespaces and watch for ProjectHelmCharts" env:"PROJECT_LABEL"`

	// SystemProjectLabelValue is the value of the ProjectLabel that identifies system namespaces. Does nothing if ProjectLabel is not provided
	// example: p-ranch
	// If both this and the above example are provided, any namespaces with label 'field.cattle.io/projectId: p-ranch' will be treated
	// as a systemNamespace, which means that no ProjectHelmChart will be allowed to select it
	SystemProjectLabelValue string `usage:"Value on project label on namespaces that marks it as a system namespace" env:"SYSTEM_PROJECT_LABEL_VALUE"`

	// DisableRBACAggregation disables starting up the RBAC controllers that automaticaly create RoleBindings on the default Operator Roles deployed
	// onto each Project Registration namespace. By default, aggregation is turned on.
	DisableRBACAggregation bool `usage:"Disables RBAC aggregation of Kubernetes Default Roles on default operator roles deployed onto Project Registration Namespaces" env:"DISABLE_RBAC_AGGREGATION"`
}
