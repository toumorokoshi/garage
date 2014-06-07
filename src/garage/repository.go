// contains information about a garage repository
// (a directory with a garagerc contained inside)
package garage

import (
	"log"
	"io/ioutil"
)

type GarageRepository struct {
	RootPath string
	Scripts []string
}

func LoadGarageRepository(rootPath string) *GarageRepository {
	entries, err := ioutil.ReadDir(rootPath)
	if err != nil {
		log.Fatal("Unable to read directory: ", err)
	}

	scripts := make([]string, len(entries), len(entries))
	for i, entry := range entries {
		scripts[i] = entry.Name()
	}

	return &GarageRepository{rootPath, scripts}
}
