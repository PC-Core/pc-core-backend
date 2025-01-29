package rerrors

import (
	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/redis/go-redis/v9"
)

func RedisErrorCaster(err error) errors.PCCError {
	internal, ok := err.(redis.Error)

	if !ok {
		return errors.NewInternalSecretError()
	}

	return NewRedisError(internal)
}
