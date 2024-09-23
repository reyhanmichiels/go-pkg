package hash

type Interface interface {
	Bcrypt() bcryptInterface
	// TODO: implement hashing algorithm with argon2 and scrypt
}

type hash struct {
	bcrypt bcryptInterface
}

func Init() Interface {
	return &hash{
		bcrypt: initBcrypt(),
	}
}

func (h *hash) Bcrypt() bcryptInterface {
	return h.bcrypt
}
