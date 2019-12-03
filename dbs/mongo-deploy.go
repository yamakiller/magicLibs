package dbs

//MongoDeplay doc
//@Summary  mongo config informat
//@Struct MongoDeplay
//@Member ([]string) Mongo DB Host Address array
//@Member (string) Mongo DB URI
//@Member (int) Mongo DB Max connection Pool of number
//@Member (int) Mongo DB Min connection Pool of number
//@Member (int) Mongo DB connection heart beat interval [Unit second]
//@Member (int) Mongo DB connection idle time [Unit second]
//@Member (int) Mongo DB Background time out  [Unit second]
//@Member (int) Mongo DB connection socket time out [Unit second]
type MongoDeplay struct {
	Host              []string `xml:"host" yaml:"host" json:"host"`
	URI               string   `xml:"uri" yaml:"uri" json:"uri"`
	DBName            string   `xml:"db-name" yaml:"db-name" json:"db-name"`
	MaxPoolSize       int      `xml:"max-pool-size" yaml:"max-pool-size" json:"max-pool-size"`
	MinPoolSize       int      `xml:"min-pool-size" yaml:"min-pool-size" json:"min-pool-size"`
	HeartbeatInterval int      `xml:"heart-beat-interval" yaml:"heart-beat-interval" json:"heart-beat-interval"`
	IdleTime          int      `xml:"idle-time" yaml:"idle-time" json:"idle-time"`
	TimeOut           int      `xml:"time-out" yaml:"time-out" json:"time-out"`
	SocketTimeout     int      `xml:"socket-timeout" yaml:"socket-timeout" json:"socket-timeout"`
}
