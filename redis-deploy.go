package magicDB

//RedisDeploy desc
//@struct RedisDeploy desc: redis configureation json
//@member (string) Redis Host from json "host"
//@member (int) Redis DB select from json "db"
//@member (int) Redis Connection max idle of number from json "max-idle"
//@member (int) Redis Connection max active of number from json "max-active"
//@member (int) Redis Connection idle time from json "idle-time" util/sec
type RedisDeploy struct {
	Host      string `json:"host"`
	DB        int    `json:"db"`
	MaxIdle   int    `json:"max-idle"`
	MaxActive int    `json:"max-active"`
	IdleTime  int    `json:"idle-time"`
}

//RedisDeployArray desc
//@struct RedisDeployArray desc: redis group configureation json
//@member ([]RedisDeploy) a redis configureation
type RedisDeployArray struct {
	Deploys []RedisDeploy `json:"redis"`
}

//DoRedisDeploy desc
//@method DoRedisDeploy desc: deploy mysql db
//@param  (*MySQLDB) mysql object
//@param  (*MySQLDeployArray) mysql config informat
//@return (error) register mysql success/fail
func DoRedisDeploy(Db *RedisDB, deploy *RedisDeploy) error {
	return Db.Init(deploy.Host, deploy.DB, deploy.MaxIdle, deploy.MaxActive, deploy.IdleTime)
}
