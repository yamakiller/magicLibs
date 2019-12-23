package dbs

//MySQLGormDeploy doc
//@Summary mysql connection deploy
//@Member (string) mysql dsn config informat
//@Member (int) mysql connection max of number
//@Member (int) mysql connection idle of number
//@Member (int) mysql connection life time (sec)
type MySQLGormDeploy struct {
	DSN  string `xml:"dsn" yaml:"dsn" json:"dsn"`
	Max  int    `xml:"max" yaml:"max conn" json:"max"`
	Idle int    `xml:"idle" yaml:"idle" json:"idle"`
	Life int    `xml:"life" yaml:"life" json:"life"`
}

//DoMySQLGormDeploy doc
//@Summary deploy mysql db
//@Method DoMySQLDeploy
//@Param  (*MySQLDB) mysql object
//@Param  (*MySQLDeploy) mysql config informat
//@Return (error) register mysql success/fail
func DoMySQLGormDeploy(Db *MySQLGORM, deploy *MySQLDeploy) error {
	return Db.Initial(deploy.DSN, deploy.Max, deploy.Idle, deploy.Life)
}
