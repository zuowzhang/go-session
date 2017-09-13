package session

import (
	"sync"
	"errors"
	"io"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"
	"log"
)

type Provider interface {
	Read(sid string) (Session, error)
	Remove(sid string) error
	Gc(maxLifeTime int64)
}

type Session interface {
	Set(interface{}, interface{}) error
	Get(interface{}) interface{}
	Delete(interface{}) error
	Id() string
}

type SessionMgr struct {
	cookieName  string
	lock        sync.Mutex
	maxLifeTime int64
	provider    Provider
}

const PROVIDER_MEMORY string = "memory"
const PROVIDER_REDIS string = "redis"

var providers = make(map[string]Provider)

var Logger = logProxy{}

func InitLog(log Log) {
	Logger.impl = log
}

func RegisterProvider(providerName string, provider Provider) {
	if provider == nil {
		panic("provider is nil")
	}
	if _, ok := providers[providerName]; ok {
		panic("duplicated registe session provider")
	}
	providers[providerName] = provider
	log.Printf("RegisterProvider %s, providers.len() = %d--%#v--%x\n", providerName, len(providers), providers, providers)
}

func NewSessionMgr(providerName, cookieName string, maxLifeTime int64) (*SessionMgr, error) {
	log.Printf("lookup provider %s, providers.len() = %d--%#v--%x\n", providerName, len(providers), providers, providers)
	provider, ok := providers[providerName]
	if !ok {
		Logger.E("provider %s not exists\n", providerName)
		return nil, errors.New("unknown session provider")
	}
	mgr := &SessionMgr{
		cookieName:cookieName,
		maxLifeTime:maxLifeTime,
		provider:provider,
	}
	go mgr.Gc()
	return mgr, nil
}

func sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (mgr *SessionMgr)SessionStart(w http.ResponseWriter, r *http.Request) Session {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	cookie, err := r.Cookie(mgr.cookieName)
	if err != nil || cookie.Value == "" {
		sessionId := sessionId()
		if sessionId != "" {
			session, err := mgr.provider.Read(sessionId)
			if err == nil {
				cookie := http.Cookie{
					Name:mgr.cookieName,
					Value:sessionId,
					Path:"/",
					HttpOnly:true,
					MaxAge:int(mgr.maxLifeTime),
				}
				http.SetCookie(w, &cookie)
				return session
			}
		}
	} else {
		sessionId := url.PathEscape(cookie.Value)
		session, err := mgr.provider.Read(sessionId)
		if err == nil {
			return session
		}
	}
	return nil
}

func (mgr *SessionMgr)SessionStop(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(mgr.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	mgr.provider.Remove(cookie.Value)
	cookie = &http.Cookie{
		Name:mgr.cookieName,
		Path:"/",
		HttpOnly:true,
		Expires:time.Now(),
		MaxAge:-1,
	}
	http.SetCookie(w, cookie)
}

func (mgr *SessionMgr)Gc() {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	mgr.provider.Gc(mgr.maxLifeTime)
	time.AfterFunc(time.Duration(mgr.maxLifeTime), func() {
		mgr.Gc()
	})
}