package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	lib "telegram-bot-go/lib/e"
	"telegram-bot-go/storage"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

var ErrNoSavedPage = errors.New("no saved page")

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = lib.WrapIfErr("can't save page", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil

}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = lib.WrapIfErr("can't pick random page", err) }()

	fPath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPage
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(fPath, file.Name()))
}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, lib.Wrap("can't remove file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)
		return false, lib.Wrap(msg, nil)
	}

	return true, nil

}

func (s Storage) Remove(p *storage.Page) error {
	//defer func() { err = lib.WrapIfErr("can't remove page", err) }()
	fName, err := fileName(p)
	if err != nil {
		return lib.Wrap("can't remove file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)
		return lib.Wrap(msg, err)
	}

	return nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, lib.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewEncoder(f).Encode(&p); err != nil {
		return nil, err
	}

	return &p, nil

}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
