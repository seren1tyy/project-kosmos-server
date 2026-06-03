USE kosmos_db;

-- 1. Таблицы вселенной
CREATE TABLE IF NOT EXISTS regions (
    id INT UNSIGNED PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS constellations (
    id INT UNSIGNED PRIMARY KEY,
    region_id INT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    CONSTRAINT fk_const_region FOREIGN KEY (region_id) REFERENCES regions(id)
);

CREATE TABLE IF NOT EXISTS solar_systems (
    id INT UNSIGNED PRIMARY KEY,
    constellation_id INT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    CONSTRAINT fk_sys_const FOREIGN KEY (constellation_id) REFERENCES constellations(id)
);

-- 2. Обновляем таблицу персонажей (добавляем текущую систему)
ALTER TABLE characters 
ADD COLUMN current_system_id INT UNSIGNED DEFAULT 42000001;

-- 3. Seed: Единственный регион, созвездие и система (1-к-1-к-1)
-- Используем диапазоны: Регионы 40M+, Созвездия 41M+, Системы 42M+
INSERT INTO regions (id, name) VALUES 
(40000001, 'The Cradle');

INSERT INTO constellations (id, region_id, name) VALUES 
(41000001, 40000001, 'Nursery');

INSERT INTO solar_systems (id, constellation_id, name) VALUES 
(42000001, 41000001, 'Starting System');

-- Обновляем существующих персонажей, чтобы они были в этой системе
UPDATE characters SET current_system_id = 42000001 WHERE current_system_id IS NULL;