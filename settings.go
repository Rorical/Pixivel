package pixivel

var leveldbConf = LevelDBSetting{
	File: "a.db",
}
var redisConf = RedisSetting{
	IdleTimeout: 240,
	Password:    "",
	redisURL:    "redis://localhost:6379/0",
	MaxIdle:     3,
}
var databaseConf = DatabaseSetting{
	Type: "sqlite3",
	URI:  "s.db",
}
