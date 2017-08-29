package lib

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

type KMS struct {
	keyID  *string
	svc    kmsiface.KMSAPI
	Type   *string `json:"type"`
	Cipher []byte  `json:"cipher"`
}

func NewKMS(sess *session.Session) *KMS {
	return &KMS{
		svc:  kms.New(sess),
		Type: aws.String("kms"),
	}
}

func NewKMSFromBinary(bin []byte, sess *session.Session) *KMS {
	var ret = KMS{}
	err := json.Unmarshal(bin, &ret)
	if err != nil {
		return nil
	}
	ret.svc = kms.New(sess)
	return &ret
}

func (k *KMS) Encrypt(plainText []byte) ([]byte, error) {
	params := &kms.EncryptInput{
		KeyId:     k.keyID,
		Plaintext: plainText,
	}
	resp, err := k.svc.Encrypt(params)
	if err != nil {
		return nil, err
	}

	k.Cipher = resp.CiphertextBlob
	return resp.CiphertextBlob, nil
}

func (k *KMS) Decrypt() ([]byte, error) {
	params := &kms.DecryptInput{
		CiphertextBlob: k.Cipher,
	}
	resp, err := k.svc.Decrypt(params)
	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil
}

func (k *KMS) SetKeyID(keyID string) *KMS {
	k.keyID = &keyID
	return k
}

func (k *KMS) String() string {
	bin, err := json.Marshal(k)
	if err != nil {
		return ""
	}
	return string(bin)
}
