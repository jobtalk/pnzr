package vault

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/jobtalk/thor/lib"
)

var flagSet = &flag.FlagSet{}

var (
	kmsKeyID    *string
	encryptFlag *bool
	decryptFlag *bool
	file        *string
	f           *string
)

func init() {
	kmsKeyID = flagSet.String("key_id", "", "Amazon KMS key ID")
	encryptFlag = flagSet.Bool("encrypt", false, "encrypt mode")
	decryptFlag = flagSet.Bool("decrypt", false, "decrypt mode")

	file = flagSet.String("file", "", "target file")
	f = flagSet.String("f", "", "target file")
}

func encrypt(keyID string, fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMS()
	_, err = kms.SetKeyID(keyID).SetRegion("ap-northeast-1").Encrypt(bin)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, []byte(kms.String()), 0644)
}

func decrypt(keyID string, fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMSFromBinary(bin)
	plainText, err := kms.SetKeyID(keyID).SetRegion("ap-northeast-1").Decrypt()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, plainText, 0644)
}

type Vault struct{}

func (c *Vault) Help() string {
	var msg string
	msg += "usage: thor vault [options ...]\n"
	msg += "options:\n"
	msg += "    -key_id\n"
	msg += "        set kms key id\n"
	msg += "    -encrypt\n"
	msg += "        use encrypt mode\n"
	msg += "    -decrypt\n"
	msg += "        use decrypt mode\n"
	msg += "    -file\n"
	msg += "        setting target file\n"
	msg += "    -f"
	msg += "        setting target file\n"
	return msg
}

func (c *Vault) Run(args []string) int {
	if err := flagSet.Parse(args); err != nil {
		log.Fatalln(err)
	}
	if *file == "" {
		file = f
	}
	if *encryptFlag == *decryptFlag {
		log.Fatalln("Choose whether to execute A or B.")
	}
	if *decryptFlag {
		err := decrypt(*kmsKeyID, *file)
		if err != nil {
			log.Fatalln(err)
		}
	} else if *encryptFlag {
		err := encrypt(*kmsKeyID, *file)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return 0
}

func (c *Vault) Synopsis() string {
	return c.Help()
}
