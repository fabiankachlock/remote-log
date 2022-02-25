package remotelog

import "github.com/google/uuid"

type (
	remotelogClient struct{}

	ConnectedClient interface {
		Close() error
	}
)

func NewClient() remotelogClient {
	return remotelogClient{}
}

func (c remotelogClient) ConnectTcp(options ServerOptions) (ConnectedClient, error) {
	client := tcpClient{
		id: uuid.NewString(),
	}

	return client, client.connect(options)
}
