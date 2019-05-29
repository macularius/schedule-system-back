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

func createTokken(userID int, time string) string {
	str := fmt.Sprintf("%ssalt%d", time, userID)
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// IsExistByTokken возвращает true, если заданный tokken существует в кэше
func IsExistByTokken(tokken string) bool {
	var exist bool
	_, exist = cache.sessions[tokken]
	return exist
}

// IsExistByLogin возвращает true, если сессия для данного login'а существует в кэше
func IsExistByLogin(login string) (string, bool) {
	for tokken, session := range cache.sessions {
		if session.Login == login {
			return tokken, true
		}
	}
	return "", false
}

// Add добавление сессии в кэш
func Add(login string) (string, map[string]entities.Session) {
	// Срабатывает при первом запуске кэша
	if cache.cleanupInterval == 0 || cache.defaultExpiration == 0 || cache.sessions == nil {
		cache.init()
	}

	// Если токкен существует вернуть его сессию с обновленным временем жизни
	if tokken, exist := IsExistByLogin(login); exist {
		ExtendSession(login)
		return tokken, cache.sessions
	}

	db, err := sql.Open("postgres", config.GetConnectionString())
	if err != nil {
		// log.Fatal("Error creating connection: ", err.Error())
		return "Error creating connection: ", nil
	}
	defer db.Close()
	rows, err := db.Query("SELECT uid FROM users WHERE login = '" + login + "'")
	if err != nil {
		// log.Fatal("Error creating role: ", err.Error())
		return "Error creating role: ", nil
	}
	defer rows.Close()

	var uid int
	err = rows.Scan(&uid)
	var created = time.Now()
	var duration = cache.defaultExpiration
	var expiration = time.Now().Add(duration).UnixNano()

	cache.Lock()
	defer cache.Unlock()

	tokken := createTokken(uid, strconv.FormatInt(created.Unix(), 10))
	constr, err := GetNewConnectionString(login)
	if err != nil {
		return err.Error(), nil
	}

	connection, err := sql.Open("postgres", constr)
	if err != nil {
		return "Error creating connection", nil
	}

	if _, exist := cache.sessions[tokken]; !exist {
		cache.sessions[tokken] = entities.Session{
			UserID:     uid,
			Connection: connection,
			Login:      login,
			Created:    created,
			Expiration: expiration,
		}
	} else {
		// log.Fatal("This tokken is exist")
		fmt.Println(cache.sessions)
		return "This tokken(" + tokken + ") is exist", nil
	}

	// session := cache.sessions[tokken]

	return tokken, cache.sessions
}

// GetConnectionByTokken получение connect'а по токену, возвращает nil, если токен не существует
func GetConnectionByTokken(tokken string) *sql.DB {
	if session, exist := cache.sessions[tokken]; exist {
		return session.Connection
	}
	return nil
}

// DeleteByTokken удаление сессии по токену
func DeleteByTokken(tokken string) {
	cache.Lock()
	defer cache.Unlock()

	if _, exist := cache.sessions[tokken]; !exist {
		log.Fatal()
	}

	delete(cache.sessions, tokken)
}

// DeleteByLogin удаление сессии по токену
func DeleteByLogin(login string) {

	var tokken string
	for key, session := range cache.sessions {
		if session.Login == login {
			tokken = key
			break
		}
	}

	cache.Lock()
	defer cache.Unlock()

	delete(cache.sessions, tokken)
}

// ExtendSession обновление сессии пользователя
func ExtendSession(login string) {
	var tokken string
	var session entities.Session
	for tokken, session = range cache.sessions {
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

	delete(cache.sessions, tokken)
	cache.sessions[tokken] = newSession
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
