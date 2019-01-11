package storage

type StorageKey struct {
	KeyAddress   uint64
	KeyType      int32
	IsArrayIndex bool
}

type StorageValue struct {
	ValueAddress uint64
	ValueType    int32
}

type StorageMapping struct {
	StorageValue  StorageValue
	StorageKey    []StorageKey
	StorageKeyMap map[string]bool
}
