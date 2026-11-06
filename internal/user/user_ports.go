package user

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(plain, hashed string) bool
}

type Dependencies struct {
	Passwords PasswordHasher
}
