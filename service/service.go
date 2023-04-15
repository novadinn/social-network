package service

import (
	"log"
	
	"github.com/novadinn/social-network/queries"

	"github.com/rs/xid"
	"github.com/gorilla/sessions"
)

type Service struct {
	Queries *queries.Queries
	Logger *log.Logger
	Store *sessions.CookieStore
	Payload map[string]interface{}
}

func New() *Service {
	que := queries.New()
	// que.Clear()
	que = queries.New()
	
	logger := log.Default()
	
	store := sessions.NewCookieStore([]byte("sn-authentication-key"))
	store.Options.HttpOnly = true
	store.Options.Secure = true

	payload := make(map[string]interface{})
	
		return &Service{Queries: que,
			Logger: logger,
			Store: store,
			Payload: payload,
		}
}

func (svc *Service) Push(key string, value interface{}) {
	svc.Payload[key] = value
}

func (svc *Service) Pop(key string) (interface{}, bool) {
	value, ok := svc.Payload[key]
	if !ok {
		return nil, false
	}
	
	delete(svc.Payload, key)
	return value, true
}

func genID() string {
	return xid.New().String()
}
