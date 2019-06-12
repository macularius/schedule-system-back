package app

import (
	"database/sql"
	"fmt"
	"myapp/app/models/entities"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var Cache = new(cache)
var config DBConfig
var db *sql.DB

// Cache структура кэша для хранения сессий пользователей
type cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	sessions          map[string]entities.Session
}

func (c *cache) init() {
	c.cleanupInterval = 500 * time.Second
	c.defaultExpiration = 600 * time.Second
	c.sessions = make(map[string]entities.Session)
	c.StartGC()
}

// IsExistByLogin возвращает sid и true, если сессия для данного login'а существует в кэше
func IsExistByLogin(login string) (string, bool) {
	for sid, session := range Cache.sessions {
		if session.Login == login {
			return sid, true
		}
	}
	return "", false
}

// IsExistBySID возвращает true, если сессия с данным sid существует в кэше
func IsExistBySID(sid string) bool {
	// Срабатывает при первом запуске кэша
	if Cache.cleanupInterval == 0 || Cache.defaultExpiration == 0 || Cache.sessions == nil {
		Cache.init()
	}

	var exist bool
	_, exist = Cache.sessions[sid]
	return exist
}

// Add добавление сессии в кэш
func Add(sid string, login string, password string, token string) error {
	// Срабатывает при первом запуске кэша
	if Cache.cleanupInterval == 0 || Cache.defaultExpiration == 0 || Cache.sessions == nil {
		Cache.init()
	}

	// Если токкен существует вернуть его сессию с обновленным временем жизни
	if _, exist := IsExistByLogin(login); exist {
		ExtendSession(login)
		return nil
	}

	db, err := sql.Open("postgres", GetConnectionString())
	if err != nil {
		return fmt.Errorf("Error creating connection: %s", err.Error())
	}
	defer db.Close()
	row := db.QueryRow(fmt.Sprintf("SELECT u.uid, e.eid FROM users as u, employees as e WHERE u.login = '%s' AND e.eid=u.eid", login))
	if err != nil {
		// log.Fatal("Error creating role: ", err.Error())
		return fmt.Errorf("Error selected users with login: %s", err.Error())
	}

	var uid int64
	var eid int64
	err = row.Scan(&uid, &eid)
	if err != nil {
		return fmt.Errorf("Error scanning uid and eid: %s", err.Error())
	}

	var created = time.Now()
	var duration = Cache.defaultExpiration
	var expiration = time.Now().Add(duration).UnixNano()

	Cache.Lock()
	defer Cache.Unlock()

	constr, err := GetNewConnectionString(login, password)
	if err != nil {
		return fmt.Errorf("GetNewConnectionString %s", err.Error())
	}

	connection, err := sql.Open("postgres", constr)
	if err != nil {
		return fmt.Errorf("Open %s", err.Error())
	}

	fmt.Printf("\nAdd\n sid: %s\ntoken: %s\nuid: %d\neid: %d\nlogin: %s", sid, token, uid, eid, login)
	if _, exist := Cache.sessions[sid]; !exist {
		Cache.sessions[sid] = entities.Session{
			Token:      token,
			UserID:     uid,
			EmployeeID: eid,
			Connection: connection,
			Login:      login,
			Created:    created,
			Expiration: expiration,
		}
	} else {
		return fmt.Errorf("not exist %s", err.Error())
	}
	fmt.Print("\n", Cache.sessions[sid], "\n")

	return nil
}

// GetConnectionBySID получение connect'а по токену, возвращает nil, если токен не существует
func GetConnectionBySID(sid string) *sql.DB {
	if session, exist := Cache.sessions[sid]; exist {
		return session.Connection
	}
	return nil
}

// DeleteBySID удаление сессии по токену
func DeleteBySID(sid string) error {
	Cache.Lock()
	defer Cache.Unlock()

	if _, exist := Cache.sessions[sid]; !exist {
		return fmt.Errorf("Сессии не существует. SID: %s", sid)
	}

	Cache.sessions[sid].Connection.Close()
	delete(Cache.sessions, sid)
	fmt.Println("logout: ", Cache.sessions[sid])
	return nil
}

// ExtendSession обновление сессии пользователя
func ExtendSession(login string) {
	var token string
	var session entities.Session
	for token, session = range Cache.sessions {
		if session.Login == login {
			break
		}
	}
	newSession := entities.Session{
		UserID:     session.UserID,
		Connection: session.Connection,
		Login:      session.Login,
		Created:    session.Created,
		Expiration: time.Now().Add(Cache.defaultExpiration).UnixNano(),
	}

	Cache.Lock()
	defer Cache.Unlock()

	delete(Cache.sessions, token)
	Cache.sessions[token] = newSession

}

// GetSessionBySID возвращает сессию по sid, если она существует
func GetSessionBySID(sid string) (entities.Session, error) {
	if IsExistBySID(sid) {
		return Cache.sessions[sid], nil
	}

	return *new(entities.Session), fmt.Errorf("Запрошанная сессия не существует(sid: %s)", sid)
}

// StartGC запуск сборки мусора в горутине
func (c *cache) StartGC() {
	go Cache.GC()
}

// GC сборка мусора
func (c *cache) GC() {
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
func (c *cache) expiredKeys() (keys []string) {
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
func (c *cache) clearItems(keys []string) {
	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		delete(c.sessions, k)
	}
}
