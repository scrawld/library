package amqpd

import (
	"fmt"
	"os"
	"testing"

	"github.com/scrawld/library/config"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	config.Get().Rabbitmq.Host = "b-030fdd34-d008-446e-bbee-8c3bb6bec241.mq.ap-southeast-1.amazonaws.com"
	config.Get().Rabbitmq.Port = "5671"
	config.Get().Rabbitmq.Username = "bra88"
	config.Get().Rabbitmq.Password = "TdSoB_e2Lnpzsfl5AzGp9bhHchdzsrjg"
	config.Get().Rabbitmq.Vhost = "bra88"
	config.Get().Rabbitmq.TlsProtocols = true

	if err := Init(); err != nil {
		fmt.Printf("init error: %s\n", err)
		os.Exit(-1)
	}
	os.Exit(m.Run())
}

func TestAmqpd(t *testing.T) {
	for i := 0; i < 3000; i++ {
		cli, err := New()
		require.NoError(t, err, "new amqpd error i: %d", i)

		err = cli.Publish("test_exchanges", "", []byte(fmt.Sprintf("%d", i)))
		require.NoError(t, err, "publish error i: %d", i)
	}
}
