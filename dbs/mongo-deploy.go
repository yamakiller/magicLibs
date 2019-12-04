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
//@Member (int) Mongo DB connectiosn socket time out [Unit second]
type MongoDeplay struct {
	Hosts             []string `xml:"hosts" yaml:"hosts" json:"hosts"`
	URI               string   `xml:"uri" yaml:"uri" json:"uri"`
	DBName            string   `xml:"db-name" yaml:"db name" json:"db-name"`
	MaxPoolSize       int      `xml:"max-pool" yaml:"max pool" json:"max-pool"`
	MinPoolSize       int      `xml:"min-pool" yaml:"min pool" json:"min-pool"`
	HeartRate 		  int      `xml:"heart-rate" yaml:"heart rate" json:"heart-beat-interval"`
	IdleTime          int      `xml:"idle-time" yaml:"idle time" json:"idle-time"`
	TimeOut           int      `xml:"time-out" yaml:"time out" json:"time-out"`
	ConnectTimeout    int      `xml:"connection-timeout" yaml:"connection timeout" json:"connection-timeout"`
}
