# go-libp2p-pubsub-rpc

[![Made by Textile](https://img.shields.io/badge/made%20by-Textile-informational.svg)](https://textile.io)
[![Chat on Slack](https://img.shields.io/badge/slack-slack.textile.io-informational.svg)](https://slack.textile.io)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg)](https://github.com/RichardLitt/standard-readme)

> RPC over libp2p pubsub with error handling

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)
- [Contributing](#contributing)
- [Changelog](#changelog)
- [License](#license)

## Background

`go-libp2p-pubsub-rpc` is an extension to [go-libp2p-pubsub](https://github.com/libp2p/go-libp2p-pubsub) that provides RPC-like functionality:

- _Request/Response_ pattern with a peer
- _Request/Multi-Response_ pattern with multiple peers

## Install

```bash
go get github.com/textileio/go-libp2p-pubsub-rpc
```

## Usage

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	core "github.com/libp2p/go-libp2p-core/peer"
	rpc "github.com/textileio/go-libp2p-pubsub-rpc"
	"github.com/textileio/go-libp2p-pubsub-rpc/peer"
)

func main() {

	//
	// PING PONG WITH TWO PEERS
	//

	// create a peer
	p1, _ := peer.New(peer.Config{
		RepoPath:   "/tmp/repo1",
		EnableMDNS: true, // can be used when peers are on the same network
	})

	// create another peer
	p2, _ := peer.New(peer.Config{
		RepoPath:   "/tmp/repo2",
		EnableMDNS: true,
	})

	eventHandler := func(from core.ID, topic string, msg []byte) {
		fmt.Printf("%s event: %s %s\n", topic, from, msg)
	}
	messageHandler := func(from core.ID, topic string, msg []byte) ([]byte, error) {
		fmt.Printf("%s message: %s %s\n", topic, from, msg)
		if string(msg) == "ping" {
			return []byte("pong"), nil
		} else {
			return nil, errors.New("invalid request")
		}
	}

	t1, _ := p1.NewTopic(context.Background(), "mytopic", true) // no need to subscribe if only publishing
	t1.SetEventHandler(eventHandler)                            // event handler reports topic membership events
	t1.SetMessageHandler(messageHandler)                        // message handle is any func that returns a response and error

	t2, _ := p2.NewTopic(context.Background(), "mytopic", true)
	t2.SetEventHandler(eventHandler)
	t2.SetMessageHandler(messageHandler) // using same message handler as peer1, but this could be anything

	time.Sleep(time.Second) // wait for mdns discovery

	// peer1 requests "pong" from peer2
	rc1, _ := t1.Publish(context.Background(), []byte("ping"))
	r1 := <-rc1
	// check r1.Err
	fmt.Printf("peer1 received \"%s\" from %s\n", r1.Data, r1.From)

	// peer2 requests "pong" from peer1
	rc2, _ := t2.Publish(context.Background(), []byte("ping"))
	r2 := <-rc2
	// check r2.Err
	fmt.Printf("peer2 received \"%s\" from %s\n", r2.Data, r2.From)

	// peers can respond with an error
	rc3, _ := t2.Publish(context.Background(), []byte("not a ping"))
	r3 := <-rc3
	fmt.Printf("peer2 received error \"%s\" from %s\n", r3.Err, r3.From)

	//
	// PING PONG WITH MULTIPLE PEERS
	//

	// create another peer
	p3, _ := peer.New(peer.Config{
		RepoPath:   "/tmp/repo3",
		EnableMDNS: true,
	})
	t3, _ := p3.NewTopic(context.Background(), "mytopic", true)
	t3.SetEventHandler(eventHandler)
	t3.SetMessageHandler(messageHandler)

	time.Sleep(time.Second) // wait for mdns discovery

	// peer1 requests "pong" from peer2 and peer3
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rc, _ := t1.Publish(ctx, []byte("ping"), rpc.WithMultiResponse(true))
	for r := range rc {
		// check r.Err
		fmt.Printf("peer1 received \"%s\" from %s\n", r.Data, r.From)
	}
}
```

## Contributing

Pull requests and bug reports are very welcome ❤️

This repository falls under the Textile [Code of Conduct](./CODE_OF_CONDUCT.md).

Feel free to get in touch by:
-   [Opening an issue](https://github.com/textileio/go-libp2p-pubsub-rpc/issues/new)
-   Joining the [public Slack channel](https://slack.textile.io/)
-   Sending an email to contact@textile.io

## Changelog

A changelog is published along with each [release](https://github.com/textileio/go-libp2p-pubsub-rpc/releases).

## License

[MIT](LICENSE)
