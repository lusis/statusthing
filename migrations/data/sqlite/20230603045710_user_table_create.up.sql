CREATE TABLE IF NOT EXISTS users
	(
		id VARCHAR(191) PRIMARY KEY,
		username VARCHAR(191) NOT NULL UNIQUE,
		password BLOB NOT NULL,
		first_name VARCHAR(191),
		last_name VARCHAR(191),
		avatar_url VARCHAR(191),
		email_address VARCHAR(191) NOT NULL UNIQUE,
		last_login INT DEFAULT NULL,
		created INT NOT NULL,
		updated INT NOT NULL,
		deleted INT DEFAULT NULL
	);