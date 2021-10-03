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
//@Member (int) Mongo DB Background time out  [Unit millisecond]
//@Member (int) Mongo DB connectiosn socket time out [Unit second]
type MongoDeplay struct {
	Hosts          []string `xml:"hosts" yaml:"hosts" json:"hosts"`
	URI            string   `xml:"uri" yaml:"uri" json:"uri"`
	DBName         string   `xml:"name" yaml:"name" json:"name"`
	MaxPoolSize    int      `xml:"maxpool" yaml:"maxpool" json:"maxpool"`
	MinPoolSize    int      `xml:"minpool" yaml:"minpool" json:"minpool"`
	HeartRate      int      `xml:"heartrate" yaml:"heartrate" json:"heartrate"`
	IdleTime       int      `xml:"idletime" yaml:"idletime" json:"idletime"`
	TimeOut        int      `xml:"timeout" yaml:"timeout" json:"timeout"`
	ConnectTimeout int      `xml:"conntimeout" yaml:"conntimeout" json:"conntimeout"`
}
