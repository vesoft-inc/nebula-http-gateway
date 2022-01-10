package nebula

import (
	"crypto/tls"
	"math"
	"time"
)

const (
	DefaultTimeout        = time.Duration(0)
	DefaultBufferSize     = 128 << 10
	DefaultFrameMaxLength = math.MaxUint32
)

type (
	Options struct {
		version      Version
		log          Logger
		graph        socketOptions
		meta         socketOptions
		storageAdmin socketOptions
	}

	socketOptions struct {
		timeout        time.Duration
		bufferSize     int
		frameMaxLength uint32
		tlsConfig      *tls.Config
	}

	Option func(o *Options)
)

func WithVersion(version Version) Option {
	return func(o *Options) {
		o.version = version
	}
}

func WithLogger(log Logger) Option {
	return func(o *Options) {
		o.log = log
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		WithGraphTimeout(timeout)
		WithMetaTimeout(timeout)
		WithStorageTimeout(timeout)
	}
}

func WithGraphTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.graph.timeout = timeout
	}
}

func WithMetaTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.meta.timeout = timeout
	}
}

func WithStorageTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.storageAdmin.timeout = timeout
	}
}

func WithBufferSize(bufferSize int) Option {
	return func(o *Options) {
		WithMetaBufferSize(bufferSize)
		WithGraphBufferSize(bufferSize)
		WithStorageBufferSize(bufferSize)
	}
}

func WithGraphBufferSize(bufferSize int) Option {
	return func(o *Options) {
		o.graph.bufferSize = bufferSize
	}
}

func WithMetaBufferSize(bufferSize int) Option {
	return func(o *Options) {
		o.meta.bufferSize = bufferSize
	}
}

func WithStorageBufferSize(bufferSize int) Option {
	return func(o *Options) {
		o.storageAdmin.bufferSize = bufferSize
	}
}

func WithFrameMaxLength(frameMaxLength uint32) Option {
	return func(o *Options) {
		WithGraphFrameMaxLength(frameMaxLength)
		WithMetaFrameMaxLength(frameMaxLength)
		WithStorageFrameMaxLength(frameMaxLength)
	}
}

func WithGraphFrameMaxLength(frameMaxLength uint32) Option {
	return func(o *Options) {
		o.graph.frameMaxLength = frameMaxLength
	}
}

func WithMetaFrameMaxLength(frameMaxLength uint32) Option {
	return func(o *Options) {
		o.meta.frameMaxLength = frameMaxLength
	}
}

func WithStorageFrameMaxLength(frameMaxLength uint32) Option {
	return func(o *Options) {
		o.storageAdmin.frameMaxLength = frameMaxLength
	}
}

func WithTLS(tlsConfig *tls.Config) Option {
	return func(o *Options) {
		WithGraphTLS(tlsConfig)
		WithMetaTLS(tlsConfig)
		WithStorageTLS(tlsConfig)
	}
}

func WithGraphTLS(tlsConfig *tls.Config) Option {
	return func(o *Options) {
		o.graph.tlsConfig = tlsConfig
	}
}

func WithMetaTLS(tlsConfig *tls.Config) Option {
	return func(o *Options) {
		o.meta.tlsConfig = tlsConfig
	}
}

func WithStorageTLS(tlsConfig *tls.Config) Option {
	return func(o *Options) {
		o.storageAdmin.tlsConfig = tlsConfig
	}
}

func (o *Options) complete() {
	defaultOpts := defaultOptions()

	if o.log == nil {
		o.log = defaultOpts.log
	}
	o.graph.complete()
	o.meta.complete()
	o.storageAdmin.complete()
}

func (o *Options) validate() error {
	return nil
}

func (o *socketOptions) complete() {
	defaultOpts := defaultSocketOptions()
	if o.timeout < 0 {
		o.timeout = defaultOpts.timeout
	}
	if o.bufferSize <= 0 {
		o.bufferSize = defaultOpts.bufferSize
	}
	if o.frameMaxLength <= 0 {
		o.frameMaxLength = defaultOpts.frameMaxLength
	}
}

func defaultOptions() Options {
	return Options{
		version:      VersionAuto,
		log:          noOpLogger{},
		graph:        defaultSocketOptions(),
		meta:         defaultSocketOptions(),
		storageAdmin: defaultSocketOptions(),
	}
}

func defaultSocketOptions() socketOptions {
	return socketOptions{
		timeout:        DefaultTimeout,
		bufferSize:     DefaultBufferSize,
		frameMaxLength: DefaultFrameMaxLength,
		tlsConfig:      nil,
	}
}
