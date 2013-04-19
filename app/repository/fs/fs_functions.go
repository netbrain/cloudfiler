package fs

import (
	"encoding/json"
	"fmt"
	"github.com/netbrain/cloudfiler/app/conf"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

var idSeq int
var storagePath, storageTmpPath, idSeqPath string

const (
	storagePathKey = "storage-path"
	storageTmpKey  = "storage-tmp-path"
)

func init() {
	var present bool
	var saveConf bool
	storagePath, present = conf.Config.Repository[storagePathKey]
	if !present {
		storagePath = filepath.Join(conf.Config.ApplicationHome, "db")
		conf.Config.Repository[storagePathKey] = storagePath
		saveConf = true
	}

	storageTmpPath, present = conf.Config.Repository[storageTmpKey]
	if !present {
		storageTmpPath = filepath.Join(conf.Config.ApplicationHome, "tmp")
		conf.Config.Repository[storageTmpKey] = storageTmpPath
		conf.SaveConfig()
	}

	if saveConf {
		conf.SaveConfig()
	}

	idSeqPath = filepath.Join(storagePath, "idseq")

	b, err := ioutil.ReadFile(idSeqPath)

	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		idSeq, err = strconv.Atoi(string(b))
		if err != nil {
			panic(err)
		}
	}
}

func generateID() int {
	idSeq++

	file, err := os.Create(idSeqPath)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	file.Write([]byte(strconv.Itoa(idSeq)))

	return idSeq
}

func serialize(data interface{}) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return b
}

func unserialize(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}

func getPath(folder string, id interface{}) string {
	parent := filepath.Join(storagePath, folder)
	file := fmt.Sprintf("%v", id)

	err := os.MkdirAll(parent, 0700)
	if err != nil {
		panic(err)
	}

	return filepath.Join(parent, file)
}

func getTempFile() *os.File {
	file, err := ioutil.TempFile(getPath(storagePath, ""), "")
	if err != nil {
		panic(err)
	}
	return file
}
