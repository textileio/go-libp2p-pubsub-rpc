package rpc

import (
	"errors"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// PublishOptions defines options for Publish.
type PublishOptions struct {
	ignoreResponse bool
	multiResponse  bool
	numRetries     int
	retryTimeout   time.Duration
	retryWait      time.Duration
	retryBackoff   float32
	pubOpts        []pubsub.PubOpt
}

var defaultOptions = PublishOptions{
	ignoreResponse: false,
	multiResponse:  false,
	numRetries:     0,
	retryTimeout:   time.Second * 5,
	retryWait:      time.Second,
	retryBackoff:   1.5,
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

// WithRetries enables publish retries.
// Retries will be canceled if the context passed to Publish expires.
// Retries are only done if the error is ErrResponseNotReceived.
// When used with WithMultiResponse, retries are only done if there are zero responses.
// Default: disabled.
func WithRetries(
	numRetries int,
	retryTimeout time.Duration,
	retryWait time.Duration,
	retryBackoff float32,
) PublishOption {
	return func(args *PublishOptions) error {
		if numRetries < 0 {
			return errors.New("numRetries must be greater than or equal to zero")
		}
		if retryTimeout <= 0 {
			return errors.New("retryTimeout must be greater than zero")
		}
		if retryWait <= 0 {
			return errors.New("retryTimeout must be greater than zero")
		}
		if retryBackoff < 1 {
			return errors.New("retryBackoff must be greater than or equal to one")
		}
		args.numRetries = numRetries
		args.retryTimeout = retryTimeout
		args.retryWait = retryWait
		args.retryBackoff = retryBackoff
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
