package v1

import (
	"github.com/ieee0824/cryptex"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/cryptex/kms"
	"os"
	"github.com/gin-gonic/gin/json"
	"github.com/jobtalk/pnzr/subcmd/vault/edit/util"
)

type Editor struct {
	c *cryptex.Cryptex
}

func New(s *session.Session, keyID string) *Editor {
	return &Editor{
		cryptex.New(kms.New(s).SetKey(keyID)),
	}
}


func (e *Editor)Edit(fileName string) error {
	cryptex.SetEditor(util.GetEditor())
	var container = &cryptex.Container{}

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	if err := json.NewDecoder(f).Decode(container); err != nil {
		return err
	}

	if _, err := e.c.Edit(container); err != nil {
		return err
	}
	return nil
}
