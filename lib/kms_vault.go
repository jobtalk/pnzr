package lib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type KMS struct {
	KeyID     *string
	awsConfig *aws.Config
}

func NewKMS() *KMS {
	return &KMS{
		awsConfig: &aws.Config{},
	}
}

func (k *KMS) Encrypt(plainText []byte) ([]byte, error) {
	svc := kms.New(session.New(), k.awsConfig)
	params := &kms.EncryptInput{
		KeyId:     k.KeyID,
		Plaintext: plainText,
	}
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
	return resp.Plaintext, nil
}

func (k *KMS) SetKeyID(keyID string) *KMS {
	k.KeyID = &keyID
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
