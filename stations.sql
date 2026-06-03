USE kosmos_db;

-- 1. Таблица станций (ID начинаются с 50_000_000)
CREATE TABLE IF NOT EXISTS stations (
    id INT UNSIGNED PRIMARY KEY,
    system_id INT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    CONSTRAINT fk_station_system FOREIGN KEY (system_id) REFERENCES solar_systems(id)
);

-- 2. Добавляем текущую станцию в таблицу персонажей
ALTER TABLE characters 
ADD COLUMN current_station_id INT UNSIGNED DEFAULT 50000001;

-- 3. Seed: Создаём стартовую станцию в нашей единственной системе
INSERT INTO stations (id, system_id, name) VALUES 
(50000001, 42000001, 'Starter Outpost');

-- 4. Привязываем всех существующих персонажей к этой станции
UPDATE characters SET current_station_id = 50000001 WHERE current_station_id IS NULL;