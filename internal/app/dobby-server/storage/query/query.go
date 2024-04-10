package query

const (
	TruncateTable     = `DELETE FROM %s;`
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
	UpdateConfigTable = `UPDATE Config SET 
		baseUrl = ?, 
		gSheetTokenFile = ?, 
		gSheetCredFile = ?, 
		gSheetId = ?;`
	EnsureConfigRow = `INSERT INTO Config (
					baseUrl, 
					gSheetTokenFile, 
					gSheetCredFile, 
					gSheetId)
					SELECT '', '', '', '' WHERE NOT EXISTS (SELECT 1 FROM Config);`

	SelectPotionSubTable = `SELECT * FROM PotionSubforumConfig;`
	CreatePotionSubTable = `CREATE TABLE IF NOT EXISTS PotionSubforumConfig (
        "url" TEXT PRIMARY KEY,
        "timeLimit" INTEGER NOT NULL,
        "turnLimit" INTEGER NOT NULL
    );`
	InsertPotionSubTable = `INSERT INTO PotionSubforumConfig (
		url, 
		timeLimit, 
		turnLimit)
		VALUES (?, ?, ?);`

	SelectPotionThrTable = `SELECT * FROM PotionThreadConfig;`
	CreatePotionThrTable = `CREATE TABLE IF NOT EXISTS PotionThreadConfig (
        "url" TEXT PRIMARY KEY,
        "timeLimit" INTEGER NOT NULL,
        "turnLimit" INTEGER NOT NULL
    );`
	InsertPotionThrTable = `INSERT INTO PotionThreadConfig (
		url,
		timeLimit,
		turnLimit)
		VALUES (?, ?, ?);`

	SelectCreationChamberSubTable = `SELECT * FROM CreationChamberSubforumConfig;`
	CreateCreationChamberSubTable = `CREATE TABLE IF NOT EXISTS CreationChamberSubforumConfig (
        "url" TEXT PRIMARY KEY,
        "timeLimit" INTEGER NOT NULL,
        "turnLimit" INTEGER NOT NULL
    );`
	InsertCreationChamberSubTable = `INSERT INTO CreationChamberSubforumConfig (
		url, 
		timeLimit, 
		turnLimit)
		VALUES (?, ?, ?);`
)
