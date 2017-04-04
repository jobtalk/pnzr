package lib

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type KMS struct {
	keyID     *string
	awsConfig *aws.Config
	Type      *string `json:"type"`
	Chipher   []byte  `json:"chipher"`
}

func NewKMS() *KMS {
	return &KMS{
		awsConfig: &aws.Config{},
		Type:      aws.String("kms"),
	}
}

func (k *KMS) Encrypt(plainText []byte) ([]byte, error) {
	svc := kms.New(session.New(), k.awsConfig)
	params := &kms.EncryptInput{
		KeyId:     k.keyID,
		Plaintext: plainText,
	}
	k.Chipher = nil
	resp, err := svc.Encrypt(params)
	if err != nil {
		return nil, err
	}

	return resp.CiphertextBlob, nil
}

func (k *KMS) Decrypt(cipherText []byte) ([]byte, error) {
	svc := kms.New(session.New(), k.awsConfig)
	params := &kms.DecryptInput{
		CiphertextBlob: cipherText,
	}
	resp, err := svc.Decrypt(params)
	if err != nil {
		return nil, err
	}
	k.Chipher = resp.Plaintext
	return resp.Plaintext, nil
}

func (k *KMS) SetKeyID(keyID string) *KMS {
	k.keyID = &keyID
	return k
}

func (k *KMS) SetRegion(region string) *KMS {
	k.awsConfig.Region = &region
	return k
}

func (k *KMS) SetAWSConfig(awsConfig *aws.Config) *KMS {
	k.awsConfig = awsConfig
	return k
}

func (k *KMS) String() string {
	bin, err := json.Marshal(k)
	if err != nil {
		return ""
	}
	return string(bin)
}
