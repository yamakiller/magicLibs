package util

import "github.com/google/uuid"

//SpawnUUID desc
//@Method SpawnUUID desc: spawn uuid
//@Return (string) uuid
func SpawnUUID() string {
	guid := uuid.New()
	return guid.String()
}
