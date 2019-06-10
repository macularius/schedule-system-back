package app

import (
	"database/sql"
	"fmt"
	"myapp/app/models/entities"
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

// IsExistByLogin возвращает sid и true, если сессия для данного login'а существует в кэше
func IsExistByLogin(login string) (string, bool) {
	for sid, session := range cache.sessions {
		if session.Login == login {
			return sid, true
		}
	}
	return "", false
}

// IsExistBySID возвращает true, если сессия с данным sid существует в кэше
func IsExistBySID(sid string) bool {
	var exist bool
	_, exist = cache.sessions[sid]
	return exist
}

// Add добавление сессии в кэш
func Add(sid string, login string, token string) (map[string]entities.Session, string) {
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
	row := db.QueryRow(fmt.Sprintf("SELECT u.uid, e.eid FROM users as u, employees as e WHERE u.login = '%s' AND e.eid=u.eid", login))
	if err != nil {
		// log.Fatal("Error creating role: ", err.Error())
		return nil, fmt.Sprintf("Error selected users with login: %s", err.Error())
	}

	var uid int
	var eid int
	err = row.Scan(&uid, &eid)
	if err != nil {
		return nil, fmt.Sprintf("Error scanning uid and eid: %s", err.Error())
	}

	var created = time.Now()
	var duration = cache.defaultExpiration
	var expiration = time.Now().Add(duration).UnixNano()

	cache.Lock()
	defer cache.Unlock()

	constr, err := GetNewConnectionString(login)
	if err != nil {
		return nil, fmt.Sprintf("GetNewConnectionString %s", err.Error())
	}

	connection, err := sql.Open("postgres", constr)
	if err != nil {
		return nil, fmt.Sprintf("Open %s", err.Error())
	}

	if _, exist := cache.sessions[sid]; !exist {
		cache.sessions[sid] = entities.Session{
			Token:      token,
			UserID:     uid,
			EmployeeID: eid,
			Connection: connection,
			Login:      login,
			Created:    created,
			Expiration: expiration,
		}
	} else {
		return nil, fmt.Sprintf("not exist %s", err.Error())
	}

	// session := cache.sessions[token]

	return cache.sessions, ""
}

// GetConnectionBySID получение connect'а по токену, возвращает nil, если токен не существует
func GetConnectionBySID(sid string) *sql.DB {
	if session, exist := cache.sessions[sid]; exist {
		return session.Connection
	}
	return nil
}

// DeleteBySID удаление сессии по токену
func DeleteBySID(sid string) error {
	cache.Lock()
	defer cache.Unlock()

	if _, exist := cache.sessions[sid]; !exist {
		return fmt.Errorf("Сессии не существует. SID: %s", sid)
	}

	delete(cache.sessions, sid)
	fmt.Println("logout: ", cache.sessions[sid])
	return nil
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

// GetSessionBySID возвращает сессию по токену, если она существует
func GetSessionBySID(sid string) (entities.Session, error) {
	if IsExistBySID(sid) {
		return cache.sessions[sid], nil
	}

	return *new(entities.Session), fmt.Errorf("Запрошанная сессия не существует(sid: %s)", sid)
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
