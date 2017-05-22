package getenv

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var TEST_DIR = os.Getenv("GOPATH") + "/src/github.com/jobtalk/pnzr/test"

func init() {
	err := godotenv.Load(fmt.Sprintf("%s/getenv/.env.test", TEST_DIR))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestConvertStringToBoolean(t *testing.T) {
	if convertStringToBoolean("") {
		t.Fatalf("値が不正: %v", convertStringToBoolean(""))
	} else if convertStringToBoolean("false") {
		t.Fatalf("値が不正: %v", convertStringToBoolean("false"))
	} else if convertStringToBoolean("False") {
		t.Fatalf("値が不正: %v", convertStringToBoolean("False"))
	} else if convertStringToBoolean("FALSE") {
		t.Fatalf("値が不正: %v", convertStringToBoolean("FALSE"))
	} else if convertStringToBoolean("0") {
		t.Fatalf("値が不正: %v", convertStringToBoolean("0"))
	} else if convertStringToBoolean("f") {
		t.Fatalf("値が不正: %v", convertStringToBoolean("f"))
	}

	if !convertStringToBoolean("1") {
		t.Fatalf("値が不正: %v", convertStringToBoolean("1"))
	} else if !convertStringToBoolean("t") {
		t.Fatalf("値が不正: %v", convertStringToBoolean("t"))
	} else if !convertStringToBoolean("true") {
		t.Fatalf("値が不正: %v", convertStringToBoolean("true"))
	} else if !convertStringToBoolean("True") {
		t.Fatalf("値が不正: %v", convertStringToBoolean("True"))
	}
}

func TestBool(t *testing.T) {
	{
		t.Log("環境変数が存在しない時デフォルト値を見るテスト")
		if !Bool("hoge", true) {
			t.Fatalf("値が不正: %v", Bool("hoge", true))
		} else if Bool("hoge", false) {
			t.Fatalf("値が不正: %v", Bool("hoge", false))
		}
	}

	{
		t.Log("envに設定されている時")
		{
			t.Log("trueの時")
			if !Bool("BOOL_0") {
				t.Log("BOOL_0")
				t.Fatalf("値が不正: %v", Bool("BOOL_0"))
			} else if !Bool("BOOL_1") {
				t.Log("BOOL_1")
				t.Fatalf("値が不正: %v", Bool("BOOL_1"))
			} else if !Bool("BOOL_2") {
				t.Log("BOOL_2")
				t.Fatalf("値が不正: %v", Bool("BOOL_2"))
			} else if !Bool("BOOL_3") {
				t.Log("BOOL_3")
				t.Fatalf("値が不正: %v", Bool("BOOL_3"))
			}
		}

		{
			t.Log("falseの時")
			if Bool("BOOL_4") {
				t.Log("BOOL_4")
				t.Fatalf("値が不正: %v", Bool("BOOL_4"))
			} else if Bool("BOOL_5") {
				t.Log("BOOL_5")
				t.Fatalf("値が不正: %v", Bool("BOOL_5"))
			} else if Bool("BOOL_6") {
				t.Log("BOOL_6")
				t.Fatalf("値が不正: %v", Bool("BOOL_6"))
			} else if Bool("BOOL_7") {
				t.Log("BOOL_7")
				t.Fatalf("値が不正: %v", Bool("BOOL_7"))
			} else if Bool("BOOL_8") {
				t.Log("BOOL_8")
				t.Fatalf("値が不正: %v", Bool("BOOL_8"))
			} else if Bool("hoge") {
				t.Log("hoge")
				t.Fatalf("値が不正: %v", Bool("BOOL_8"))
			}
		}

	}
}

func TestString(t *testing.T) {
	if String("") != "" {
		t.Fatalf("値が不正: %v", String(""))
	} else if String("", "hoge") != "hoge" {
		t.Fatalf("値が不正: %v", String("", "hoge"))
	}

	if String("PANZER") != "VOR" {
		t.Fatalf("値が不正: %v", String("PANZER"))
	}
}
