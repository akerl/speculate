package executors

import (
	"github.com/akerl/speculate/creds"
)

// Executor defines the interface for requesting a new set of AWS creds
type Executor interface {
	Execute() (creds.Creds, error)
	ExecuteWithCreds(creds.Creds) (creds.Creds, error)
}
