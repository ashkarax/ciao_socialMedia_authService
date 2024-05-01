package interface_hash_authSvc

type IhashPassword interface {
	HashPassword(password string) string
	CompairPassword(hashedPassword string, plainPassword string) error
}
