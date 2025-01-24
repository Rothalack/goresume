CREATE TABLE IF NOT EXISTS games (
    id INT AUTO_INCREMENT PRIMARY KEY,
    game_name VARCHAR(255) NOT NULL,
    api_url VARCHAR(255) NOT NULL,
    note VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO games (game_name, api_url, note)
VALUES
("Retail", "https://www.warcraftlogs.com/", "Retail WoW"),
("Classic", "https://classic.warcraftlogs.com/", "Main Line Classic"),
("Classic Fresh", "https://fresh.warcraftlogs.com/", "Classic Fresh"),
("Season of Discovery", "https://sod.warcraftlogs.com/", "Season of Discovery"),
("Vanilla", "https://vanilla.warcraftlogs.com/", "Vanilla");