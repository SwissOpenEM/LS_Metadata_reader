package configuration

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func Getconfig() []byte {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	configFilePath := filepath.Join(usr.HomeDir, ".config", "LS_reader.conf")

	_, err1 := os.Stat(configFilePath)
	if err1 != nil {
		Setconfig(configFilePath, usr)
	}
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return content
}

func Setconfig(path string, usr *user.User) {
	fmt.Println("We discovered you don't have your instrument config set - please do so now:")
	fmt.Println("What is your instruments spherical aberration (CS)?")
	reader := bufio.NewReader(os.Stdin)
	input1, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	fmt.Println("And what is the rotation or flipping that needs to be done when importing the gain reference to cryosparc?")
	reader2 := bufio.NewReader(os.Stdin)
	input2, err := reader2.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	configmap := make(map[string]string)
	configmap["CS"] = strings.TrimSpace(input1)
	configmap["Gainref_FlipRotate"] = strings.TrimSpace(input2)
	config, err := json.MarshalIndent(configmap, "", "    ")
	if err != nil {
		fmt.Println("Error generating config:", err)
		return
	}
	configfolder := filepath.Join(usr.HomeDir, ".config")
	_, errexist := os.Stat(configfolder)
	if os.IsNotExist(errexist) {
		os.Mkdir(configfolder, 0755)
	}
	err = os.WriteFile(path, config, 0644)
	if err != nil {
		fmt.Println("Error generating config:", err)
	}
	fmt.Println("Generated config at", path)
}
