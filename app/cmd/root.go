package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "k3d-gitea",
		Short: "cli for gitea installation and configuration for use with postgres.",
		Long: `cli for gitea installation and configuration for use with postgres.`,
	}
	loginfo = log.New(os.Stdout, "[INFO]", log.Ldate|log.Ltime|log.Lshortfile)
	logerror = log.New(os.Stderr, "[ERROR]", log.Ldate|log.Ltime|log.Lshortfile)

)
  
func Execute() error {
	return rootCmd.Execute()
}
