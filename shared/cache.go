package shared

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
	"github.com/vegidio/umd"
)

func init() {
	gob.Register(time.Time{})
}

func LoadCache(filePath string) (*umd.Response, error) {
	var response umd.Response

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&response)
	return &response, err
}

func SaveCache(filePath string, response *umd.Response) error {
	dir := filepath.Dir(filePath)
	if exists, _ := afero.DirExists(fs, dir); !exists {
		_ = fs.MkdirAll(dir, 0755)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(response)
}
