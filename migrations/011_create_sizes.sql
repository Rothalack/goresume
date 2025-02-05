CREATE TABLE IF NOT EXISTS sizes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    difficulty_id INT NOT NULL,
    zone_id INT NOT NULL,
    size INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY (difficulty_id, zone_id, size),
    FOREIGN KEY (difficulty_id) REFERENCES difficulties(difficulty_id) ON DELETE CASCADE
);