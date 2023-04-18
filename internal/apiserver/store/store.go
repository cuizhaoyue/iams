package store

var client Factory

type Factory interface {
	Users() UserStore
	Secrets() SecretStore
	Polices() PolicyStore
	PolicyAudit() PolicyAuditStore
	Close() error
}

// Client 返回一个Factory实例.
func Client() Factory {
	return client
}

// SetClient 设置Factory实例.
func SetClient(factory Factory) {
	client = factory
}
