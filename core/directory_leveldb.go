package core

import (
	"encoding/binary"
	"github.com/syndtr/goleveldb/leveldb"
	"path/filepath"
)

type LeveldbDirectory struct {
	db   *leveldb.DB
	path string // leveldb 文件存放路径
}

func NewLeveldbDirectory(dir string) (d *LeveldbDirectory, err error) {
	d = new(LeveldbDirectory)
	d.path = filepath.Join(dir, "index") // TODO all volumes in one directory. Future: one volume one directory
	d.db, err = leveldb.OpenFile(d.path, nil)
	return
}

func (d *LeveldbDirectory) Get(id uint64) (n *Needle, err error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, id)
	data, err := d.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	n, err = NeedleUnmarshal(data)

	return
}

func (d *LeveldbDirectory) Set(n *Needle) (err error) {
	data, err := NeedleMarshal(n)
	if err != nil {
		return err
	}
	return d.db.Put(data[:8], data, nil)
}

func (d *LeveldbDirectory) Has(id uint64) (has bool) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, id)
	_, err := d.db.Get(key, nil)
	return err == nil
}

func (d *LeveldbDirectory) Del(id uint64) (err error) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, id)
	return d.db.Delete(key, nil)
}
