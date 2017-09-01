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

func (t tag) String() string {
	bin, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return ""
	}
	return string(bin)
}

func checkENV() bool {
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		return false
	}
	if runtime.GOARCH != "amd64" {
		return false
	}
	return true
}

type UpdateCommand struct{}

func (c *UpdateCommand) Run(args []string) int {
	var platform string
	tags := []tag{}
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
	if err := json.Unmarshal(bin, &tags); err != nil {
		log.Println(err)
		return 255
	}

	if !checkENV() {
		fmt.Printf("Sorry, this architecture not supported (%s %s).\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println("Please try manual update.")
		return 255
	}

	if runtime.GOOS == "darwin" {
		platform = "darwin-amd64"
	} else if runtime.GOOS == "linux" {
		platform = "linux-amd64"
	}
	latest := tags[0].Name
	if vars.VERSION == latest {
		fmt.Println("this version is latest")
		return 0
	}
	if latest == "" {
		fmt.Println("can not get latest version")
		return 255
	}
	binaryURL := fmt.Sprintf("https://github.com/jobtalk/pnzr/releases/download/%s/pnzr-%s", latest, platform)

	dir, err := os.Executable()
	if err != nil {
		log.Println(err)
		return 255
	}

	latestResp, err := client.Get(binaryURL)
	if err != nil {
		log.Println(err)
		return 255
	}
	latestBin, err := ioutil.ReadAll(latestResp.Body)
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
