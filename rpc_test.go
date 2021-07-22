package rpc_test

import (
	"context"
	"testing"
	"time"

	core "github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rpc "github.com/textileio/go-libp2p-pubsub-rpc"
	"github.com/textileio/go-libp2p-pubsub-rpc/finalizer"
	"github.com/textileio/go-libp2p-pubsub-rpc/peer"
	golog "github.com/textileio/go-log/v2"
	logging "github.com/textileio/go-log/v2"
	"go.uber.org/zap/zapcore"
)

func init() {
	if err := setLogLevels(map[string]golog.LogLevel{
		"psrpc":      golog.LevelDebug,
		"psrpc/peer": golog.LevelDebug,
		"psrpc/mdns": golog.LevelDebug,
	}); err != nil {
		panic(err)
	}
}

func TestPingPong(t *testing.T) {
	fin := finalizer.NewFinalizer()

	p1, err := peer.New(peer.Config{
		RepoPath:   t.TempDir(),
		EnableMDNS: true,
	})
	require.NoError(t, err)
	fin.Add(p1)

	p2, err := peer.New(peer.Config{
		RepoPath:   t.TempDir(),
		EnableMDNS: true,
	})
	require.NoError(t, err)
	fin.Add(p2)

	eventHandler := func(from core.ID, topic string, msg []byte) {
		t.Logf("%s event: %s %s", topic, from, msg)
	}
	messageHandler := func(from core.ID, topic string, msg []byte) ([]byte, error) {
		t.Logf("%s message: %s %s", topic, from, msg)
		return []byte("pong"), nil
	}

	t1, err := p1.NewTopic(context.Background(), "topic", true)
	require.NoError(t, err)
	t1.SetEventHandler(eventHandler)
	t1.SetMessageHandler(messageHandler)
	fin.Add(t1)

	t2, err := p2.NewTopic(context.Background(), "topic", true)
	require.NoError(t, err)
	t2.SetEventHandler(eventHandler)
	t2.SetMessageHandler(messageHandler)
	fin.Add(t2)

	time.Sleep(time.Second) // wait for mdns discovery

	// peer1 requests "pong" from peer2
	rc1, err := t1.Publish(context.Background(), []byte("ping"))
	require.NoError(t, err)
	r1 := <-rc1
	require.NotNil(t, r1)
	require.NoError(t, r1.Err)
	assert.Equal(t, "pong", string(r1.Data))
	assert.NotEmpty(t, r1.ID)
	assert.Equal(t, p2.Host().ID().String(), r1.From.String())

	// peer2 requests "pong" from peer1
	rc2, err := t2.Publish(context.Background(), []byte("ping"))
	require.NoError(t, err)
	r2 := <-rc2
	require.NotNil(t, r2)
	require.NoError(t, r2.Err)
	assert.Equal(t, "pong", string(r2.Data))
	assert.NotEmpty(t, r2.ID)
	assert.Equal(t, p1.Host().ID().String(), r2.From.String())

	require.NoError(t, fin.Cleanup(nil))
}

func TestMultiPingPong(t *testing.T) {
	fin := finalizer.NewFinalizer()

	p1, err := peer.New(peer.Config{
		RepoPath:   t.TempDir(),
		EnableMDNS: true,
	})
	require.NoError(t, err)
	fin.Add(p1)

	p2, err := peer.New(peer.Config{
		RepoPath:   t.TempDir(),
		EnableMDNS: true,
	})
	require.NoError(t, err)
	fin.Add(p2)

	p3, err := peer.New(peer.Config{
		RepoPath:   t.TempDir(),
		EnableMDNS: true,
	})
	require.NoError(t, err)
	fin.Add(p3)

	eventHandler := func(from core.ID, topic string, msg []byte) {
		t.Logf("%s event: %s %s", topic, from, msg)
	}
	messageHandler := func(from core.ID, topic string, msg []byte) ([]byte, error) {
		t.Logf("%s message: %s %s", topic, from, msg)
		return []byte("pong"), nil
	}

	t1, err := p1.NewTopic(context.Background(), "topic", false)
	require.NoError(t, err)
	t1.SetEventHandler(eventHandler)
	t1.SetMessageHandler(messageHandler)
	fin.Add(t1)

	t2, err := p2.NewTopic(context.Background(), "topic", true)
	require.NoError(t, err)
	t2.SetEventHandler(eventHandler)
	t2.SetMessageHandler(messageHandler)
	fin.Add(t2)

	t3, err := p3.NewTopic(context.Background(), "topic", true)
	require.NoError(t, err)
	t3.SetEventHandler(eventHandler)
	t3.SetMessageHandler(messageHandler)
	fin.Add(t3)

	time.Sleep(time.Second) // wait for mdns discovery

	// peer1 requests "pong" from peer2 and peer3
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rc, err := t1.Publish(ctx, []byte("ping"), rpc.WithMultiResponse(true))
	require.NoError(t, err)
	var pongs []struct{}
	for r := range rc {
		require.NotNil(t, r)
		require.NoError(t, r.Err)
		assert.Equal(t, "pong", string(r.Data))
		assert.NotEmpty(t, r.ID)
		pongs = append(pongs, struct{}{})
	}
	assert.Len(t, pongs, 2)

	require.NoError(t, fin.Cleanup(nil))
}

func setLogLevels(systems map[string]logging.LogLevel) error {
	for sys, level := range systems {
		l := zapcore.Level(level)
		if sys == "*" {
			for _, s := range logging.GetSubsystems() {
				if err := logging.SetLogLevel(s, l.CapitalString()); err != nil {
					return err
				}
			}
		}
		if err := logging.SetLogLevel(sys, l.CapitalString()); err != nil {
			return err
		}
	}
	return nil
}
