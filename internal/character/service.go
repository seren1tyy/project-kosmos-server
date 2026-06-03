package character

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"sync"

	"project-kosmos-server/internal/items"
	"project-kosmos-server/internal/protocol"
	"project-kosmos-server/internal/universe"
	"project-kosmos-server/pkg/crypto"
)

type Service struct {
	DB          *sql.DB
	Key         []byte
	mu          sync.Mutex
	itemsSvc    *items.Service
	universeSvc *universe.Service
}

func NewService(db *sql.DB, key []byte, itemsSvc *items.Service, universeSvc *universe.Service) *Service {
	return &Service{DB: db, Key: key, itemsSvc: itemsSvc, universeSvc: universeSvc}
}

func (s *Service) HandleCreate(conn net.Conn, accountID int, req map[string]interface{}) {
	log.Printf("[CharCreate] ENTERED. AccountID: %d | ReqType: %v", accountID, req["type"])

	if len(s.Key) == 0 {
		log.Printf("[CharCreate] Crypto key is EMPTY! Check config.json & main.go")
		protocol.SendPacket(conn, map[string]interface{}{"type": "create_result", "success": false, "message": "server_config_error"})
		return
	}

	b64, ok := req["data"].(string)
	if !ok {
		protocol.SendPacket(conn, map[string]interface{}{"type": "create_result", "success": false, "message": "invalid_payload"})
		return
	}

	dec, err := crypto.DecryptXOR(b64, s.Key)
	if err != nil {
		log.Printf("[CharCreate] Decrypt failed: %v", err)
		protocol.SendPacket(conn, map[string]interface{}{"type": "create_result", "success": false, "message": "decrypt_failed"})
		return
	}

	var data struct {
		Name      string `json:"name"`
		Race      string `json:"race"`
		Bloodline string `json:"bloodline"`
	}
	if err := json.Unmarshal([]byte(dec), &data); err != nil {
		protocol.SendPacket(conn, map[string]interface{}{"type": "create_result", "success": false, "message": "invalid_data_format"})
		return
	}

	if len(data.Name) < 3 || len(data.Name) > 30 {
		protocol.SendPacket(conn, map[string]interface{}{"type": "create_result", "success": false, "message": "invalid_name_length"})
		return
	}

	var exists int
	if err := s.DB.QueryRow("SELECT COUNT(*) FROM characters WHERE name = ?", data.Name).Scan(&exists); err != nil {
		log.Printf("[CharCreate] DB check failed: %v", err)
		protocol.SendPacket(conn, map[string]interface{}{"type": "create_result", "success": false, "message": "db_error"})
		return
	}
	if exists > 0 {
		protocol.SendPacket(conn, map[string]interface{}{"type": "create_result", "success": false, "message": "name_taken"})
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var maxID sql.NullInt64
	if err := s.DB.QueryRow("SELECT MAX(id) FROM characters").Scan(&maxID); err != nil {
		log.Printf("[CharCreate] MAX(id) query failed: %v", err)
		protocol.SendPacket(conn, map[string]interface{}{"type": "create_result", "success": false, "message": "db_error"})
		return
	}

	var newID int64 = 10_000_000
	if maxID.Valid {
		newID = maxID.Int64 + 1
	}
	if newID > 19_999_999 {
		protocol.SendPacket(conn, map[string]interface{}{"type": "create_result", "success": false, "message": "id_range_exhausted"})
		return
	}

	_, err = s.DB.Exec(
		"INSERT INTO characters (id, account_id, name, race, bloodline) VALUES (?, ?, ?, ?, ?)",
		newID, accountID, data.Name, data.Race, data.Bloodline,
	)
	if err != nil {
		log.Printf("[CharCreate] INSERT failed: %v", err)
		protocol.SendPacket(conn, map[string]interface{}{"type": "create_result", "success": false, "message": "creation_failed"})
		return
	}

	log.Printf("[CharCreate] SUCCESS: %s (ID:%d, Account:%d)", data.Name, newID, accountID)
	protocol.SendPacket(conn, map[string]interface{}{
		"type":    "create_result",
		"success": true,
		"message": "Character created successfully",
		"char_id": newID,
	})

	_, err = s.DB.Exec(
		"INSERT INTO character_assets (character_id, type_id, quantity, location) VALUES (?, 582, 1, 'dock')",
		newID,
	)
	if err != nil {
		log.Printf("[CharCreate] Failed to grant starter Wisp: %v", err)
		// Не критично — персонаж создан, но без корабля
	} else {
		log.Printf("[CharCreate] Granted starter Wisp to character %d", newID)
	}
}

func (s *Service) HandleFetch(conn net.Conn, accountID int, req map[string]interface{}) {
	page := 0
	perPage := 8

	// Диагностика: что пришло от клиента
	log.Printf("[CharFetch] RAW req keys: %v", mapKeys(req))
	if b64, ok := req["data"].(string); ok {
		log.Printf("[CharFetch] data field present, len=%d, first50=%s", len(b64), truncate(b64, 50))

		dec, err := crypto.DecryptXOR(b64, s.Key)
		if err != nil {
			log.Printf("[CharFetch] DecryptXOR failed: %v", err)
		} else {
			log.Printf("[CharFetch] Decrypted JSON: %s", string(dec))

			var data map[string]interface{}
			if err := json.Unmarshal([]byte(dec), &data); err != nil {
				log.Printf("[CharFetch] JSON parse failed: %v", err)
			} else {
				log.Printf("🔎 [CharFetch] Parsed data: %v", data)
				if p, ok := data["page"].(float64); ok {
					page = int(p)
					log.Printf("[CharFetch] page extracted: %d", page)
				} else {
					log.Printf("[CharFetch] page not found or wrong type")
				}
				if pp, ok := data["per_page"].(float64); ok && pp > 0 {
					perPage = int(pp)
					log.Printf("[CharFetch] per_page extracted: %d", perPage)
				}
			}
		}
	} else {
		log.Printf("[CharFetch] 'data' field missing or not string")
	}

	if page < 0 {
		page = 0
	}
	if perPage < 1 || perPage > 50 {
		perPage = 8
	}

	log.Printf("[CharFetch] FINAL: account=%d, page=%d, per_page=%d", accountID, page, perPage)

	var total int
	if err := s.DB.QueryRow("SELECT COUNT(*) FROM characters WHERE account_id = ?", accountID).Scan(&total); err != nil {
		log.Printf("[CharFetch] COUNT error: %v", err)
		protocol.SendPacket(conn, map[string]interface{}{"type": "fetch_chars_result", "success": false, "characters": []map[string]interface{}{}, "total": 0, "page": 0})
		return
	}

	offset := page * perPage
	rows, err := s.DB.Query(
		"SELECT id, name, race, bloodline FROM characters WHERE account_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?",
		accountID, perPage, offset,
	)
	if err != nil {
		log.Printf("[CharFetch] SELECT error: %v", err)
		protocol.SendPacket(conn, map[string]interface{}{"type": "fetch_chars_result", "success": false, "characters": []map[string]interface{}{}, "total": total, "page": page})
		return
	}
	defer rows.Close()

	chars := []map[string]interface{}{}
	for rows.Next() {
		var id int64
		var name, race, bloodline string
		rows.Scan(&id, &name, &race, &bloodline)
		chars = append(chars, map[string]interface{}{
			"id": id, "name": name, "race": race, "bloodline": bloodline,
		})
	}

	log.Printf("[CharFetch] Returned %d chars (total=%d)", len(chars), total)
	protocol.SendPacket(conn, map[string]interface{}{
		"type":       "fetch_chars_result",
		"success":    true,
		"characters": chars,
		"total":      total,
		"page":       page,
		"per_page":   perPage,
	})
}

// Вспомогательные функции (добавь ниже в том же файле)
func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
func (s *Service) HandleEnter(conn net.Conn, req map[string]interface{}) {
	b64, ok := req["data"].(string)
	if !ok {
		protocol.SendPacket(conn, map[string]interface{}{"type": "enter_game_result", "success": false, "message": "invalid_payload"})
		return
	}

	dec, err := crypto.DecryptXOR(b64, s.Key)
	if err != nil {
		protocol.SendPacket(conn, map[string]interface{}{"type": "enter_game_result", "success": false, "message": "decrypt_failed"})
		return
	}

	var payload struct {
		CharID int64 `json:"char_id"`
	}
	if err := json.Unmarshal([]byte(dec), &payload); err != nil {
		protocol.SendPacket(conn, map[string]interface{}{"type": "enter_game_result", "success": false, "message": "invalid_data"})
		return
	}

	// TODO: добавить проверку, что char_id принадлежит accountID сессии
	log.Printf("Entering game with char ID: %d", payload.CharID)
	protocol.SendPacket(conn, map[string]interface{}{
		"type":    "enter_game_result",
		"success": true,
		"char_id": payload.CharID,
	})

	var charName string
	s.DB.QueryRow("SELECT name FROM characters WHERE id = ?", payload.CharID).Scan(&charName)

	protocol.SendPacket(conn, map[string]interface{}{
		"type":      "enter_game_result",
		"success":   true,
		"char_id":   payload.CharID,
		"char_name": charName,
	})

}

func (s *Service) HandleCheckName(conn net.Conn, req map[string]interface{}) {
	b64, ok := req["data"].(string)
	if !ok {
		protocol.SendPacket(conn, map[string]interface{}{"type": "check_name_result", "is_taken": false, "error": "invalid_payload"})
		return
	}
	dec, err := crypto.DecryptXOR(b64, s.Key)
	if err != nil {
		protocol.SendPacket(conn, map[string]interface{}{"type": "check_name_result", "is_taken": false, "error": "decrypt_failed"})
		return
	}

	var payload struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(dec), &payload); err != nil {
		protocol.SendPacket(conn, map[string]interface{}{"type": "check_name_result", "is_taken": false, "error": "invalid_data"})
		return
	}

	var count int
	if err := s.DB.QueryRow("SELECT COUNT(*) FROM characters WHERE name = ?", payload.Name).Scan(&count); err != nil {
		protocol.SendPacket(conn, map[string]interface{}{"type": "check_name_result", "is_taken": false, "error": "db_error"})
		return
	}

	log.Printf("Name check: '%s' -> taken=%v", payload.Name, count > 0)
	protocol.SendPacket(conn, map[string]interface{}{
		"type":     "check_name_result",
		"is_taken": count > 0,
	})
}

