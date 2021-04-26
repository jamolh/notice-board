package db

const createTableNoticesQuery = `CREATE TABLE IF NOT EXISTS notices (
	id UUID PRIMARY KEY,
	title VARCHAR(256) NOT NULL,
	description TEXT NOT NULL,
	price NUMERIC(10, 2) DEFAULT 0,
	created_at TIMESTAMP DEFAULT now(),
	images JSONB DEFAULT '[]',
	UNIQUE(title)
);`

const dropTableNoticesQuery = `DROP TABLE IF EXISTS notices;`
