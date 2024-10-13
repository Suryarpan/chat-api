package apiconf

import "os"

func DBUrlConfig() string {
	dbUrl, ok := os.LookupEnv("CHAT_API_DB_URL")
	if !ok {
		panic("Error: could not find CHAT_API_DB_URL")
	}
	return dbUrl
}
