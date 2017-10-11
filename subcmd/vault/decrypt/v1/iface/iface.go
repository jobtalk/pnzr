package v1_api

type API interface {
	Decrypt(string) error
}
