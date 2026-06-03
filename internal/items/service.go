package items

import (
	"database/sql"
	"fmt"
)

type ItemType struct {
	ID          int     `json:"id"`
	GroupID     int     `json:"group_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Volume      float64 `json:"volume"`
	Mass        float64 `json:"mass"`
	IsSystem    bool    `json:"is_system"`
}

type Asset struct {
	InstanceID int64    `json:"instance_id"`
	TypeID     int      `json:"type_id"`
	Quantity   int      `json:"quantity"`
	Location   string   `json:"location"`
	Item       ItemType `json:"item"`
}

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) GetTypeByID(id int) (*ItemType, error) {
	var t ItemType
	var isSys int
	err := s.DB.QueryRow(`
		SELECT id, group_id, name, IFNULL(description,''), volume, mass, is_system
		FROM item_types WHERE id = ?
	`, id).Scan(&t.ID, &t.GroupID, &t.Name, &t.Description, &t.Volume, &t.Mass, &isSys)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item %d not found", id)
	}
	if err != nil {
		return nil, err
	}
	t.IsSystem = isSys == 1
	return &t, nil
}

func (s *Service) GetCharacterAssets(characterID int) ([]Asset, error) {
	rows, err := s.DB.Query(`
		SELECT ca.instance_id, ca.type_id, ca.quantity, ca.location
		FROM character_assets ca
		WHERE ca.character_id = ?
		ORDER BY ca.location, ca.instance_id
	`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.InstanceID, &a.TypeID, &a.Quantity, &a.Location); err != nil {
			return nil, err
		}
		// Загружаем информацию о типе предмета
		t, err := s.GetTypeByID(a.TypeID)
		if err != nil {
			continue // пропускаем сломанные записи
		}
		a.Item = *t
		assets = append(assets, a)
	}
	return assets, nil
}
