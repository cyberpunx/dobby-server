package query

const (
	SelectConfigTable = `SELECT * FROM Config LIMIT 1;`

	CreateConfigTable = `CREATE TABLE IF NOT EXISTS Config (
        "baseUrl" TEXT,
        "gSheetTokenFile" TEXT,
        "gSheetCredFile" TEXT,
        "gSheetModeracionId" TEXT
    );`
	InsertConfigTable = `INSERT INTO Config (
		baseUrl,
		gSheetTokenFile,
		gSheetCredFile,
		gSheetModeracionId)
		VALUES (?, ?, ?, ?);`
)
