package config

type Config struct {
	DB     *DB     `yaml:"db"`
	Static *Static `yaml:"static"`
	Server *Server `yaml:"server"`
}
type DB struct {
	Account         string `yaml:"account"`
	Password        string `yaml:"password"`
	Address         string `yaml:"address"`
	Port            string `yaml:"port"`
	ConnectProtocol string `yaml:"connectProtocol"`
	Database        string `yaml:"database"`
	CharSet         string `yaml:"charSet"`
	TimeSet         string `yaml:"timeSet"`
}

func (db *DB) GetMySQLDNS() string {
	return db.Account + ":" + db.Password + "@" + db.ConnectProtocol + "(" + db.Address + ")/" + db.Database + "?" + db.CharSet + "&" + db.TimeSet
}

type Static struct {
	HttpFolder     string `yaml:"httpFolder"`
	LocalStorePath string `yaml:"localStorePath"`
}

type Server struct {
	Protocol string `yaml:"protocol"`
	IP       string `yaml:"ip"`
	Port     string `yaml:"port"`
}

func (s *Server) GetPortLikeInDomain() string {
	return ":" + s.Port
}

//GetStaticUrl the prefix of the url of video
func (s *Static) GetStaticUrl() string {
	return CServer.Protocol + "://" + CServer.IP + ":" + CServer.Port + "/" + s.HttpFolder
}
