package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PC-Core/pc-core-backend/internal/auth"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/models"
	"github.com/PC-Core/pc-core-backend/internal/redis/rerrors"
	"github.com/redis/go-redis/v9"
)

const (
	UserIDKey = "userID"
	//ExpTimeHrs = 24 * 30 * time.Hour
)

const IntErrorCode = 0xffffffffffffffff

type RedisController struct {
	client *redis.Client
}

func NewRedisController(client *redis.Client) *RedisController {
	return &RedisController{
		client,
	}
}

func (c *RedisController) GetNextID() (uint64, errors.PCCError) {
	id := c.client.Incr(context.Background(), UserIDKey)

	if err := id.Err(); err != nil {
		return IntErrorCode, rerrors.RedisErrorCaster(err)
	}

	value := id.Val()

	if value < 0 {
		return IntErrorCode, rerrors.NewRedisErrorWrongValue()
	}

	return uint64(value), nil
}

func (c *RedisController) CreateTempUser(auth auth.Auth) (interface{}, errors.PCCError) {
	id, err := c.GetNextID()

	if err != nil {
		return nil, err
	}

	tu := models.NewPublicUser(int(id), "", "", models.Temporary)
	b, jerr := json.Marshal(tu)

	if jerr != nil {
		err = errors.NewInternalSecretError()
	}

	if err != nil {
		return nil, err
	}

	dur, err := c.getUserIDTTL()

	if err != nil {
		return nil, err
	}

	rerr := c.client.Set(context.Background(), fmt.Sprintf("user:%d", id), b, dur).Err()

	if rerr != nil {
		return nil, rerrors.RedisErrorCaster(rerr)
	}

	return auth.AuthentificateWithDur(tu, dur, dur)
}

func (c *RedisController) GetCart(user_id uint64) ([]models.TempCartItem, errors.PCCError) {
	cartobj := c.client.Get(context.Background(), fmt.Sprintf("cart:%d", user_id))

	err := cartobj.Err()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, rerrors.RedisErrorCaster(err)
	}

	cart := make([]models.TempCartItem, 0)

	err = json.Unmarshal([]byte(cartobj.Val()), &cart)

	if err != nil {
		return nil, errors.NewJsonUnmarshalError()
	}

	return cart, nil
}

func (c *RedisController) CreateCartAndPut(user_id uint64, product_id uint64, quantity uint) (uint64, errors.PCCError) {
	cart := []models.TempCartItem{
		*models.NewTempCartItem(product_id, quantity),
	}

	json_cart, jerr := json.Marshal(cart)

	if jerr != nil {
		return IntErrorCode, errors.NewJsonMarshalError()
	}

	ttl, err := c.getUserIDTTL()

	if err != nil {
		return IntErrorCode, err
	}

	res := c.client.Set(context.Background(), fmt.Sprintf("cart:%d", user_id), string(json_cart), ttl)

	if err := res.Err(); err != nil {
		return IntErrorCode, nil
	}

	return product_id, nil
}

func (c *RedisController) AddToCart(user_id uint64, product_id uint64, quantity uint) (uint64, errors.PCCError) {
	record := fmt.Sprintf("cart:%d", user_id)

	tu := c.client.Get(context.Background(), record)

	err := tu.Err()

	if err == redis.Nil {
		return c.CreateCartAndPut(user_id, product_id, quantity)
	}

	if err != nil {
		return IntErrorCode, rerrors.RedisErrorCaster(err)
	}

	var cart []models.TempCartItem

	err = json.Unmarshal([]byte(tu.Val()), &cart)

	if err != nil {
		return IntErrorCode, errors.NewJsonUnmarshalError()
	}

	ttl, terr := c.getUserIDTTL()

	if terr != nil {
		return IntErrorCode, terr
	}

	if !c.checkCartForCollisionsAndAppend(cart, product_id, quantity) {
		cart = append(cart, *models.NewTempCartItem(product_id, quantity))
	}

	newcart, err := json.Marshal(cart)

	if err != nil {
		return IntErrorCode, errors.NewJsonMarshalError()
	}

	err = c.client.Set(context.Background(), record, newcart, ttl).Err()

	if err != nil {
		return IntErrorCode, rerrors.RedisErrorCaster(err)
	}

	return product_id, nil
}

func (c *RedisController) appendToCart(item *models.TempCartItem, quantity uint) bool {
	item.Quantity += quantity
	return true
}

func (c *RedisController) checkCartForCollisionsAndAppend(cart []models.TempCartItem, product_id uint64, quantity uint) bool {
	for i := range cart {
		if cart[i].ProductID == product_id {
			return c.appendToCart(&cart[i], quantity)
		}
	}
	return false
}

func (c *RedisController) getUserIDTTL() (time.Duration, errors.PCCError) {
	res := c.client.TTL(context.Background(), UserIDKey)

	if err := res.Err(); err != nil {
		return time.Duration(0), rerrors.RedisErrorCaster(err)
	}

	return res.Val(), nil
}
