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
	fmt.Println("And what is the rotation or flipping that needs to be done when importing the gain reference to e.g cryosparc?")
	reader2 := bufio.NewReader(os.Stdin)
	input2, err := reader2.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	fmt.Println("If available what is the path where EPU mirrors your data output and dumps its metadata .csvs (typically this is on the microscope computer). If you dont know/ cant reacht that folder leave this empty.")
	reader3 := bufio.NewReader(os.Stdin)
	input3, err := reader3.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	configmap := make(map[string]string)
	configmap["CS"] = strings.TrimSpace(input1)
	configmap["Gainref_FlipRotate"] = strings.TrimSpace(input2)
	configmap["MPCPATH"] = strings.TrimSpace(input3)
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

func Changeconfig() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
	}
	configFilePath := filepath.Join(usr.HomeDir, ".config", "LS_reader.conf")
	_, err1 := os.Stat(configFilePath)
	if err1 != nil {
		Setconfig(configFilePath, usr)
		os.Exit(0)
	}

	fmt.Println("What is your instruments spherical aberration (CS)?")
	reader := bufio.NewReader(os.Stdin)
	input1, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	fmt.Println("And what is the rotation or flipping that needs to be done when importing the gain reference to e.g cryosparc?")
	reader2 := bufio.NewReader(os.Stdin)
	input2, err := reader2.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	fmt.Println("If available what is the path where EPU mirrors your data output and dumps its metadata .xmls (typically this is on the microscope computer). If you dont know/ cant reach that folder leave this empty. For optimal usage of this tool this is however required.")
	reader3 := bufio.NewReader(os.Stdin)
	input3, err := reader3.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	configmap := make(map[string]string)
	configmap["CS"] = strings.TrimSpace(input1)
	configmap["Gainref_FlipRotate"] = strings.TrimSpace(input2)
	configmap["MPCPATH"] = strings.TrimSpace(input3)
	config, err := json.MarshalIndent(configmap, "", "    ")
	if err != nil {
		fmt.Println("Error generating config:", err)
		return
	}
	err = os.WriteFile(configFilePath, config, 0644)
	if err != nil {
		fmt.Println("Error generating config:", err)
	}
	fmt.Println("Generated config at", configFilePath)
}
