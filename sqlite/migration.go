package sqlite

var migrations = []string{
	`CREATE TABLE groups (
		id CHAR(27) PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,

	`CREATE TABLE individuals (
		id CHAR(27) PRIMARY KEY,
		group_id CHAR(27) NOT NULL,
		name VARCHAR(255) NOT NULL,
		response INT(1) NOT NULL,
		has_responded INT(1) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		FOREIGN KEY (group_id) REFERENCES groups (id)
	);`,
}
