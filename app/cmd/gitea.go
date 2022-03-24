package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

var (
    giteaCmd = &cobra.Command{
        Use: "gitea",
        Short: "Install gitea.",
        Long: `Install gitea using go helm driver.`,
        
        Run: func(cmd *cobra.Command, args []string) {
            helmInstall()
        },
    }
    dryrun bool
    chartPath string
)

func init() {
	rootCmd.AddCommand(giteaCmd)
    giteaCmd.PersistentFlags().BoolVar(&dryrun, "dry-run", false, "flags whether to run the command or simply mock run it.")
    rootCmd.AddCommand(giteaCmd)
    giteaCmd.PersistentFlags().StringVar(&chartPath, "chartPath", "", "Absolute path to the file you wish to load.")
}

//
// helmInstall installs gitea from a local helm chart.
// TODO: Wait for helm to fix parsing of values file, 
// as currently values do not get parsed properly.
func helmInstall() error {
    namespace := "gitea"
    releaseName := "gitea"

    settings := cli.New()
    settings.SetNamespace(namespace)
    actionConfig := new(action.Configuration)

    if err := actionConfig.Init(settings.RESTClientGetter(), namespace,
        os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
        logerror.Printf("%+v", err)
        return err
    }

    // load values specific to this helm installation
    vals, err := generateGiteaValues()

    if err != nil {
        logerror.Println(err)
        return err
    }

    // load chart from cli chartPath arg
    chart, err := loader.Load(chartPath)
    if err != nil {
        logerror.Println(err)
        return err
    }

    client := action.NewInstall(actionConfig)
    client.CreateNamespace = true
    client.Namespace = namespace
    client.ReleaseName = releaseName
    client.DryRun = dryrun

    // install the chart here
    rel, err := client.Run(chart, vals)
    if err != nil {
        logerror.Println(err)
        return err
    }

    log.Printf("Installed Chart with name: %s in namespace: %s from path: %s\n", rel.Name, rel.Namespace, chartPath)

    return nil
}

func generateGiteaValues() (map[string]interface{}, error) {
    loginfo.Println("Attemping to load values for gitea.")
    out := make(map[string]interface{})
    if err := yaml.Unmarshal([]byte(giteaValuesYaml), &out); err != nil {
        logerror.Println(err)
        return nil, err
    }
    loginfo.Printf("Parsed obj. :\n\t -> %v\n", out)
    return out, nil
}



var giteaValuesYaml string = `
gitea:
  admin:
    username: superuser
    password: password
    email: "gitea@local.domain"

  metrics:
    enabled: false
    serviceMonitor:
      enabled: false

  config: {}

postgresql:
  enabled: true
  global:
    postgresql:
      postgresqlDatabase: gitea
      postgresqlUsername: gitea
      postgresqlPassword: gitea
      servicePort: '5432'
  persistence:
    size: 1Gi

memcached:
  enabled: false
  service:
    port: 11211

checkDeprecation: true
`