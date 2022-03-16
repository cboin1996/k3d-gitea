package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	psqlSecret string
	rootCmd = &cobra.Command{
		Use:   "k3d-gitea",
		Short: "helm installation of gitea with go configuration of postgres",
		Long: `helm installation of gitea with go configuration of postgres`,
	}
	loginfo = log.New(os.Stdout, "[INFO]", log.Ldate|log.Ltime|log.Lshortfile)
	logerror = log.New(os.Stderr, "[ERROR]", log.Ldate|log.Ltime|log.Lshortfile)
)
  
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	loginfo.Println("Initializing.")
	rootCmd.PersistentFlags().StringVar(&psqlSecret, "psqlsecret", "", "postgres root user secret.")
	rootCmd.MarkPersistentFlagRequired("psqlsecret")
}
