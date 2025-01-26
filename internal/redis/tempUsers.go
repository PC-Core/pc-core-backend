package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Core-Mouse/cm-backend/internal/auth"
	"github.com/Core-Mouse/cm-backend/internal/models"
	"github.com/redis/go-redis/v9"
)

const (
	UserIDKey  = "userID"
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

func (c *RedisController) GetNextID() (uint64, error) {
	id := c.client.Incr(context.Background(), UserIDKey)

	if err := id.Err(); err != nil {
		return IntErrorCode, err
	}

	value := id.Val()

	if value < 0 {
		return IntErrorCode, fmt.Errorf("redis returned negative value")
	}

	return uint64(value), nil
}

func (c *RedisController) CreateTempUser(auth auth.Auth) (interface{}, error) {
	id, err := c.GetNextID()

	if err != nil {
		return nil, err
	}

	tu := models.NewPublicUser(int(id), "", "", models.Temporary)
	b, err := json.Marshal(tu)

	if err != nil {
		return nil, err
	}

	dur, err := c.getUserIDTTL();

	if err != nil {
		return nil, err
	}

	err = c.client.Set(context.Background(), fmt.Sprintf("user:%d", id), b, dur).Err()

	if err != nil {
		return nil, err
	}

	return auth.AuthentificateWithDur(tu, dur, dur)
}

func (c *RedisController) GetCart(user_id uint64) ([]models.TempCartItem, error) {
	cartobj := c.client.Get(context.Background(), fmt.Sprintf("cart:%d", user_id))

	err := cartobj.Err()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	cart := make([]models.TempCartItem, 0)

	err = json.Unmarshal([]byte(cartobj.Val()), &cart)

	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (c *RedisController) CreateCartAndPut(user_id uint64, product_id uint64, quantity uint) (uint64, error) {
	cart := []models.TempCartItem{
		*models.NewTempCartItem(product_id, quantity),
	}

	json_cart, err := json.Marshal(cart)

	if err != nil {
		return IntErrorCode, err
	}

	ttl, err := c.getUserIDTTL();

	if err != nil {
		return IntErrorCode, err
	}

	res := c.client.Set(context.Background(), fmt.Sprintf("cart:%d", user_id), string(json_cart), ttl)

	if err := res.Err(); err != nil {
		return IntErrorCode, nil
	}

	return product_id, nil
}

func (c *RedisController) AddToCart(user_id uint64, product_id uint64, quantity uint) (uint64, error) {
	record := fmt.Sprintf("cart:%d", user_id)

	tu := c.client.Get(context.Background(), record)

	err := tu.Err()

	if err == redis.Nil {
		return c.CreateCartAndPut(user_id, product_id, quantity)
	}

	if err != nil {
		return IntErrorCode, err
	}

	var cart []models.TempCartItem

	err = json.Unmarshal([]byte(tu.Val()), &cart)

	if err != nil {
		return IntErrorCode, err
	}

	ttl, err := c.getUserIDTTL();

	if err != nil {
		return IntErrorCode, err
	}

	if !c.checkCartForCollisionsAndAppend(cart, product_id, quantity) {
		cart = append(cart, *models.NewTempCartItem(product_id, quantity))
	}

	newcart, err := json.Marshal(cart)

	if err != nil {
		return IntErrorCode, err
	}

	err = c.client.Set(context.Background(), record, newcart, ttl).Err()

	if err != nil {
		return IntErrorCode, err
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

func (c *RedisController) getUserIDTTL() (time.Duration, error) {
	res := c.client.TTL(context.Background(), UserIDKey)

	if err := res.Err(); err != nil {
		return time.Duration(0), err
	}

	return res.Val(), nil
}
