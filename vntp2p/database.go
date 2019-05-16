// Copyright 2019 The go-vnt Authors
// This file is part of the go-vnt library.
//
// The go-vnt library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-vnt library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-vnt library. If not, see <http://www.gnu.org/licenses/>.

package vntp2p

import (
	"fmt"

	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	"github.com/syndtr/goleveldb/leveldb"
)

// LevelDB vntdb object
type LevelDB struct {
	path string
	db   *leveldb.DB
}

// 最好使用单例，因为leveldb只能有一个打开句柄，而且这个句柄是线程安全的。
var vntdb *LevelDB

// GetDatastore singleton design pattern
func GetDatastore(path string) (*LevelDB, error) {
	if vntdb != nil && vntdb.path == path {
		return vntdb, nil
	}
	vntdb, err := newDatastore(path)
	if err != nil {
		return nil, err
	}
	return vntdb, nil
}

func newDatastore(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{
		path: path,
		db:   db,
	}, nil
}

// Put implement Put() of ds.Batching interface
func (d *LevelDB) Put(key ds.Key, value interface{}) (err error) {
	byteKey := []byte(key.String())
	byteVal, ok := value.([]byte)
	if !ok {
		return ds.ErrInvalidType
	}
	err = d.db.Put(byteKey, byteVal, nil)
	if err != nil {
		fmt.Printf("leveldb put error = %s\n", err)
		return err
	}
	return nil
}

// Get implement Get() of ds.Batching interface
func (d *LevelDB) Get(key ds.Key) (value interface{}, err error) {
	byteKey := []byte(key.String())
	byteVal, err := d.db.Get(byteKey, nil)
	if err == leveldb.ErrNotFound {
		return nil, ds.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return byteVal, nil
}

// Has implement Has() of ds.Batching interface
func (d *LevelDB) Has(key ds.Key) (exists bool, err error) {
	byteKey := []byte(key.String())
	exists, err = d.db.Has(byteKey, nil)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Delete implement Delete() of ds.Batching interface
func (d *LevelDB) Delete(key ds.Key) (err error) {
	byteKey := []byte(key.String())
	err = d.db.Delete(byteKey, nil)
	if err != nil {
		return err
	}
	return nil
}

// Query implement Query() of ds.Batching interface
func (d *LevelDB) Query(q query.Query) (query.Results, error) {
	var re []query.Entry
	iter := d.db.NewIterator(nil, nil)
	for iter.Next() {
		keyByte := iter.Key()
		valueByte := iter.Value()
		re = append(re, query.Entry{Key: string(keyByte), Value: valueByte})
	}
	r := query.ResultsWithEntries(q, re)
	r = query.NaiveQueryApply(q, r)
	return r, nil
}

// Close implement Close() of ds.Batching interface
func (d *LevelDB) Close() error {
	err := d.db.Close()
	if err != nil {
		return err
	}
	return nil
}

// Batch implement Batch() of ds.Batching interface
func (d *LevelDB) Batch() (ds.Batch, error) {
	return ds.NewBasicBatch(d), nil
}
