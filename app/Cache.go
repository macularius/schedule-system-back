package app

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"log"
	"myapp/app/models/entities"
	"strconv"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var cache Cache
var config DBConfig
var db *sql.DB

// Cache структура кэша для хранения сессий пользователей
type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	sessions          map[string]entities.Session
}

func (c *Cache) init() {
	c.cleanupInterval = 500 * time.Second
	c.defaultExpiration = 600 * time.Second
	c.sessions = make(map[string]entities.Session)
	c.StartGC()
}

func createToken(userID int, time string) string {
	str := fmt.Sprintf("%ssalt%d", time, userID)
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// IsExistByToken возвращает true, если заданный token существует в кэше
func IsExistByToken(token string) bool {
	var exist bool
	_, exist = cache.sessions[token]
	return exist
}

// IsExistByLogin возвращает true, если сессия для данного login'а существует в кэше
func IsExistByLogin(login string) (string, bool) {
	for token, session := range cache.sessions {
		if session.Login == login {
			return token, true
		}
	}
	return "", false
}

// Add добавление сессии в кэш
func Add(login string) (map[string]entities.Session, string) {
	// Срабатывает при первом запуске кэша
	if cache.cleanupInterval == 0 || cache.defaultExpiration == 0 || cache.sessions == nil {
		cache.init()
	}

	// Если токкен существует вернуть его сессию с обновленным временем жизни
	if _, exist := IsExistByLogin(login); exist {
		ExtendSession(login)
		return cache.sessions, ""
	}

	db, err := sql.Open("postgres", GetConnectionString())
	if err != nil {
		// log.Fatal("Error creating connection: ", err.Error())
		return nil, fmt.Sprintf("Error creating connection: %s", err.Error())
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT uid FROM users WHERE login = '%s'", login))
	if err != nil {
		// log.Fatal("Error creating role: ", err.Error())
		return nil, fmt.Sprintf("Error selected users with login: %s", err.Error())
	}
	defer rows.Close()

	var uid int
	if rows.Next() {
		err = rows.Scan(&uid)
	}
	if err != nil {
		// log.Fatal("Error creating role: ", err.Error())
		return nil, fmt.Sprintf("Error scanning uid: %s", err.Error())
	}
	var created = time.Now()
	var duration = cache.defaultExpiration
	var expiration = time.Now().Add(duration).UnixNano()

	cache.Lock()
	defer cache.Unlock()

	token := createToken(uid, strconv.FormatInt(created.Unix(), 10))
	constr, err := GetNewConnectionString(login)
	if err != nil {
		return nil, fmt.Sprintf("1 %s", err.Error())
	}

	connection, err := sql.Open("postgres", constr)
	if err != nil {
		return nil, fmt.Sprintf("2 %s", err.Error())
	}

	if _, exist := cache.sessions[token]; !exist {
		cache.sessions[token] = entities.Session{
			UserID:     uid,
			Connection: connection,
			Login:      login,
			Created:    created,
			Expiration: expiration,
		}
	} else {
		// log.Fatal("This token is exist")
		fmt.Println(cache.sessions)
		return nil, fmt.Sprintf("3 %s", err.Error())
	}

	// session := cache.sessions[token]

	return cache.sessions, ""
}

// GetConnectionByToken получение connect'а по токену, возвращает nil, если токен не существует
func GetConnectionByToken(token string) *sql.DB {
	if session, exist := cache.sessions[token]; exist {
		return session.Connection
	}
	return nil
}

// DeleteByToken удаление сессии по токену
func DeleteByToken(token string) {
	cache.Lock()
	defer cache.Unlock()

	if _, exist := cache.sessions[token]; !exist {
		log.Fatal()
	}

	delete(cache.sessions, token)
}

// DeleteByLogin удаление сессии по токену
func DeleteByLogin(login string) {

	var token string
	for key, session := range cache.sessions {
		if session.Login == login {
			token = key
			break
		}
	}

	cache.Lock()
	defer cache.Unlock()

	delete(cache.sessions, token)
}

// ExtendSession обновление сессии пользователя
func ExtendSession(login string) {
	var token string
	var session entities.Session
	for token, session = range cache.sessions {
		if session.Login == login {
			break
		}
	}
	newSession := entities.Session{
		UserID:     session.UserID,
		Connection: session.Connection,
		Login:      session.Login,
		Created:    session.Created,
		Expiration: time.Now().Add(cache.defaultExpiration).UnixNano(),
	}

	cache.Lock()
	defer cache.Unlock()

	delete(cache.sessions, token)
	cache.sessions[token] = newSession
}

// StartGC запуск сборки мусора в горутине
func (c *Cache) StartGC() {
	go cache.GC()
}

// GC сборка мусора
func (c *Cache) GC() {
	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(c.cleanupInterval)

		if c.sessions == nil {
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

// expiredKeys возвращает список "просроченных" ключей
func (c *Cache) expiredKeys() (keys []string) {
	c.RLock()
	defer c.RUnlock()

	for k, i := range c.sessions {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearItems(keys []string) {
	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		delete(c.sessions, k)
	}
}
