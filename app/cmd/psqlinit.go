package cmd

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var rbypass bool
var psqlSecret string

func init() {
	rootCmd.AddCommand(psqlinit)
	psqlinit.PersistentFlags().StringVar(&psqlSecret, "psqlsecret", "", "postgres root user secret.")
	psqlinit.MarkPersistentFlagRequired("psqlsecret")
	psqlinit.PersistentFlags().BoolVar(&rbypass, "rbypass", false, "bypass gitea role creation.")
	
}

var psqlinit = &cobra.Command{
	Use: "psqlinit",
	Short: "Initialize postgres.",
	Long: "Initialize postgres to be configured for use with gitea in a k8s cluster.",
	
	Run: func(cmd *cobra.Command, args []string) {
		db, err := connect()
		if err == nil {
			initializeDb(db)
		}
	},
}

/*
Establish a connection with the database.
*/
func connect() (*sql.DB, error) {
	// make connection string using password from cli
	connStr := fmt.Sprintf("user=postgres password=%s sslmode=disable", psqlSecret) 
	loginfo.Printf("Attempting to establish connection" +
			   "using postgres driver with con. string: %s\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logerror.Println(fmt.Sprintf("Connection failed, error: %s", err))
		return nil, nil
	}
	loginfo.Println("Connection successful.")
	return db, err
}

func initializeDb(db *sql.DB) {
	if !rbypass {
		rows, err := db.Query("CREATE ROLE gitea WITH LOGIN PASSWORD 'gitea'")
		if err != nil {
			logerror.Println("Error while creating gitea role: ", err)
			return
		}
		loginfo.Print("Recieved rows while creating gitea: \n\t - ", rows)
		loginfo.Println("Role creation was a success!")
	} else {
		loginfo.Println("Bypassing role creation.")
	}
	loginfo.Println("Attempting to create giteadb!")
	giteaDbCreateQuery  := `CREATE DATABASE giteadb WITH OWNER gitea TEMPLATE template0 ENCODING UTF8 LC_COLLATE 'en_US.UTF-8' LC_CTYPE 'en_US.UTF-8';`
	rows, err := db.Query(giteaDbCreateQuery)
	if err != nil {
		logerror.Println("Error while creating giteadb: ", err)
		return
	}
	loginfo.Print("Recieved rows while creating giteadb: \n\t - ", rows)
	loginfo.Println("giteadb creation was a success!")
}
