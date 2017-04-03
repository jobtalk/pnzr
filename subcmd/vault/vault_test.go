package vault

import "testing"

func TestVaultParamValidate(t *testing.T) {
	var pass = "hoge"
	var path = "huga"
	var emptyPass = ""
	var emptyPath = ""
	t.Log("<--validate nil param test-->")
	nilParam := &vaultParam{}
	if err := nilParam.validate(); err == nil {
		t.Errorf("ng: param is nil. but error is nil")
	} else {
		t.Logf("ok: error is %v", err)
	}

	nilParam.Pass = &pass
	nilParam.Path = nil
	if err := nilParam.validate(); err == nil {
		t.Errorf("ng: param is nil. but error is nil")
	} else {
		t.Logf("ok: error is %v", err)
	}

	nilParam.Pass = nil
	nilParam.Path = &path
	if err := nilParam.validate(); err == nil {
		t.Errorf("ng: param is nil. but error is nil")
	} else {
		t.Logf("ok: error is %v", err)
	}

	t.Log("<--validate empty param test-->")
	emptyParam := &vaultParam{
		&emptyPass,
		&path,
	}
	if err := emptyParam.validate(); err == nil {
		t.Errorf("ng: param is empty. but error is nil")
	} else {
		t.Logf("ok: error is %v", err)
	}

	emptyParam.Pass = &pass
	emptyParam.Path = &emptyPath
	if err := emptyParam.validate(); err == nil {
		t.Errorf("ng: param is empty. but error is nil")
	} else {
		t.Logf("ok: error is %v", err)
	}
}
