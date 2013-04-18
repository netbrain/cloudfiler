package conf

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io"
	"log"
	"os"
	"reflect"
)

var Config *config
var userHome, configHome, configFile string

type config struct {
	ConfigFilePath               string            `json:"-"`
	ApplicationHome              string            `json:"applicationHome"`
	ServerAddr                   string            `json:"serverAddr"`
	CookieStoreAuthenticationKey []byte            `json:"cookieStoreAuthenticationKey"`
	CookieStoreEncryptionKey     []byte            `json:"cookieStoreEncryptionKey"`
	Repository                   map[string]string `json:"repository"`
}

func init() {
	userHome = os.Getenv("HOME")
	configHome = userHome + "/.cloudfiler"
	configFile = configHome + "/config"

	wd, _ := os.Getwd()
	Config = &config{
		ConfigFilePath:               configFile,
		ApplicationHome:              wd,
		ServerAddr:                   "127.0.0.1:8080",
		CookieStoreAuthenticationKey: generateRandomKey(32),
		CookieStoreEncryptionKey:     generateRandomKey(32),
		Repository:                   make(map[string]string),
	}

	os.MkdirAll(configHome, 0755)

	file, err := os.Open(configFile)
	if err != nil {
		log.Println("Creating configuration file")
		if file, err = os.Create(configFile); err != nil {
			panic(err)
		}
		file.Close()
		SaveConfig()
	} else {
		defer file.Close()
		log.Println("Loading configuration file")
		unmarshalToConfig(file)
	}
	printConfigToStdout()
}

func getMarshalledConfig() []byte {
	b, _ := json.MarshalIndent(Config, "", "  ")
	return b
}

func unmarshalToConfig(file *os.File) {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(file)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(buffer.Bytes(), Config)
}

func SaveConfig() {
	file, err := os.OpenFile(configFile, os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(getMarshalledConfig())
}

func printConfigToStdout() {
	cVal := reflect.ValueOf(Config).Elem()
	cType := cVal.Type()
	for x := 0; x < cType.NumField(); x++ {
		cFieldType := cType.Field(x)
		cFieldVal := cVal.Field(x)
		log.Printf("%s: %v", cFieldType.Name, cFieldVal.Interface())
	}
}

func generateRandomKey(strength int) []byte {
	k := make([]byte, strength)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}
