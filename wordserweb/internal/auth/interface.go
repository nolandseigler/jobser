package auth

type KeyValStorer interface {
	Insert(key string, value string) error
	Delete(key string) error
	Get(key string) (string, bool)
}

type UserVerifier interface {
	Verify(username string, password string) error
}
