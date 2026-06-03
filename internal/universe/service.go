package universe

import (
	"database/sql"
	"fmt"
)

type LocationInfo struct {
	StationID         int    `json:"station_id"`
	StationName       string `json:"station_name"`
	SystemID          int    `json:"system_id"`
	SystemName        string `json:"system_name"`
	ConstellationID   int    `json:"constellation_id"`
	ConstellationName string `json:"constellation_name"`
	RegionID          int    `json:"region_id"`
	RegionName        string `json:"region_name"`
}

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

// GetLocationByStation вытягивает всю иерархию: Станция -> Система -> Созвездие -> Регион
func (s *Service) GetLocationByStation(stationID int) (*LocationInfo, error) {
	var info LocationInfo
	err := s.DB.QueryRow(`
		SELECT 
			st.id, st.name,
			ss.id, ss.name,
			c.id, c.name,
			r.id, r.name
		FROM stations st
		JOIN solar_systems ss ON st.system_id = ss.id
		JOIN constellations c ON ss.constellation_id = c.id
		JOIN regions r ON c.region_id = r.id
		WHERE st.id = ?
	`, stationID).Scan(
		&info.StationID, &info.StationName,
		&info.SystemID, &info.SystemName,
		&info.ConstellationID, &info.ConstellationName,
		&info.RegionID, &info.RegionName,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("station %d not found", stationID)
	}
	return &info, err
}
