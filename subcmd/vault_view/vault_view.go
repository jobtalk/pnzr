package vview

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/lib"
)

var flagSet = &flag.FlagSet{}

var (
	kmsKeyID *string
	file     *string
	f        *string
)

func init() {
	kmsKeyID = flagSet.String("key_id", getenv.String("KMS_KEY_ID"), "Amazon KMS key ID")

	file = flagSet.String("file", "", "target file")
	f = flagSet.String("f", "", "target file")
}

func decrypt(keyID string, fileName string) ([]byte, error) {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	kms := lib.NewKMSFromBinary(bin)
	if kms == nil {
		return nil, errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.SetKeyID(keyID).Decrypt()
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

type VaultView struct{}

func (c *VaultView) Help() string {
	var msg string
	msg += "usage: pnzr vault-view [options ...]\n"
	msg += "options:\n"
	msg += "    -key_id\n"
	msg += "        set kms key id\n"
	msg += "    -file\n"
	msg += "        setting target file\n"
	msg += "    -f"
	msg += "        setting target file\n"
	msg += "    -profile\n"
	msg += "        aws credential name\n"
	msg += "    -region\n"
	msg += "        aws region name\n"
	msg += "    -aws-access-key-id\n"
	msg += "        setting aws access key id\n"
	msg += "    -aws-secret-key-id\n"
	msg += "        setting aws secret key id\n"
	msg += "===================================================\n"
	return msg
}

func (c *VaultView) Synopsis() string {
	return c.Help()
}

func (c *VaultView) Run(args []string) int {
	session, args := lib.GetSession(args)
	if err := flagSet.Parse(args); err != nil {
		log.Fatalln(err)
	}

	if *f == "" && *file == "" && len(flagSet.Args()) != 0 {
		targetName := flagSet.Args()[0]
		file = &targetName
	}

	if *file == "" {
		file = f
	}

	plain, err := decrypt(session, *kmsKeyID, *file)
	if err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command("less")
	cmd.Stdin = bytes.NewReader(plain)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	return 0
}
