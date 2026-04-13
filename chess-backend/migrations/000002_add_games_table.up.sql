CREATE TABLE IF NOT EXISTS games (
                                     id VARCHAR(36) PRIMARY KEY,
                                     white_player_id VARCHAR(36) REFERENCES users(id) ON DELETE SET NULL,
                                     black_player_id VARCHAR(36) REFERENCES users(id) ON DELETE SET NULL,
                                     status VARCHAR(20) NOT NULL DEFAULT 'waiting',
                                     fen TEXT,
                                     turn VARCHAR(5) NOT NULL,
                                     moves JSONB NOT NULL DEFAULT '[]',
                                     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_games_status ON games(status);