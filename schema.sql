USE kosmos_db;

-- 🆔 Аккаунты (2,000,000 – 2,999,999)
CREATE TABLE IF NOT EXISTS accounts (
    id INT UNSIGNED PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_account_id CHECK (id BETWEEN 2000000 AND 2999999)
);

-- 👤 Персонажи (3,000,000 – 3,999,999)
CREATE TABLE IF NOT EXISTS characters (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT, -- AUTO_INCREMENT разрешён, CHECK всё равно сработает
    account_id INT UNSIGNED NOT NULL,
    name VARCHAR(100) UNIQUE NOT NULL,
    race VARCHAR(50) DEFAULT 'Unknown',
    bloodline VARCHAR(50) DEFAULT 'Unknown',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_char_id CHECK (id BETWEEN 3000000 AND 3999999),
    CONSTRAINT fk_char_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

-- 📦 Типы предметов (1 – 500,000)
CREATE TABLE IF NOT EXISTS item_types (
    id INT UNSIGNED PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    volume DECIMAL(10,3) DEFAULT 0.0,
    CONSTRAINT chk_item_id CHECK (id BETWEEN 1 AND 500000)
);

-- 🏢 NPC Корпорации (1,000,000 – 1,999,999)
CREATE TABLE IF NOT EXISTS npcs_corps (
    id INT UNSIGNED PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    faction_id INT UNSIGNED,
    CONSTRAINT chk_corp_id CHECK (id BETWEEN 1000000 AND 1999999)
);

-- 🔹 Seed-данные для теста
INSERT IGNORE INTO accounts (id, login, password) VALUES 
(2000001, 'commander', 'eve_online_2026'),
(2000002, 'testpilot', '12345');