package network

import (
	"encoding/json"
	"io"
	"log"
	"net"

	"project-kosmos-server/internal/auth"
	"project-kosmos-server/internal/character" // 👈 Новый импорт
	"project-kosmos-server/internal/protocol"
)

// Добавь передачу character.Service в HandleClient (обновим main.go ниже)
func HandleClient(conn net.Conn, authSvc *auth.Service, charSvc *character.Service) {
	defer conn.Close()
	if tcp, ok := conn.(*net.TCPConn); ok {
		tcp.SetNoDelay(true)
	}
	log.Printf("🔌 Connected: %s", conn.RemoteAddr())

	for {
		msg, err := protocol.ReadPacket(conn)
		if err != nil {
			if err != io.EOF {
				log.Printf("📉 Read: %v", err)
			}
			return
		}

		var req map[string]interface{}
		if err := json.Unmarshal(msg, &req); err != nil {
			continue
		}

		switch req["type"] {
		case "ping", "ready", "auth":
			authSvc.Handle(conn, req)

		case "check_name":
			charSvc.HandleCheckName(conn, req)

		case "create_character":
			tokenStr, _ := req["token"].(string)
			log.Printf("🔀 Routing to Create: token_prefix=%s", tokenStr[:min(len(tokenStr), 8)]+"...")

			if authSvc.SessionMgr == nil {
				protocol.SendPacket(conn, map[string]interface{}{"type": "create_result", "success": false, "message": "session_mgr_nil"})
				break
			}
			sess, ok := authSvc.SessionMgr.GetByToken(tokenStr)
			if !ok || sess == nil {
				log.Printf("❌ Create failed: invalid session (token=%s)", tokenStr)
				protocol.SendPacket(conn, map[string]interface{}{"type": "create_result", "success": false, "message": "invalid_session"})
				break
			}

			log.Printf("✅ Session valid (AccountID: %d). Calling HandleCreate.", sess.UserID)
			charSvc.HandleCreate(conn, sess.UserID, req)

		case "fetch_characters":
			tokenStr, _ := req["token"].(string)
			sess, ok := authSvc.SessionMgr.GetByToken(tokenStr)
			if !ok || sess == nil {
				protocol.SendPacket(conn, map[string]interface{}{"type": "fetch_chars_result", "success": false, "characters": []map[string]interface{}{}, "total": 0, "page": 0})
				break
			}
			charSvc.HandleFetch(conn, sess.UserID, req)

		case "enter_game":
			tokenStr, _ := req["token"].(string)
			sess, ok := authSvc.SessionMgr.GetByToken(tokenStr)
			if !ok || sess == nil {
				protocol.SendPacket(conn, map[string]interface{}{"type": "enter_game_result", "success": false, "message": "invalid_session"})
				break
			}
			charSvc.HandleEnter(conn, req)
		case "get_inventory":
			tokenStr, _ := req["token"].(string)
			sess, ok := authSvc.SessionMgr.GetByToken(tokenStr)
			if !ok || sess == nil {
				protocol.SendPacket(conn, map[string]interface{}{
					"type":    "inventory_result",
					"success": false,
					"assets":  []map[string]interface{}{},
				})
				break
			}
			// characterID берём из GameSession клиента, но для надёжности — из запроса
			charID, _ := req["char_id"].(float64)
			charSvc.HandleGetInventory(conn, int(charID))

		case "get_location":
			tokenStr, _ := req["token"].(string)
			sess, ok := authSvc.SessionMgr.GetByToken(tokenStr)
			if !ok || sess == nil {
				protocol.SendPacket(conn, map[string]interface{}{"type": "location_result", "success": false})
				break
			}
			charID, _ := req["char_id"].(float64)
			charSvc.HandleGetLocation(conn, int(charID))

		default:
			log.Printf("⚠️ Unknown packet type: %v", req["type"])
			protocol.SendPacket(conn, map[string]interface{}{"type": "error", "message": "unknown"})
		}
	}
}
