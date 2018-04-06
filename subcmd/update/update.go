package update

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/jobtalk/pnzr/vars"
)

const (
	GITHUB_API = "https://api.github.com/repos/jobtalk/pnzr/tags"
)

var client = &http.Client{}

type commit struct {
	Sha string `json:"sha"`
	URL string `json:"url"`
}

type tag struct {
	Name       string `json:"name"`
	ZipballURL string `json:"zipball_url"`
	TarballURL string `json:"tarball_url"`
	Commit     commit `json:"commit"`
}

type Env struct {
	os   string
	arch string
}

type OsEnv interface {
	checkENV() bool
	detectPlatform() (string, error)
}

func (t tag) String() string {
	bin, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return ""
	}
	return string(bin)
}

func (e *Env) checkENV() bool {
	if e.os != "darwin" && e.os != "linux" {
		return false
	}
	if e.arch != "amd64" {
		return false
	}
	return true
}

func (e *Env) detectPlatform() (string, error) {
	if e.os == "darwin" {
		return "darwin-amd64", nil
	} else if e.os == "linux" {
		return "linux-amd64", nil
	}
	return "", fmt.Errorf("This is not %s", "darwin or linux")
}

func checkVersion(latestVar string) (int, string) {
	if vars.VERSION == latestVar {
		return 0, "this version is latest"
	}
	if latestVar == "" {
		return 255, "can not get latest versiont"
	}
	return -1, latestVar
}

type UpdateCommand struct{}

func (c *UpdateCommand) Run(args []string) int {
	resp, err := client.Get(GITHUB_API)
	if err != nil {
		log.Println(err)
		return 255
	}
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return 255
	}
	defer resp.Body.Close()

	tags := []tag{}
	if err := json.Unmarshal(bin, &tags); err != nil {
		log.Println(err)
		return 255
	}
	e := Env{
		runtime.GOOS,
		runtime.GOARCH,
	}
	if !e.checkENV() {
		fmt.Printf("Sorry, this architecture not supported (%s %s).\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println("Please try manual update.")
		return 255
	}

	platform, err := e.detectPlatform()
	if err != nil {
		panic(err)
	}

	exitStatus, binVer := checkVersion(tags[0].Name)
	if exitStatus != -1 {
		fmt.Println(binVer)
		return exitStatus
	}

	binURL := fmt.Sprintf("https://github.com/jobtalk/pnzr/releases/download/%s/pnzr-%s", binVer, platform)

	latestResp, err := client.Get(binURL)
	if err != nil {
		log.Println(err)
		return 255
	}
	latestBin, err := ioutil.ReadAll(latestResp.Body)
	if err != nil {
		log.Println(err)
		return 255
	}

	dir, err := os.Executable()
	if err != nil {
		log.Println(err)
		return 255
	}

	if err := ioutil.WriteFile(dir, latestBin, 0755); err != nil {
		fmt.Printf("Check parmission: %s\n", dir)
		return 255
	}

	fmt.Println("update successed")

	return 0
}

func (c *UpdateCommand) Synopsis() string {
	return "Update pnzr to the latest version."
}

func (c *UpdateCommand) Help() string {
	msg := "\n\n"
	msg += "    Usage:\n"
	msg += "        no option\n\n"
	msg += "    Description:\n"
	msg += "        update pnzr to the latest version.\n"
	msg += "\n"
	return msg
}
