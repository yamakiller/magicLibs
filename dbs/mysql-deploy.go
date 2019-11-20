package dbs

//MySQLDeploy desc
//@struct MySQLDeploy desc: mysql connection deploy
//@member (string) mysql dsn config informat
//@member (int) mysql connection max of number
//@member (int) mysql connection idle of number
//@member (int) mysql connection life time (sec)
type MySQLDeploy struct {
	DSN  string `json:"dsn"`
	Max  int    `json:"max-conn"`
	Idle int    `json:"idle-conn"`
	Life int    `json:"life-time"`
}

//DoMySQLDeploy desc
//@method DoMySQLDeploy desc: deploy mysql db
//@param  (*MySQLDB) mysql object
//@param  (*MySQLDeploy) mysql config informat
//@return (error) register mysql success/fail
func DoMySQLDeploy(Db *MySQLDB, deploy *MySQLDeploy) error {
	return Db.Initial(deploy.DSN, deploy.Max, deploy.Idle, deploy.Life)
}
