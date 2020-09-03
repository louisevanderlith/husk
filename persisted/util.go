package persisted

import (
	"os"
	"path/filepath"
)

//CreateDirectory will create the specified path if it doesn't exist
func CreateDirectory(folderPath string) error {
	_, err := os.Stat(folderPath)

	if os.IsNotExist(err) {
		return os.Mkdir(folderPath, 0755)
	}

	return err
}

//DestroyContents will remove the file at the given path
func DestroyContents(path string) error {
	d, err := os.Open(path)

	if err != nil {
		return err
	}

	defer d.Close()

	names, err := d.Readdirnames(-1)

	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path, name))

		if err != nil {
			return err
		}
	}

	return nil
}
