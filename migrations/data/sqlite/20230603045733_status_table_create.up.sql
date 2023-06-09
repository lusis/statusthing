CREATE TABLE IF NOT EXISTS status
	(
		id VARCHAR(191) PRIMARY KEY,
		name VARCHAR(191) NOT NULL UNIQUE,
		kind VARCHAR(191) NOT NULL,
		description VARCHAR(191) DEFAULT NULL,
		color VARCHAR(191) DEFAULT NULL,
		created INT NOT NULL,
		updated INT NOT NULL,
		deleted INT DEFAULT NULL
	);