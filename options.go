package rpc

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// PublishOptions defines options for Publish.
type PublishOptions struct {
	ignoreResponse bool
	multiResponse  bool
	republish      bool
	pubOpts        []pubsub.PubOpt
}

var defaultOptions = PublishOptions{
	ignoreResponse: false,
	multiResponse:  false,
	republish:      false,
}

// PublishOption defines a Publish option.
type PublishOption func(*PublishOptions) error

// WithIgnoreResponse indicates whether or not Publish will wait for a response(s) from the receiver(s).
// Default: disabled.
func WithIgnoreResponse(enable bool) PublishOption {
	return func(args *PublishOptions) error {
		args.ignoreResponse = enable
		return nil
	}
}

// WithMultiResponse indicates whether or not Publish will wait for multiple responses before returning.
// Default: disabled.
func WithMultiResponse(enable bool) PublishOption {
	return func(args *PublishOptions) error {
		args.multiResponse = enable
		return nil
	}
}

// WithRepublishing indicates whether or not Publish will continue republishing to newly joined peers as long
// as the context hasn't expired or is not canceled.
// Default: disabled.
func WithRepublishing(enable bool) PublishOption {
	return func(args *PublishOptions) error {
		args.republish = enable
		return nil
	}
}

// WithPubOpts sets native pubsub.PubOpt options.
func WithPubOpts(opts ...pubsub.PubOpt) PublishOption {
	return func(args *PublishOptions) error {
		args.pubOpts = opts
		return nil
	}
}
