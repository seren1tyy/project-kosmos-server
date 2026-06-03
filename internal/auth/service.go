package auth

import (
	"encoding/json"
	"log"
	"net"
	"project-kosmos-server/internal/db"
	"project-kosmos-server/internal/protocol"
	"project-kosmos-server/internal/session"
	"project-kosmos-server/pkg/crypto"
)

type Service struct {
	DB         db.DB
	Key        []byte
	SessionMgr *session.Manager
}

func (s *Service) Handle(conn net.Conn, req map[string]interface{}) {
	switch req["type"] {
	case "ping":
		protocol.SendPacket(conn, map[string]interface{}{"type": "pong", "status": "ok"})
	case "ready":
		protocol.SendPacket(conn, map[string]interface{}{"type": "server_status", "message": "online"})
	case "auth":
		s.handleAuth(conn, req)
	default:
		protocol.SendPacket(conn, map[string]interface{}{"type": "error", "message": "unknown"})
	}
}

func (s *Service) handleAuth(conn net.Conn, req map[string]interface{}) {
	b64, ok := req["data"].(string)
	if !ok {
		protocol.SendPacket(conn, map[string]interface{}{"type": "auth_result", "success": false, "message": "missing_data"})
		return
	}

	dec, err := crypto.DecryptXOR(b64, s.Key)
	if err != nil {
		log.Printf("Decrypt error: %v", err)
		protocol.SendPacket(conn, map[string]interface{}{"type": "auth_result", "success": false, "message": "decrypt_failed"})
		return
	}

	var creds struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal([]byte(dec), &creds); err != nil {
		log.Printf("Auth JSON parse error: %v", err)
		protocol.SendPacket(conn, map[string]interface{}{"type": "auth_result", "success": false, "message": "invalid_payload"})
		return
	}

	log.Printf("Auth attempt: login='%s'", creds.Login)
	u, ok, err := s.DB.FindByLogin(creds.Login)
	if err != nil {
		log.Printf("DB error: %v", err)
		protocol.SendPacket(conn, map[string]interface{}{"type": "auth_result", "success": false, "message": "db_error"})
		return
	}
	if !ok || u.Password != creds.Password {
		log.Printf("Auth failed: %s", creds.Login)
		protocol.SendPacket(conn, map[string]interface{}{"type": "auth_result", "success": false, "message": "invalid_credentials"})
		return
	}

	// 1. Проверяем наличие персонажей
	hasChars, err := s.DB.HasCharacters(u.ID)
	if err != nil {
		log.Printf("Char check error: %v (falling back to false)", err)
	}

	// 2. Создаём/обновляем сессию (автоматически кикнет старый клиент, если залогинен)
	token, err := s.SessionMgr.CreateSession(u.ID, conn)
	if err != nil {
		log.Printf("Session create error: %v", err)
		protocol.SendPacket(conn, map[string]interface{}{"type": "auth_result", "success": false, "message": "session_error"})
		return
	}

	log.Printf("Auth success: %s (ID:%d) | Has chars: %v", creds.Login, u.ID, hasChars)

	// 3. Отправляем ответ с токеном
	protocol.SendPacket(conn, map[string]interface{}{
		"type":          "auth_result",
		"success":       true,
		"user_id":       u.ID,
		"session_token": token,
		"has_character": hasChars,
	})
}
