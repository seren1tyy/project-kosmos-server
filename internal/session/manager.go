package session

import (
	"crypto/rand"
	"encoding/base64"
	"net"
	"sync"

	"project-kosmos-server/internal/protocol"
)

type Session struct {
	UserID          int
	Token           string
	Conn            net.Conn
	Active          bool
	CharacterID     int
	CurrentSystemID int
	ShipTypeID      int
}

type Manager struct {
	mu      sync.RWMutex
	byUser  map[int]*Session
	byToken map[string]*Session
}

func NewManager() *Manager {
	return &Manager{
		byUser:  make(map[int]*Session),
		byToken: make(map[string]*Session),
	}
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (m *Manager) CreateSession(userID int, conn net.Conn) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// ✅ КИКАЕМ СТАРУЮ СЕССИЮ С ПРЕДУПРЕЖДЕНИЕМ
	if old, exists := m.byUser[userID]; exists && old.Active {
		old.Active = false
		protocol.SendPacket(old.Conn, map[string]interface{}{
			"type":    "session_kicked",
			"message": "account_logged_in_elsewhere",
		})
		old.Conn.Close()
		delete(m.byToken, old.Token)
	}

	token, err := GenerateToken()
	if err != nil {
		return "", err
	}

	sess := &Session{
		UserID: userID,
		Token:  token,
		Conn:   conn,
		Active: true,
	}
	m.byUser[userID] = sess
	m.byToken[token] = sess
	return token, nil
}

func (m *Manager) GetByToken(token string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sess, ok := m.byToken[token]
	return sess, ok && sess.Active
}

func (m *Manager) Invalidate(token string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sess, ok := m.byToken[token]; ok {
		sess.Active = false
		if sess.Conn != nil {
			sess.Conn.Close()
		}
		delete(m.byToken, token)
		delete(m.byUser, sess.UserID)
	}
}

func (m *Manager) UpdateCharacter(token string, charID, systemID, shipTypeID int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if sess, ok := m.byToken[token]; ok && sess.Active {
		sess.CharacterID = charID
		sess.CurrentSystemID = systemID
		sess.ShipTypeID = shipTypeID
		return true
	}
	return false
}

func (m *Manager) GetSession(token string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sess, ok := m.byToken[token]
	return sess, ok && sess.Active
}
