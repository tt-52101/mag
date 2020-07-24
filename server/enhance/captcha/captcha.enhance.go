// @Title: captcha.enhance.go
// @Author: key7men@gmail.com
// @Description: 扩展dchest/captcha包，允许将captcha存入redis
// @Update: 2020/7/23 9:15 AM
package captcha

import (
	"encoding/hex"
	"time"

	"github.com/dchest/captcha"
	"github.com/go-redis/redis"
)

// Logger 任意实现Logger
type Logger interface {
	Printf(format string, args ...interface{})
}

func NewRedisStore(opts *redis.Options, expiration time.Duration, out Logger, prefix ...string) captcha.Store {
	if opts == nil {
		panic("options cannot be nil")
	}
	return NewRedisStoreWithCli(
			redis.NewClient(opts),
			expiration,
			out,
			prefix...,
		)
}

func NewRedisStoreWithCli(cli *redis.Client, expiration time.Duration, out Logger, prefix ...string) captcha.Store {
	store := &redisStore{
		cli: 		cli,
		expiration: expiration,
		out: 		out,
	}
	if len(prefix) > 0 {
		store.prefix = prefix[0]
	}
	return store
}

type clienter interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(keys ...string) *redis.IntCmd
}

type redisStore struct {
	cli 		clienter
	prefix 		string
	out			Logger
	expiration	time.Duration
}

func (r *redisStore) getKey(id string) string {
	return r.prefix + id
}

func (r *redisStore) printf(format string, args ...interface{}) {
	if r.out != nil {
		r.out.Printf(format, args...)
	}
}

func (r *redisStore) Set(id string, digits []byte) {
	cmd := r.cli.Set(r.getKey(id), hex.EncodeToString(digits), r.expiration)
	if err := cmd.Err(); err != nil {
		r.printf("redis execution set command error: %s", err.Error())
	}
	return
}

func (r *redisStore) Get(id string, clear bool) []byte {
	key := r.getKey(id)
	cmd := r.cli.Get(key)
	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			return nil
		}
		r.printf("redis execution get command error: %s", err.Error())
		return nil
	}

	b, err := hex.DecodeString(cmd.Val())
	if err != nil {
		r.printf("hex decoding error: %s", err.Error())
		return nil
	}

	if clear {
		cmd := r.cli.Del(key)
		if err := cmd.Err(); err != nil {
			r.printf("redis execution del command error: %s", err.Error())
			return nil
		}
	}

	return b
}

