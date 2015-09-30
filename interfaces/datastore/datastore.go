/*
Package datastore contains an fake implementation of a SQL database
that (for whatever reason) needs to talk to the session manager, but we don't want to
import that package for a couple reasons:

1. The session manager implementation also needs to use the database, and Go doesn't allow
   cyclic imports
2. It's just bad form to tightly couple your concrete implementations - you'd get
   bit by a different problem if cyclic imports *were* allowed.

Don't look too far into why the classes are talking with each other - it's a
contrived example to demonstrate Go's implicit interfaces.
*/
package datastore

import (
	"fmt"
)

// SessionManager defines the methods we need to access from our session manager, without
// knowing or caring what the concrete type is. We're interfacing out a type without telling it!
type SessionManager interface {
	// Validate a session key, returning the username and user id
	ValidateSession(sessionKey string) (string, int64, error)
}

// SQLDatabase is a data store that uses SQL
type SQLDatabase struct {
	sessionManager SessionManager // interface that defines the methods we need from our Session Manager
	sqlHost        string
	sqlLogin       string
	sqlPassword    string
}

// NewSQLDatabase builds a new SQLDatabase
// Note that we don't pass in the sessionManager - that's because our session manager and data store
// need each other. The top-level code that wires them together should set each other's properties
// immediately after instantiating them.
func NewSQLDatabase(sqlHost string, sqlLogin string, sqlPassword string) (*SQLDatabase, error) {
	ret := SQLDatabase{
		sqlHost:     sqlHost,
		sqlLogin:    sqlLogin,
		sqlPassword: sqlPassword,
	}
	return &ret, nil
}

// SetSessionManager sets the SessionManager
func (db *SQLDatabase) SetSessionManager(sessionManager SessionManager) error {
	db.sessionManager = sessionManager
	return nil
}

// ValidateUser checks if a username/password is valid, returning the unique user ID
func (db *SQLDatabase) ValidateUser(username string, password string) (int64, error) {
	// TODO: stubbed out
	return 42, nil
}

// LockOutUser records that a user account can no longer be accessed, due to bad login attempts
func (db *SQLDatabase) LockOutUser(username string) error {
	// TODO: stubbed out
	return nil
}

// ChangePassword changes a user's password, making sure the person is allowed to do this.
// Note: this is a bad example - the DB manager shouldn't be reaching out to the session manager,
// but I need an example, so...
func (db *SQLDatabase) ChangePassword(sessionKey string, username string, newPassword string) error {
	if db.sessionManager == nil {
		panic("You forgot to set the SessionManager!")
	}

	// make sure the user is logged in
	authenticatedUsername, _, authErr := db.sessionManager.ValidateSession(sessionKey)
	if authErr != nil {
		return fmt.Errorf("Error authenticating user %s by session key: %s", username, authErr)
	}

	// make sure the user is changing own password
	if username != authenticatedUsername {
		return fmt.Errorf("Authenticated user %s: trying to change password for %s", authenticatedUsername, username)
	}

	// TODO: do actual database stuff
	return nil
}
