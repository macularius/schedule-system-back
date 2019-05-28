package app

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"myapp/app/models/entities"
	"strconv"
	"sync"
	"time"
)

// Cache структура кэша для хранения сессий пользователей
type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	sessions          map[string]entities.Session
}

func createTokken(userID int, time string) string {
	str := fmt.Sprintf("%ssalt%d", time, userID)
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// New инициализация кэша
// defaultExpiration — время жизни кеша по-умолчанию,
// 		если установлено значение меньше или равно 0 — время жизни кеша бессрочно.
// cleanupInterval — интервал между удалением просроченного кеша.
//		При установленном значении меньше или равно 0 — очистка и удаление просроченного кеша не происходит.
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

	// инициализируем карту(map) в паре ключ(string)/значение(Session)
	sessions := make(map[string]entities.Session)

	cache := Cache{
		sessions:          sessions,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	// Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

// Add добавляет сессию в кэш
func (c *Cache) Add(login string, userID int, duration time.Duration) {

	var expiration int64
	var created = time.Now()

	// Если продолжительность жизни равна 0 - используется значение по-умолчанию
	if duration == 0 {
		duration = c.defaultExpiration
	}

	// Устанавливаем время истечения кеша
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	c.Lock()
	defer c.Unlock()

	connection, err := sql.Open("postgres", GetNewConnectionString(login))
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	key := createTokken(userID, strconv.FormatInt(created.Unix(), 10))

	if _, exist := c.sessions[key]; exist {
		c.sessions[key] = entities.Session{
			UserID:     userID,
			Connection: connection,
			Login:      login,
			Created:    created,
			Expiration: expiration,
		}
	} else {
		log.Fatal("This tokken is exist")
	}
}

// Get получение сессии из кеша
func (c *Cache) Get(key string) (interface{}, bool) {

	c.RLock()

	defer c.RUnlock()

	session, found := c.sessions[key]

	// ключ не найден
	if !found {
		return nil, false
	}

	// Проверка на установку времени истечения, в противном случае он бессрочный
	if session.Expiration > 0 {

		// Если в момент запроса кеш устарел возвращаем nil
		if time.Now().UnixNano() > session.Expiration {
			return nil, false
		}

	}

	return session.Connection, true
}

// Delete удаление из кэша сессиии с токеном
func (c *Cache) Delete(key string) error {

	c.Lock()

	defer c.Unlock()

	if _, found := c.sessions[key]; !found {
		return errors.New("Key not found")
	}

	delete(c.sessions, key)

	return nil
}

// StartGC запуск сборки мусора в горутине
func (c *Cache) StartGC() {
	go c.GC()
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
