package server

import (
	authClient "app/services/auth/client"
	"app/services/gate/dto"
	"context"
	"fmt"
	"log"
	"time"
)

const cacheTTL = time.Minute * 5

func (s *pbServer) authInfo(ctx context.Context) (*authClient.UserInfo, error) {
	info, ok := ctx.Value(userInfoKey).(*authClient.UserInfo)
	if !ok {
		return nil, fmt.Errorf("not authorized")
	}

	return info, nil
}

func (s *pbServer) authInfoRaw(ctx context.Context, token string) (*authClient.UserInfo, error) {
	redisStart := time.Now()

	redisData, err := s.redis.Get(token)

	redisFinish := time.Now()
	registerCacheHandle("redis", redisFinish.Sub(redisStart))

	if err != nil {
		// Ошибка отсутствия значения также логируется для отладки
		log.Printf("%s error from redis: %s\n", token, err.Error())
	} else {
		return &authClient.UserInfo{
			ID:    redisData.ID,
			Token: token,
		}, nil
	}

	authStart := time.Now()

	info, err := s.auth.Info(ctx, token)

	authFinish := time.Now()
	registerCacheHandle("auth", authFinish.Sub(authStart))

	if err != nil {
		return nil, err
	}

	// В данном случае кешер сеттится специально здесь, а не в сервисе авторизации
	err = s.redis.Set(token, dto.UserInfo{ID: info.ID}, cacheTTL)
	if err != nil {
		log.Println(err)
	}

	return info, nil
}
