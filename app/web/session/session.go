package session

import (
	"github.com/gorilla/sessions"
	. "github.com/netbrain/cloudfiler/app/conf"
	"net/http"
)

const (
	DefaultSessionName = "cloudfiler"
)

type Session struct {
	w            http.ResponseWriter
	r            *http.Request
	sessionStore *sessions.FilesystemStore
	session      *sessions.Session
}

func NewSession(w http.ResponseWriter, r *http.Request, name ...string) *Session {
	var sessionName string

	if len(name) == 0 {
		sessionName = DefaultSessionName
	} else {
		sessionName = name[0]
	}

	s := &Session{
		w: w,
		r: r,
		sessionStore: sessions.NewFilesystemStore(
			"",
			Config.SessionStoreAuthenticationKey,
			Config.SessionStoreEncryptionKey,
		),
	}

	s.sessionStore.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 30 * 60,
	}

	s.session, _ = s.sessionStore.Get(r, sessionName)

	return s
}

func (s *Session) Get(key interface{}) (v interface{}, ok bool) {
	v, ok = s.session.Values[key]
	return
}

func (s *Session) Set(key, val interface{}) {
	s.session.Values[key] = val
	s.save()
}

func (s *Session) Remove(key interface{}) {
	delete(s.session.Values, key)
	s.save()
}

func (s *Session) AddFlash(v interface{}) {
	s.session.AddFlash(v)
	s.save()
}

func (s *Session) Flash() (flashes []interface{}) {
	flashes = s.session.Flashes()
	s.save()
	return
}

func (s *Session) Destroy(w http.ResponseWriter, r *http.Request) {
	s.session.Options.MaxAge = -1
	s.save()
}

func (s *Session) save() {
	if err := s.session.Save(s.r, s.w); err != nil {
		panic(err)
	}
}
