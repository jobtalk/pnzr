package vedit

import (
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

func encrypt(keyID string, fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMS()
	_, err = kms.SetKeyID(keyID).Encrypt(bin)
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
	if kms == nil {
		return errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.SetKeyID(keyID).Decrypt()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, plainText, 0644)
}

type VaultEdit struct{}

func (c *VaultEdit) Help() string {
	var msg string
	msg += "usage: pnzr vault-edit [options ...]\n"
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

func (c *VaultEdit) Synopsis() string {
	return c.Help()
}

func getEditor() string {
	if e := os.Getenv("PNZR_EDITOR"); e != "" {
		return e
	}

	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}

	return "nano"
}

func (c *VaultEdit) Run(args []string) int {
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

	if err := decrypt(session, *kmsKeyID, *file); err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command(getEditor(), *file)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	if err := encrypt(session, *kmsKeyID, *file); err != nil {
		log.Fatalln(err)
	}

	return 0
}
