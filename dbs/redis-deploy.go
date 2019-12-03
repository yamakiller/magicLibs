package dbs

//RedisDeploy desc
//@Struct RedisDeploy desc: redis configureation json
//@Member (string) Redis Host from json "host"
//@Member (int) Redis DB select from json "db"
//@Member (int) Redis Connection max idle of number from json "max-idle"
//@Member (int) Redis Connection max active of number from json "max-active"
//@Member (int) Redis Connection idle time from json "idle-time" util/sec
type RedisDeploy struct {
	Host      string `xml:"host" yaml:"host" json:"host"`
	DB        int    `xml:"db" yaml:"db" json:"db"`
	MaxIdle   int    `xml:"max-idle" yaml:"max-idle" json:"max-idle"`
	MaxActive int    `xml:"max-active" yaml:"max-active" json:"max-active"`
	IdleTime  int    `xml:"idle-time" yaml:"idle-time" json:"idle-time"`
}

//RedisDeployArray desc
//@Struct RedisDeployArray desc: redis group configureation json
//@Member ([]RedisDeploy) a redis configureation
type RedisDeployArray struct {
	Deploys []RedisDeploy `xml:"redis" yaml:"redis" json:"redis"`
}

//DoRedisDeploy desc
//@Method DoRedisDeploy desc: deploy mysql db
//@Param  (*MySQLDB) mysql object
//@Param  (*MySQLDeployArray) mysql config informat
//@Return (error) register mysql success/fail
func DoRedisDeploy(Db *RedisDB, deploy *RedisDeploy) error {
	return Db.Initial(deploy.Host, deploy.DB, deploy.MaxIdle, deploy.MaxActive, deploy.IdleTime)
}
