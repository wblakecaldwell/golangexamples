/*
Package session contains an fake implementation of a Redis session manager
that (for whatever reason) needs to talk to the Database, but we don't want to
import that package for a couple reasons:

1. The Database implementation also needs to use the session manager, and Go
   doesn't allow cyclic imports
2. It's just bad form to tightly couple your concrete implementations - you'd get
   bit by a different problem if cyclic imports *were* allowed.

Don't look too far into why the classes are talking with each other - it's a
contrived example to demonstrate Go's implicit interfaces.
*/
package session

import (
	"fmt"
)

// DataStore defines the methods we need to access from our DataStore, without
// knowing or caring what the concrete type is. We're interfacing out a type
// without telling that it!
type DataStore interface {
	// ValidateUser checks if a username/password is valid, returning the unique user ID
	ValidateUser(username string, password string) (int64, error)

	// LockOutUser records that a user account can no longer be accessed, due to bad login attempts
	LockOutUser(username string) error
}

// RedisSessionManager is a session implementation that validates logins with the database, and stores the
// session info in Redis
type RedisSessionManager struct {
	dataStore     DataStore // avoid coupling to the concrete Database type by using the DataStore interface
	redisHost     string    // internal Redis stuff
	redisLogin    string    // internal Redis stuff
	redisPassword string    // internal Redis stuff
}

// NewRedisSessionManager creates a new RedisSessionManager
// Note that we don't pass in the sessionManager - that's because our session manager and data store
// need each other. The top-level code that wires them together should set each other's properties
// immediately after instantiating them.
func NewRedisSessionManager(redisHost string, redisLogin string, redisPassword string) (*RedisSessionManager, error) {
	ret := RedisSessionManager{
		redisHost:     redisHost,
		redisLogin:    redisLogin,
		redisPassword: redisPassword,
	}
	return &ret, nil
}

// SetDataStore sets the DataStore
func (rs *RedisSessionManager) SetDataStore(dataStore DataStore) error {
	rs.dataStore = dataStore
	return nil
}

// AuthenticateUser authenticates the input user, returning the session key, and
// error when not valid or a problem occurred.
func (rs *RedisSessionManager) AuthenticateUser(username string, password string) (string, error) {
	if rs.dataStore == nil {
		panic("You forgot to set the DataStore!")
	}

	_, authErr := rs.dataStore.ValidateUser(username, password)
	if authErr != nil {
		// this system doesn't allow a single bad login, because: reasons
		if lockOutErr := rs.dataStore.LockOutUser(username); lockOutErr != nil {
			// TODO: log this
		}
		return "", fmt.Errorf("Error validating user %s: %s", username, authErr)
	}

	// valid login - TODO: actually store this in Redis

	return "your-session-key", nil
}

// ValidateSession checks if the input session key is a valid logged-in user, and returns the
// username and id if so.
func (rs *RedisSessionManager) ValidateSession(sessionKey string) (string, int64, error) {
	// TODO: this is just stubbed-out
	return "fprefect", 42, nil
}
