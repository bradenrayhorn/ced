package sqlite

var migrations = []string{
	`CREATE TABLE groups (
		id CHAR(27) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		max_attendees INT(4) NOT NULL,
		attendees INT(4) NOT NULL,
		has_responded INT(1) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,
	`ALTER TABLE groups ADD COLUMN search_hints VARCHAR(255) DEFAULT "";`,
}
