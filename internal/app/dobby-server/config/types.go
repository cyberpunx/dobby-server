package config

type Config struct {
	BaseUrl         *string `json:"baseUrl" meta-obscure:"default"`
	GSheetTokenFile *string `json:"gSheetTokenFile" meta-obscure:"default"`
	GSheetCredFile  *string `json:"gSheetCredFile" meta-obscure:"default"`
	GSheetId        *string `json:"gSheetId" meta-obscure:"default"`
	TursoDbUrl      *string `json:"tursoDbUrl" meta-obscure:"default"`
	TursoDbToken    *string `json:"tursoDbToken" meta-obscure:"default"`
	ServerPort      *string `json:"serverPort" meta-obscure:"default"`
	UnicodeOutput   *bool   `json:"unicodeOutput" meta-obscure:"default"`
}

type Task struct {
	Urls      *[]string `json:"urls" meta-obscure:"default"`
	Method    *string   `json:"method" meta-obscure:"default"`
	TimeLimit *int      `json:"timeLimit" meta-obscure:"default"`
	TurnLimit *int      `json:"turnLimit" meta-obscure:"default"`
}

type PotionSubforumConfig struct {
	Url       *string `json:"url" meta-obscure:"default"`
	TimeLimit *int    `json:"timeLimit" meta-obscure:"default"`
	TurnLimit *int    `json:"turnLimit" meta-obscure:"default"`
}

type CreationChamberSubforumConfig struct {
	Url       *string `json:"url" meta-obscure:"default"`
	TimeLimit *int    `json:"timeLimit" meta-obscure:"default"`
	TurnLimit *int    `json:"turnLimit" meta-obscure:"default"`
}

type PotionThreadConfig struct {
	Url       *string `json:"url" meta-obscure:"default"`
	TimeLimit *int    `json:"timeLimit" meta-obscure:"default"`
	TurnLimit *int    `json:"turnLimit" meta-obscure:"default"`
}

const (
	LogTagInfo    = "Blinky!"
	LogTagPotions = "potionsClub"
	passphrase    = "yourEncryptionKey"
)

var config *Config
var loadedConfigPath string
var Reset = ""
var Red = ""
var Green = ""
var Yellow = ""
var Blue = ""
var Purple = ""
var Cyan = ""
var Gray = ""
var White = ""
var CheckEmoji = ""
var CrossEmoji = ""
var RightArrowEmoji = ""

func (c *Config) Validate() {
	if c.TursoDbUrl == nil {
		panic("Missing TursoDbUrl")
	}
	if c.TursoDbToken == nil {
		panic("Missing TursoDbToken")
	}
}

func MergeConfigs(fileConfig, conf Config) Config {
	// Inicializar una nueva instancia de Config para almacenar los resultados combinados.
	var mergedConfig Config

	// Combinar BaseUrl
	if conf.BaseUrl != nil {
		mergedConfig.BaseUrl = conf.BaseUrl
	} else {
		mergedConfig.BaseUrl = fileConfig.BaseUrl
	}

	// Combinar GSheetTokenFile
	if conf.GSheetTokenFile != nil {
		mergedConfig.GSheetTokenFile = conf.GSheetTokenFile
	} else {
		mergedConfig.GSheetTokenFile = fileConfig.GSheetTokenFile
	}

	// Combinar GSheetCredFile
	if conf.GSheetCredFile != nil {
		mergedConfig.GSheetCredFile = conf.GSheetCredFile
	} else {
		mergedConfig.GSheetCredFile = fileConfig.GSheetCredFile
	}

	// Combinar GSheetId
	if conf.GSheetId != nil {
		mergedConfig.GSheetId = conf.GSheetId
	} else {
		mergedConfig.GSheetId = fileConfig.GSheetId
	}

	// Combinar TursoDbUrl
	if conf.TursoDbUrl != nil {
		mergedConfig.TursoDbUrl = conf.TursoDbUrl
	} else {
		mergedConfig.TursoDbUrl = fileConfig.TursoDbUrl
	}

	// Combinar TursoDbToken
	if conf.TursoDbToken != nil {
		mergedConfig.TursoDbToken = conf.TursoDbToken
	} else {
		mergedConfig.TursoDbToken = fileConfig.TursoDbToken
	}

	// Combinar ServerPort
	if conf.ServerPort != nil {
		mergedConfig.ServerPort = conf.ServerPort
	} else {
		mergedConfig.ServerPort = fileConfig.ServerPort
	}

	return mergedConfig
}
