package creds

import (
	"fmt"
)

func validateLifetime(lifetime int64) error {
	logger.InfoMsg(fmt.Sprintf("validating lifetime: %d", lifetime))
	if lifetime != 0 && (lifetime < 900 || lifetime > 3600) {
		return fmt.Errorf("lifetime must be between 900 and 3600: %d", lifetime)
	}
	return nil
}
