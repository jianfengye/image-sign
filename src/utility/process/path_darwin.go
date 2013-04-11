package process

import (
	"os"
	"os/exec"
	"path"
)

func RootPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}

	dir, _ := path.Split(file)

	os.Chdir(dir + "/../")
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd, nil
}
