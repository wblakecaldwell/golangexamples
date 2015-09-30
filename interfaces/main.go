/*
Package main wires together two structures that reference each other through
interfaces so they don't create a cyclic package dependency. The different
structures don't need to know what the concrete types are, so they each
define their own interface around the other.
*/
package main

import (
	"fmt"
	"github.com/wblakecaldwell/golangexamples/interfaces/datastore"
	"github.com/wblakecaldwell/golangexamples/interfaces/session"
	"os"
)

func main() {
	// create the database
	db, err := datastore.NewSQLDatabase("localhost", "sqluser", "hunter2")
	if err != nil {
		fmt.Printf("Fatal error: can't create SQLDatabase: %s\n", err)
		os.Exit(1)
	}

	// create the session manager
	sessionManager, err := session.NewRedisSessionManager("localhost", "redisuser", "hunter2")
	if err != nil {
		fmt.Printf("Fatal error: can't create RedisSessionManager: %s\n", err)
		os.Exit(1)
	}

	// tie the two together
	db.SetSessionManager(sessionManager)
	sessionManager.SetDataStore(db)

	// TODO: start the app!
	fmt.Println("Success!")
}