func (s *Service) HandleGetInventory(conn net.Conn, characterID int) {
	assets, err := s.itemsSvc.GetCharacterAssets(characterID)
	if err != nil {
		log.Printf("[Inventory] DB error: %v", err)
		protocol.SendPacket(conn, map[string]interface{}{
			"type":    "inventory_result",
			"success": false,
			"assets":  []map[string]interface{}{},
		})
		return
	}

	// Преобразуем в массив для JSON
	var out []map[string]interface{}
	for _, a := range assets {
		out = append(out, map[string]interface{}{
			"instance_id": a.InstanceID,
			"type_id":     a.TypeID,
			"quantity":    a.Quantity,
			"location":    a.Location,
			"name":        a.Item.Name,
			"volume":      a.Item.Volume,
			"description": a.Item.Description,
		})
	}

	log.Printf("[Inventory] char=%d, assets=%d", characterID, len(out))
	protocol.SendPacket(conn, map[string]interface{}{
		"type":    "inventory_result",
		"success": true,
		"assets":  out,
	})
}

func (s *Service) HandleGetLocation(conn net.Conn, characterID int) {
	var stationID int
	err := s.DB.QueryRow("SELECT current_station_id FROM characters WHERE id = ?", characterID).Scan(&stationID)
	if err != nil {
		log.Printf("[Location] DB error for char %d: %v", characterID, err)
		protocol.SendPacket(conn, map[string]interface{}{"type": "location_result", "success": false})
		return
	}

	loc, err := s.universeSvc.GetLocationByStation(stationID)
	if err != nil {
		log.Printf("[Location] Universe error: %v", err)
		protocol.SendPacket(conn, map[string]interface{}{"type": "location_result", "success": false})
		return
	}

	log.Printf("[Location] Char %d is at: %s, %s", characterID, loc.StationName, loc.SystemName)

	protocol.SendPacket(conn, map[string]interface{}{
		"type":    "location_result",
		"success": true,
		"location": map[string]interface{}{
			"station":       loc.StationName,
			"system":        loc.SystemName,
			"constellation": loc.ConstellationName,
			"region":        loc.RegionName,
		},
	})
}
