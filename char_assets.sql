DROP TABLE IF EXISTS character_assets;

CREATE TABLE character_assets (
    instance_id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    character_id INT UNSIGNED NOT NULL,
    type_id INT UNSIGNED NOT NULL,
    quantity INT UNSIGNED DEFAULT 1,
    location VARCHAR(50) DEFAULT 'dock',  -- 'dock' / 'cargo' / 'fitted'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_asset_character FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    CONSTRAINT fk_asset_type FOREIGN KEY (type_id) REFERENCES item_types(id)
) ENGINE=InnoDB;

-- Индексы для быстрой выборки
CREATE INDEX idx_assets_character ON character_assets(character_id);
CREATE INDEX idx_assets_type ON character_assets(type_id);