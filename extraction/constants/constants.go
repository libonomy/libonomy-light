package constants

type commonConstants struct {
	Port     string
	Protocol string
}

var (
	//Server Constants
	Server = commonConstants{
		Port:     ":4500",
		Protocol: "tcp",
	}
)
