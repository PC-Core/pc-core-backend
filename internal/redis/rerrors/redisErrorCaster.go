package rerrors

import (
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/redis/go-redis/v9"
)

func RedisErrorCaster(err error) errors.PCCError {
	internal, ok := err.(redis.Error)

	if !ok {
		return errors.NewInternalSecretError()
	}

	return NewRedisError(internal)
}
