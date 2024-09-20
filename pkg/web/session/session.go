package session

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/J4yTr1n1ty/meal-planner/pkg/config"
	"github.com/J4yTr1n1ty/meal-planner/pkg/redis"
)

const (
	ContextKey = "session"
)

// Session cookie object that is kept server side in Redis.
type Session struct {
	UUID     string    `json:"-"` // not stored
	LoggedIn bool      `json:"loggedIn"`
	LastSeen time.Time `json:"lastSeen"`
}

func New() *Session {
	return &Session{
		UUID: uuid.New().String(),
	}
}

func LoadOrNew(r *http.Request) *Session {
	session := New()

	sessionId, err := r.Cookie(config.SessionKey)
	if err != nil {
		log.Println("session.LoadOrNew: cookie error, new session: ", err)
		return session
	}

	session.UUID = sessionId.Value
	key := fmt.Sprintf(config.SessionRedisKeyFormat, session.UUID)

	err = redis.Get(key, session)
	if err != nil {
		log.Println(fmt.Sprintf("session.LoadOrNew: didn find session %s in Redis: %s", key, err))
	}

	return session
}
