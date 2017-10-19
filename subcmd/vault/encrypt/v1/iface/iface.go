package v1_api

type API interface {
	Encrypt(keyID, fileName string) error
}
