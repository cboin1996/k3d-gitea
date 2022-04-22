package cmd

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)
var host string
var port int
var dbname string
var rbypass bool
var psqluname string
var psqlpasswd string
var rolename string
var rolepasswd string
var rdbname string

func init() {
	rootCmd.AddCommand(psqlinit)
	psqlinit.PersistentFlags().StringVar(&host, "host","localhost", "postgres host")
	psqlinit.MarkPersistentFlagRequired("host")
	psqlinit.PersistentFlags().IntVar(&port, "port", 5432, "postgres server port")
	psqlinit.MarkPersistentFlagRequired("port")
	psqlinit.PersistentFlags().StringVar(&dbname, "dbname", "postgres", "postgres server dbname")
	psqlinit.PersistentFlags().StringVar(&psqluname, "psqluname","postgres", "postgres root username (only the username, no username@host required.")
	psqlinit.MarkPersistentFlagRequired("psqluname")
	psqlinit.PersistentFlags().StringVar(&psqlpasswd, "psqlpasswd", "postgres", "postgres root user password.")
	psqlinit.MarkPersistentFlagRequired("psqlpasswd")
	psqlinit.PersistentFlags().BoolVar(&rbypass, "rbypass", false, "bypass gitea role creation.")
	psqlinit.MarkPersistentFlagRequired("rolename")
	psqlinit.PersistentFlags().StringVar(&rolename, "rolename", "", "the role to create")
	psqlinit.MarkPersistentFlagRequired("rolepasswd")
	psqlinit.PersistentFlags().StringVar(&rolepasswd, "rolepasswd", "", "the role to create")
	psqlinit.MarkPersistentFlagRequired("rdbname")
	psqlinit.PersistentFlags().StringVar(&rdbname, "rdbname", "", "the db to create")
}

var psqlinit = &cobra.Command{
	Use: "psqlinit",
	Short: "Initialize postgres.",
	Long: "Initialize postgres to be configured for use with gitea in a k8s cluster.",
	
	Run: func(cmd *cobra.Command, args []string) {
		db, err := connect(host, port, psqluname, psqlpasswd, dbname)
		if err == nil {
			if !rbypass {
				createRole(db, rolename, rolepasswd) 
				grantRole(db, rolename, psqluname)
			}

			initializeDb(db, rdbname, rolename)
		}
	},
}

/*
Establish a connection with the database.
*/
func connect(host string, port int, psqluname string, psqlpasswd string, dbname string) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s@%s:%s@%s:%d/%s?sslmode=require", 
		psqluname, host, url.QueryEscape(psqlpasswd), host, port, dbname)
	loginfo.Printf("Attempting to establish connection " +
			   "using postgres driver with con. string: %s\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logerror.Printf("Connection failed, error: %s\n", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		logerror.Printf("Connection failed, error: %s", err)
		return nil, err
	}
	loginfo.Println("Connection successful.")
	return db, err
}

func createRole(db *sql.DB, role string, password string) error {
	query := fmt.Sprintf("CREATE ROLE %s WITH LOGIN PASSWORD '%s'", role, password)
	loginfo.Printf("Attempting query: %s", query)
	rows, err := db.Query(query)
	if err != nil {
		logerror.Printf("Error returned for query '%s': %v\n", query, err)
		return err
	}
	loginfo.Printf("Recieved rows for query '%s': \n\t - %v\n", query, rows)
	return nil
}

// Adds role1 as a member of role2
func grantRole(db *sql.DB, role1 string, role2 string) error {
	query := fmt.Sprintf("GRANT %s to %s", role1, role2)
	loginfo.Printf("Attempting query: %s", query)
	rows, err := db.Query(query)
	if err != nil {
		logerror.Printf("Error returned for query %s: %v\n", query, err)
		return err
	}
	loginfo.Printf("Recieved rows for query %s: \n\t - %v\n", query, rows)
	return nil
}

func initializeDb(db *sql.DB, dbname string, role string) error {
	query := fmt.Sprintf(
		`CREATE DATABASE %s WITH OWNER %s TEMPLATE template0 ENCODING UTF8 LC_COLLATE "en-US" LC_CTYPE "en-US";`, dbname, role)
	loginfo.Printf("Attempting query: %s", query)
	rows, err := db.Query(query)
	if err != nil {
		logerror.Printf("Error returned for query %s: %v", query, err)
		return err
	}
	loginfo.Printf("Recieved rows for query '%s': \n\t - %v\n", query, rows)
	return nil
}
