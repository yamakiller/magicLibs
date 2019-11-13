package magicDB

//MongoDeplay desc
//@struct MongoDeplay desc: mongo config informat
//@member ([]string) Mongo DB Host Address array
//@member (string) Mongo DB URI
//@member (int) Mongo DB Max connection Pool of number
//@member (int) Mongo DB Min connection Pool of number
//@member (int) Mongo DB connection heart beat interval [Unit second]
//@member (int) Mongo DB connection idle time [Unit second]
//@member (int) Mongo DB Background time out  [Unit second]
//@member (int) Mongo DB connection socket time out [Unit second]
type MongoDeplay struct {
	Host              []string `json:"host"`
	URI               string   `json:"uri"`
	DBName            string   `json:"db-name"`
	MaxPoolSize       int      `json:"max-pool-size"`
	MinPoolSize       int      `json:"min-pool-size"`
	HeartbeatInterval int      `json:"heart-beat-interval"`
	IdleTime          int      `json:"idle-time"`
	TimeOut           int      `json:"time-tou"`
	SocketTimeout     int      `json:"socket-timeout"`
}
