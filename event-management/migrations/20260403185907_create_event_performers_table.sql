-- +goose Up
CREATE TABLE IF NOT EXISTS event_performers (
	event_id BIGINT NOT NULL,
	performer_id BIGINT NOT NULL,
	PRIMARY KEY (event_id, performer_id),
	CONSTRAINT fk_event_performers_event FOREIGN KEY (event_id)
		REFERENCES events (id)
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	CONSTRAINT fk_event_performers_performer FOREIGN KEY (performer_id)
		REFERENCES performers (id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_event_performers_performer_id ON event_performers (performer_id);

-- +goose Down
DROP TABLE IF EXISTS event_performers;
