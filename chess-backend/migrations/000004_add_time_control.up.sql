-- Time control and clock columns
ALTER TABLE games
    ADD COLUMN time_base INT NOT NULL DEFAULT 0,
    ADD COLUMN time_increment INT NOT NULL DEFAULT 0,
    ADD COLUMN white_remaining BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN black_remaining BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN last_move_at TIMESTAMP NOT NULL DEFAULT NOW();
