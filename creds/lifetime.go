package creds

import (
	"fmt"
)

// Constants for the session duration parameter in calls to sts:AssumeRole and
// sts:GetSessionToken.
const (
	SessionLifetimeMin     = 900
	SessionLifetimeMax     = 3600 * 12
	SessionLifetimeDefault = 3600
)

func validateLifetime(lifetime int64) error {
	logger.InfoMsgf("validating lifetime: %d", lifetime)
	if lifetime != 0 && (lifetime < SessionLifetimeMin || lifetime > SessionLifetimeMax) {
		return fmt.Errorf("lifetime must be between %d and %d: %d",
			SessionLifetimeMin,
			SessionLifetimeMax,
			lifetime)
	}
	return nil
}
