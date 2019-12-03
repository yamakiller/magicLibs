package dbs

//MySQLDeploy doc
//@Summary mysql connection deploy
//@Struct MySQLDeploy
//@Member (string) mysql dsn config informat
//@Member (int) mysql connection max of number
//@Member (int) mysql connection idle of number
//@Member (int) mysql connection life time (sec)
type MySQLDeploy struct {
	DSN  string `xml:"dsn" yaml:"dsn" json:"dsn"`
	Max  int    `xml:"max-conn" yaml:"max-conn" json:"max-conn"`
	Idle int    `xml:"idle-conn" yaml:"idle-conn" json:"idle-conn"`
	Life int    `xml:"life-time" yaml:"life-time" json:"life-time"`
}

//DoMySQLDeploy doc
//@Summary deploy mysql db
//@Method DoMySQLDeploy
//@Param  (*MySQLDB) mysql object
//@Param  (*MySQLDeploy) mysql config informat
//@Return (error) register mysql success/fail
func DoMySQLDeploy(Db *MySQLDB, deploy *MySQLDeploy) error {
	return Db.Initial(deploy.DSN, deploy.Max, deploy.Idle, deploy.Life)
}
