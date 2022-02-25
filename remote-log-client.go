package remotelog

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/google/uuid"
)

type (
	remotelogClient struct {
		logger *log.Logger
	}

	ConnectedClient interface {
		Close() error
	}
)

func NewClient(options InstanceOptions) remotelogClient {
	instance := remotelogClient{}
	if options.EnableLogging {
		instance.logger = log.New(os.Stdout, "", log.Ltime|log.Ldate)
	} else {
		instance.logger = log.New(ioutil.Discard, "", 0)
	}
	return instance
}

func (c remotelogClient) ConnectTcp(options ServerOptions) (ConnectedClient, error) {
	client := tcpClient{
		id:     uuid.NewString(),
		logger: c.logger,
	}

	return client, client.connect(options)
}
