package common

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type HTTPRequestMux interface {
	Handle(pattern string, handler http.Handler)
}

// Options defines options that can be set on initializing the HelmProjectOperator
type Options struct {
	HelmApiVersion   string
	ReleaseName      string
	SystemNamespaces []string
	ChartContent     string
	Singleton        bool

	// HTTPRequestMux is the HTTP request multiplexer to register the /healthz endpoint at
	// If it is not provided, it is assumed the mux is http.DefaultServeMux
	HTTPRequestMux HTTPRequestMux

	ProjectLabel            string
	SystemProjectLabelValue string
	SystemDefaultRegistry   string
	CattleURL               string
	ClusterID               string
	NodeName                string
	AdminClusterRole        string
	EditClusterRole         string
	ViewClusterRole         string

	HelmJobImage string
}

func (opts Options) Validate() error {
	if len(opts.HelmApiVersion) == 0 {
		return errors.New("must provide a spec.helmApiVersion that this project operator is being initialized for")
	}

	if len(opts.ReleaseName) == 0 {
		return errors.New("must provide name of Helm release that this project operator should deploy")
	}

	if len(opts.SystemNamespaces) > 0 {
		logrus.Infof("Marking the following namespaces as system namespaces: %s", opts.SystemNamespaces)
	}

	if len(opts.ChartContent) == 0 {
		return errors.New("cannot instantiate Project Operator without bundling a Helm chart to provide for the HelmChart's spec.ChartContent")
	}

	if opts.Singleton {
		logrus.Infof("Note: Operator only supports a single ProjectHelmChart per project registration namespace")
		if len(opts.ProjectLabel) == 0 {
			logrus.Warnf("It is only recommended to run a singleton Project Operator when --project-label is provided (currently not set). The current configuration of this operator would only allow a single ProjectHelmChart to be managed by this Operator.")
		}
	}

	if len(opts.ProjectLabel) > 0 {
		logrus.Infof("Creating dedicated project registration namespaces to discover ProjectHelmCharts based on the value found for the project label %s on all namespaces in the cluster, excluding system namespaces; these namespaces will need to be manually cleaned up if they have the label '%s: \"true\"'", opts.ProjectLabel, HelmProjectOperatedNamespaceOrphanedLabel)
		if len(opts.SystemProjectLabelValue) > 0 {
			logrus.Infof("assuming namespaces tagged with %s=%s are also system namespaces", opts.ProjectLabel, opts.SystemProjectLabelValue)
		}
		if len(opts.ClusterID) > 0 {
			logrus.Infof("Marking project registration namespaces with %s=%s:<projectID>", opts.ProjectLabel, opts.ClusterID)
		}
	}

	if len(opts.HelmJobImage) > 0 {
		logrus.Infof("Using %s as spec.JobImage on all generated HelmChart resources", opts.HelmJobImage)
	}

	if len(opts.NodeName) > 0 {
		logrus.Infof("Marking events as being sourced from node %s", opts.NodeName)
	}

	for subjectRole, defaultClusterRoleName := range GetDefaultClusterRoles(opts) {
		logrus.Infof("RoleBindings will automatically be created for Roles in the Project Release Namespace marked with '%s': '<helm-release>' "+
			"and '%s': '%s' based on ClusterRoleBindings or RoleBindings in the Project Registration namespace tied to ClusterRole %s",
			HelmProjectOperatorProjectHelmChartRoleLabel, HelmProjectOperatorProjectHelmChartRoleAggregateFromLabel, subjectRole, defaultClusterRoleName,
		)
	}

	return nil
}
