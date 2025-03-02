package rerrors

import (
	"fmt"

	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/redis/go-redis/v9"
)

const (
	RE_DEFAULT_MESSAGE_FORMAT = "Redis error with code: %d"
	RE_MESSAGE_WRONG_VALUE    = "Redis returned the wrong value"

	RE_SAFE_MESSAGE = "Redis error occured"
)

// RedisError represents the redis error
type RedisError struct {
	Inner   redis.Error
	Code    errors.ErrorCode
	Kind    errors.ErrorKind
	Message string
}

func NewRedisError(inner redis.Error) *RedisError {
	ec := parseErrorCode(inner)
	return &RedisError{
		inner, ec, errors.EK_REDIS, fmt.Sprintf(RE_DEFAULT_MESSAGE_FORMAT, ec),
	}
}

func NewRedisErrorWrongValue() *RedisError {
	return &RedisError{
		nil, errors.EC_REDIS_WRONG_VALUE, errors.EK_REDIS, RE_MESSAGE_WRONG_VALUE,
	}
}

func (r *RedisError) Error() string {
	return r.Message
}

func (r *RedisError) GetErrorCode() errors.ErrorCode {
	return r.Code
}

func (r *RedisError) GetErrorKind() errors.ErrorKind {
	return r.Kind
}

func (r *RedisError) IntoPublic() *errors.PublicPCCError {
	return errors.NewPublicPCCError(r.Code, r.Kind, nil, RE_SAFE_MESSAGE)
}

func parseErrorCode(err redis.Error) errors.ErrorCode {
	switch err {
	case redis.Nil:
		return errors.EC_REDIS_NIL
	default:
		return errors.EC_REDIS_OTHER
	}
}
