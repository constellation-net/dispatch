package env

import (
	"os"
	"strconv"
	"time"

	"github.com/constellation-net/dispatch/utils"
)

// Vars is populated at runtime with all the environment variables required by this program
var Vars *Environment

// Environment wraps all env vars recognised by this program
type Environment struct {
	dispatchIntervalSeconds int
	Relay                   Relay
	Upstream                Upstream
}

func (e *Environment) DispatchInterval() time.Duration {
	return time.Duration(e.dispatchIntervalSeconds) * time.Second
}

// Relay holds information about the local SMTP server that clients will use a relay
type Relay struct {
	Host     string
	Port     uint16
	Username string
	Password []byte
}

func (r *Relay) Address() string {
	return r.Host + ":" + strconv.Itoa(int(r.Port))
}

// Upstream holds environment variables about the upstream SMTP server to forward incoming messages to
type Upstream struct {
	Host     string
	Port     uint16
	Username string
	Password string
	From     string
	Replyto  string
}

func (u *Upstream) Address() string {
	return u.Host + ":" + strconv.Itoa(int(u.Port))
}

func init() {
	rPort, err := utils.ConvertPort(os.Getenv("RELAY_PORT"))
	if err != nil {
		panic(err)
	}

	upPort, err := utils.ConvertPort(os.Getenv("UPSTREAM_PORT"))
	if err != nil {
		panic(err)
	}

	disInterval, err := strconv.Atoi(os.Getenv("DISPATCH_INTERVAL"))
	if err != nil {
		panic(err)
	}

	Vars = &Environment{
		dispatchIntervalSeconds: disInterval,
		Relay: Relay{
			Host:     os.Getenv("RELAY_HOST"),
			Port:     rPort,
			Username: os.Getenv("RELAY_USERNAME"),
			Password: []byte(os.Getenv("RELAY_PASSWORD")),
		},
		Upstream: Upstream{
			Host:     os.Getenv("UPSTREAM_HOST"),
			Port:     upPort,
			Username: os.Getenv("UPSTREAM_USERNAME"),
			Password: os.Getenv("UPSTREAM_PASSWORD"),
			From:     os.Getenv("UPSTREAM_FROM"),
			Replyto:  os.Getenv("UPSTREAM_REPLYTO"),
		},
	}
}
