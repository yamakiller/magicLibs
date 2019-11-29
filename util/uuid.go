package util

import "github.com/google/uuid"

//SpawnUUID desc
//@method SpawnUUID desc: spawn uuid
//@return (string) uuid
func SpawnUUID() string {
	guid := uuid.New()
	return guid.String()
}
