package v1

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/cryptex"
	"github.com/ieee0824/cryptex/kms"
	"github.com/jobtalk/pnzr/subcmd/vault/edit/util"
	"os"
)

type Editor struct {
	c          *cryptex.Cryptex
	editorName string
}

func New(s *session.Session, keyID string) *Editor {
	return &Editor{
		cryptex.New(kms.New(s).SetKey(keyID)),
		util.GetEditor(),
	}
}

func (e *Editor) Edit(fileName string) error {
	cryptex.SetEditor(e.editorName)
	var container = &cryptex.Container{}

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	if err := json.NewDecoder(f).Decode(container); err != nil {
		return err
	}
	result, err := e.c.Edit(container)
	if err != nil {
		return err
	}

	w, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(result); err != nil {
		return err
	}

	return nil
}
