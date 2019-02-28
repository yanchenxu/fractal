package sqlite

var (
	CREATE_SQL = `
	CREATE TABLE IF NOT EXISTS userinfo(
	uid INTEGER PRIMARY KEY AUTOINCREMENT,
	username VARCHAR(64) NULL,
	departname VARCHAR(64) NULL,
	created DATE NULL);
`
)
