USE kosmos_db;

-- Удаление старых таблиц в правильном порядке (FK-зависимости)
DROP TABLE IF EXISTS item_types;
DROP TABLE IF EXISTS item_groups;
DROP TABLE IF EXISTS item_categories;

-- =========================
-- КАТЕГОРИИ (ID 1 - 999)
-- =========================
CREATE TABLE item_categories (
    id INT UNSIGNED PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description VARCHAR(500),
    CONSTRAINT chk_category_id CHECK (id BETWEEN 1 AND 999)
);

-- =========================
-- ГРУППЫ (ID 1000 - 9999)
-- =========================
CREATE TABLE item_groups (
    id INT UNSIGNED PRIMARY KEY,
    category_id INT UNSIGNED NOT NULL,
    name VARCHAR(100) UNIQUE NOT NULL,
    description VARCHAR(500),
    CONSTRAINT chk_group_id CHECK (id BETWEEN 1000 AND 9999),
    CONSTRAINT fk_group_category FOREIGN KEY (category_id) REFERENCES item_categories(id)
);

-- =========================
-- ТИПЫ ПРЕДМЕТОВ (ID 1 - 4_999_999)
-- =========================
CREATE TABLE item_types (
    id INT UNSIGNED PRIMARY KEY,
    group_id INT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1000),
    volume DECIMAL(10,3) DEFAULT 0.0,       -- место в трюме (м³)
    mass DECIMAL(15,3) DEFAULT 0.0,         -- масса (кг)
    is_system TINYINT(1) DEFAULT 0,         -- 1 = системный предмет (не крафтится, не продаётся)
    CONSTRAINT chk_item_id CHECK (id BETWEEN 1 AND 4999999),
    CONSTRAINT fk_item_group FOREIGN KEY (group_id) REFERENCES item_groups(id)
);

-- =========================
-- SEED: КАТЕГОРИИ
-- =========================
INSERT INTO item_categories (id, name, description) VALUES
(1, 'System',       'Internal system objects (regions, systems, stations)'),
(2, 'Corporation',  'Corporate entities'),
(3, 'Celestial',    'Stars, planets, moons, belts'),
(4, 'Ship',         'Player-controllable vessels'),
(5, 'Module',       'Ship fittings (weapons, armor, electronics)'),
(6, 'Charge',       'Ammo and consumables for modules'),
(7, 'Drone',        'Remote-controlled fighters'),
(8, 'Structure',    'Stations, citadels, POSes');

-- =========================
-- SEED: ГРУППЫ
-- =========================
INSERT INTO item_groups (id, category_id, name, description) VALUES
-- Системные группы (по одной на каждый "предмет-концепцию" из запроса)
(1000, 1, 'Entity',        'Generic system entity marker (ID 1)'),
(1001, 2, 'Corporation',   'Player/NPC corporation (ID 2)'),
(1002, 3, 'Region',        'Galactic region (ID 3)'),
(1003, 3, 'Constellation', 'Constellation (ID 4)'),
(1004, 3, 'Solar System',  'Solar system (ID 5)'),
(1005, 3, 'Star',          'Star — various classes (ID 6-10)'),
(1006, 3, 'Planet',        'Planet — Temperate/Barren/Gas Giant (ID 11-13)'),
(1007, 3, 'Moon',          'Moon (ID 14)'),
(1008, 3, 'Asteroid Belt', 'Asteroid belt (ID 15)'),

-- Игровые группы
(1100, 4, 'Frigate',       'Light combat vessels'),
(1101, 4, 'Cruiser',       'Medium combat vessels'),
(1102, 4, 'Battleship',    'Heavy combat vessels'),
(1103, 4, 'Industrial',    'Haulers and miners'),
(1200, 5, 'Weapon',        'Turrets and launchers'),
(1201, 5, 'Armor',         'Armor plates and repairers'),
(1202, 5, 'Shield',        'Shield extenders and boosters'),
(1203, 5, 'Electronics',   'Sensor, propulsion modules');

-- =========================
-- SEED: СИСТЕМНЫЕ ПРЕДМЕТЫ (ID 1-15, is_system=1)
-- =========================
INSERT INTO item_types (id, group_id, name, description, is_system) VALUES
(1,  1000, '#System',        'Root system marker', 1),
(2,  1001, 'Corporation',    'Corporation marker', 1),
(3,  1002, 'Region',         'Region marker', 1),
(4,  1003, 'Constellation',  'Constellation marker', 1),
(5,  1004, 'Solar System',   'Solar system marker', 1),
-- Звёзды (разные классы)
(6,  1005, 'Star (G-Class)',      'Yellow main-sequence star', 1),
(7,  1005, 'Star (K-Class)',      'Orange dwarf star', 1),
(8,  1005, 'Star (M-Class)',      'Red dwarf star', 1),
(9,  1005, 'Star (F-Class)',      'White-yellow star', 1),
(10, 1005, 'Star (Black Hole)',   'Collapsed stellar remnant', 1),
-- Планеты (разные типы)
(11, 1006, 'Planet (Temperate)',  'Earth-like habitable world', 1),
(12, 1006, 'Planet (Barren)',     'Rocky dead world', 1),
(13, 1006, 'Planet (Gas Giant)',  'Massive gas planet', 1),
(14, 1007, 'Moon',                'Natural satellite', 1),
(15, 1008, 'Asteroid Belt',       'Field of rocky debris', 1);

-- =========================
-- SEED: СТАРТОВЫЙ КОРАБЛЬ (ID 582, Wisp)
-- =========================
INSERT INTO item_types (id, group_id, name, description, volume, mass) VALUES
(582, 1100, 'Wisp', 
 'A nimble starter frigate. Perfect for new pilots learning the ropes of space travel. Low cost, high agility.', 
 10.5, 1050.0);