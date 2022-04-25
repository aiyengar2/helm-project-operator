package main

import (
	_ "embed"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/aiyengar2/helm-project-operator/pkg/controllers/common"
	"github.com/aiyengar2/helm-project-operator/pkg/operator"
	"github.com/aiyengar2/helm-project-operator/pkg/version"
	command "github.com/rancher/wrangler-cli"
	_ "github.com/rancher/wrangler/pkg/generated/controllers/apiextensions.k8s.io"
	_ "github.com/rancher/wrangler/pkg/generated/controllers/networking.k8s.io"
	"github.com/rancher/wrangler/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

const (
	DummyHelmApiVersion = "dummy.cattle.io/v1alpha1"
	DummyReleaseName    = "dummy"
)

var (
	DummySystemNamespaces = []string{"kube-system"}

	//go:embed bin/example-chart/example-chart.tgz.base64
	base64TgzChart string

	debugConfig command.DebugConfig
)

type DummyOperator struct {
	// Note: all Project Operator are expected to provide these RuntimeOptions
	common.RuntimeOptions

	Kubeconfig string `usage:"Kubeconfig file"`
}

func (o *DummyOperator) Run(cmd *cobra.Command, args []string) error {
	go func() {
		// required to set up healthz and pprof handlers
		log.Println(http.ListenAndServe("localhost:80", nil))
	}()
	debugConfig.MustSetupDebug()

	cfg := kubeconfig.GetNonInteractiveClientConfig(o.Kubeconfig)

	ctx := cmd.Context()

	if err := operator.Init(ctx, o.Namespace, cfg, common.Options{
		// These fields are provided by the Project Operator
		HelmApiVersion:   DummyHelmApiVersion,
		ReleaseName:      DummyReleaseName,
		SystemNamespaces: DummySystemNamespaces,
		ChartContent:     base64TgzChart,
		Singleton:        false,

		// These fields are provided on runtime for all project operators
		ProjectLabel:            o.ProjectLabel,
		SystemProjectLabelValue: o.SystemProjectLabelValue,
		SystemDefaultRegistry:   o.SystemDefaultRegistry,
		CattleURL:               o.CattleURL,
		ClusterID:               o.ClusterID,
		NodeName:                o.NodeName,
	}); err != nil {
		return err
	}

	<-cmd.Context().Done()
	return nil
}

func main() {
	cmd := command.Command(&DummyOperator{}, cobra.Command{
		Version: version.FriendlyVersion(),
	})
	cmd = command.AddDebug(cmd, &debugConfig)
	command.Main(cmd)
}
