USE othello;

CREATE TABLE IF NOT EXISTS games (
    play_id VARCHAR(36) PRIMARY KEY,
    black_count INT DEFAULT NULL,
    white_count INT DEFAULT NULL,
    result ENUM('black_win', 'white_win', 'draw') DEFAULT NULL,
    host_secret VARCHAR(36) DEFAULT NULL,
    guest_secret VARCHAR(36) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS moves (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    play_id VARCHAR(36) NOT NULL,
    color ENUM('black', 'white') NOT NULL,
    col TINYINT NOT NULL,
    `row` TINYINT NOT NULL,
    move_order INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (play_id) REFERENCES games(play_id),
    UNIQUE KEY uk_play_move (play_id, move_order)
);

USE othello_test;

CREATE TABLE IF NOT EXISTS games (
    play_id VARCHAR(36) PRIMARY KEY,
    black_count INT DEFAULT NULL,
    white_count INT DEFAULT NULL,
    result ENUM('black_win', 'white_win', 'draw') DEFAULT NULL,
    host_secret VARCHAR(36) DEFAULT NULL,
    guest_secret VARCHAR(36) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS moves (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    play_id VARCHAR(36) NOT NULL,
    color ENUM('black', 'white') NOT NULL,
    col TINYINT NOT NULL,
    `row` TINYINT NOT NULL,
    move_order INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (play_id) REFERENCES games(play_id),
    UNIQUE KEY uk_play_move (play_id, move_order)
);
