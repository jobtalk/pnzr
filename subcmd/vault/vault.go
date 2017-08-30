package vault

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"

	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/lib"
	"os"
	"os/exec"
)

func getEditor() string {
	if e := os.Getenv("PNZR_EDITOR"); e != "" {
		return e
	}

	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}

	return "nano"
}

type mode struct {
	encrypt *bool
	decrypt *bool
	edit    *bool
	view    *bool
}

func (m *mode)checkMultiFlagSet() bool {
	return (*m.encrypt && *m.decrypt) ||
		(*m.encrypt && *m.edit) ||
		(*m.encrypt && *m.view) ||
		(*m.decrypt && *m.edit) ||
		(*m.decrypt && *m.view) ||
		(*m.edit && *m.view)
}

type VaultCommand struct {
	sess           *session.Session
	kmsKeyID       *string
	vaultMode      *mode
	file           *string
	profile        *string
	region         *string
	awsAccessKeyID *string
	awsSecretKeyID *string
}

func (v *VaultCommand) parseArgs(args []string) {
	var (
		flagSet = new(flag.FlagSet)
		f       *string
	)
	v.vaultMode = new(mode)

	v.kmsKeyID = flagSet.String("key_id", getenv.String("KMS_KEY_ID"), "Amazon KMS key ID")
	v.vaultMode.encrypt = flagSet.Bool("encrypt", getenv.Bool("ENCRYPT", false), "encrypt mode")
	v.vaultMode.decrypt = flagSet.Bool("decrypt", getenv.Bool("DECRYPT", false), "decrypt mode")
	v.vaultMode.view = flagSet.Bool("view", false, "view mode")
	v.vaultMode.edit = flagSet.Bool("edit", false, "edit mode")
	v.profile = flagSet.String("profile", getenv.String("AWS_PROFILE_NAME", "default"), "aws credentials profile name")
	v.region = flagSet.String("region", getenv.String("AWS_REGION", "ap-northeast-1"), "aws region")
	v.awsAccessKeyID = flagSet.String("aws-access-key-id", getenv.String("AWS_ACCESS_KEY_ID"), "aws access key id")
	v.awsSecretKeyID = flagSet.String("aws-secret-key-id", getenv.String("AWS_SECRET_KEY_ID"), "aws secret key id")
	v.file = flagSet.String("file", "", "target file")
	f = flagSet.String("f", "", "target file")

	if err := flagSet.Parse(args); err != nil {
		panic(err)
	}

	if *f == "" && *v.file == "" && len(flagSet.Args()) != 0 {
		targetName := flagSet.Args()[0]
		v.file = &targetName
	}

	if *v.file == "" {
		v.file = f
	}

	var awsConfig = aws.Config{}

	if *v.awsAccessKeyID != "" && *v.awsSecretKeyID != "" && *v.profile == "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(*v.awsAccessKeyID, *v.awsSecretKeyID, "")
		awsConfig.Region = v.region
	}

	v.sess = session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 *v.profile,
		Config:                  awsConfig,
	}))
}

func (v *VaultCommand) encrypt(keyID string, fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMS(v.sess)
	_, err = kms.SetKeyID(keyID).Encrypt(bin)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, []byte(kms.String()), 0644)
}

func (v *VaultCommand) decrypt(keyID string, fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMSFromBinary(bin, v.sess)
	if kms == nil {
		return errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.SetKeyID(keyID).Decrypt()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, plainText, 0644)
}

func (v *VaultCommand) decryptTemporary(keyID string, fileName string) ([]byte, error) {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	kms := lib.NewKMSFromBinary(bin, v.sess)
	if kms == nil {
		return nil, errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.SetKeyID(keyID).Decrypt()
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

func (c *VaultCommand) Help() string {
	var msg string
	msg += "usage: pnzr vault [options ...]\n"
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

func (v *VaultCommand) Run(args []string) int {
	v.parseArgs(args)

	if v.vaultMode.checkMultiFlagSet() {
		panic("Multiple vault options are selected.")
	} else if *v.vaultMode.encrypt {
		if err := v.encrypt(*v.kmsKeyID, *v.file); err != nil {
			panic(err)
		}
	} else if *v.vaultMode.decrypt {
		if err := v.decrypt(*v.kmsKeyID, *v.file); err != nil {
			panic(err)
		}
	} else if *v.vaultMode.edit {
		if err := v.decrypt(*v.kmsKeyID, *v.file); err != nil {
			panic(err)
		}
		defer func() {
			if err := v.encrypt(*v.kmsKeyID, *v.file); err != nil {
				panic(err)
			}
		}()
		cmd := exec.Command(getEditor(), *v.file)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			panic(err)
		}
	} else if *v.vaultMode.view {
		plain, err := v.decryptTemporary(*v.kmsKeyID, *v.file)
		if err != nil {
			panic(err)
		}

		cmd := exec.Command("less")
		cmd.Stdin = bytes.NewReader(plain)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			panic(err)
		}
	} else {
		panic("Vault mode is not selected.")
	}

	return 0
}

func (c *VaultCommand) Synopsis() string {
	return c.Help()
}
