package amqpd

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	GlobalConfig = Config{
		Host:         "127.0.0.1",
		Port:         "5672",
		Username:     "guest",
		Password:     "guest",
		Vhost:        "",
		TlsProtocols: false,
	}

	if err := Init(); err != nil {
		fmt.Printf("init error: %s\n", err)
		os.Exit(-1)
	}
	os.Exit(m.Run())
}

func TestAmqpdPublish(t *testing.T) {
	cli, err := New()
	require.NoError(t, err)
	defer cli.Close()

	for i := 0; i < 10000; i++ {
		err = cli.Publish("test_exchanges", "", []byte(fmt.Sprintf("%d", i)))
		require.NoError(t, err, "publish error i: %d", i)
	}
}

func TestAmqpdConsumer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	fn := func(msg []byte) error {
		fmt.Println(string(msg))
		return nil
	}

	ac, err := NewAmqpdConsumer()
	require.NoError(t, err)
	ac.AddFunc("test_queues", "test-consumer", fn)

	go func() {
		time.Sleep(time.Millisecond * 100)
		<-ac.Stop().Done()
		fmt.Println("ac stop...")
		cancel()
	}()
	ac.Start()

	<-ctx.Done()
}
