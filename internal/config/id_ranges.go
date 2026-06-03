package config

const (
	// Категории предметов
	MinCategoryID, MaxCategoryID = 1, 999

	// Группы предметов
	MinGroupID, MaxGroupID = 1000, 9999

	// Типы предметов (включая системные ID 1-15)
	MinItemID, MaxItemID = 1, 4_999_999

	// Фракции
	MinFactionID, MaxFactionID = 5_000_000, 5_999_999

	// NPC-корпорации
	MinNPCCorpID, MaxNPCCorpID = 6_000_000, 6_999_999

	// Альянсы (зарезервировано)
	MinAllianceID, MaxAllianceID = 7_000_000, 7_999_999

	// Персонажи
	MinCharacterID, MaxCharacterID = 10_000_000, 19_999_999

	// Аккаунты
	MinAccountID, MaxAccountID = 20_000_000, 29_999_999

	// Станции/структуры
	MinStationID, MaxStationID = 30_000_000, 39_999_999

	// Системные ID-маркеры (зарезервированы)
	SystemRootID          = 1
	SystemCorporationID   = 2
	SystemRegionID        = 3
	SystemConstellationID = 4
	SystemSolarID         = 5
)
