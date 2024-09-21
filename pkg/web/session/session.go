package session

import (
	"context"
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

func FromContext(ctx context.Context) *Session {
	return ctx.Value(ContextKey).(*Session)
}

func LoadOrNew(r *http.Request) *Session {
	session := New()

	sessionId, err := r.Cookie(config.SessionCookieName)
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

// Save the session and send a cookie header.
func (s *Session) Save(w http.ResponseWriter) {
	// Roll a UUID session_id value.
	if s.UUID == "" {
		s.UUID = uuid.New().String()
	}

	// Ensure it is a valid UUID.
	if _, err := uuid.Parse(s.UUID); err != nil {
		log.Println(fmt.Sprintf("Error: Session.Save: got an invalid UUID session_id: %s", err))
		s.UUID = uuid.New().String()
	}

	// Ping last seen.
	s.LastSeen = time.Now()

	// Save their session object in Redis.
	key := fmt.Sprintf(config.SessionRedisKeyFormat, s.UUID)
	if err := redis.Set(key, s, config.SessionCookieMaxAge*time.Second); err != nil {
		log.Printf(fmt.Sprintf("Session.Save: couldn't write to Redis: %s", err))
	}

	cookie := &http.Cookie{
		Name:     config.SessionCookieName,
		Value:    s.UUID,
		MaxAge:   config.SessionCookieMaxAge,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

// Get the session from the current HTTP request context.
func Get(r *http.Request) *Session {
	if r == nil {
		panic("session.Get: http.Request is required")
	}

	ctx := r.Context()
	if sess, ok := ctx.Value(ContextKey).(*Session); ok {
		return sess
	}

	// If the session isn't on the request, it means I broke something.
	log.Println("session.Get(): didn't find session in request context!")
	return nil
}

// LoginUser marks a session as logged in to an account.
func LoginUser(w http.ResponseWriter, r *http.Request) error {
	sess := Get(r)
	sess.LoggedIn = true
	sess.Save(w)

	return nil
}
