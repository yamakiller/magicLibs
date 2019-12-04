package dbs

//RedisDeploy doc
//@Summary redis configureation json
//@Struct RedisDeploy
//@Member (string) Redis Host from json "host"
//@Member (int) Redis DB select from json "db"
//@Member (int) Redis Connection max idle of number from json "idle"
//@Member (int) Redis Connection max active of number from json "active"
//@Member (int) Redis Connection idle time from json "idle-time" util/sec
type RedisDeploy struct {
	Host      string `xml:"host" yaml:"host" json:"host"`
	DB        int    `xml:"db" yaml:"db" json:"db"`
	MaxIdle   int    `xml:"idle" yaml:"idle" json:"idle"`
	MaxActive int    `xml:"active" yaml:"active" json:"active"`
	IdleTime  int    `xml:"idletime" yaml:"idletime" json:"idletime"`
}

//RedisDeployArray doc
//@Summary redis group configureation json
//@Struct RedisDeployArray
//@Member ([]RedisDeploy) a redis configureation
type RedisDeployArray struct {
	Deploys []RedisDeploy `xml:"redis" yaml:"redis" json:"redis"`
}

//DoRedisDeploy doc
//@Summary deploy mysql db
//@Method DoRedisDeploy
//@Param  (*MySQLDB) mysql object
//@Param  (*MySQLDeployArray) mysql config informat
//@Return (error) register mysql success/fail
func DoRedisDeploy(Db *RedisDB, deploy *RedisDeploy) error {
	return Db.Initial(deploy.Host, deploy.DB, deploy.MaxIdle, deploy.MaxActive, deploy.IdleTime)
}
