package creds

import (
	"fmt"
)

// LifetimeLimits describes the minimum, maximum, and default values for
// credential lifespan
type LifetimeLimits struct {
	Min, Max, Default int64
}

// SessionTokenLifetimeLimits describes the min, max, and default lifespan for
// the sts:GetSessionToken call
var SessionTokenLifetimeLimits = LifetimeLimits{Min: 900, Max: 3600 * 36, Default: 3600}

// AssumeRoleLifetimeLimits describes the min, max, and default lifespan for
// the sts:AssumeRole call
var AssumeRoleLifetimeLimits = LifetimeLimits{Min: 900, Max: 3600 * 12, Default: 3600}

func validateLifetime(l int64, limits LifetimeLimits) (int64, error) {
	lifetime := l
	logger.InfoMsgf("validating lifetime: %d", lifetime)
	if lifetime == 0 {
		logger.InfoMsgf("setting lifetime to default: %d", limits.Default)
		lifetime = limits.Default
	}
	if lifetime != 0 && (lifetime < limits.Min || lifetime > limits.Max) {
		return 0, fmt.Errorf("lifetime must be between %d and %d: %d",
			limits.Min,
			limits.Max,
			lifetime)
	}
	return lifetime, nil
}
