package record

import (
	"strings"
	"fmt"
)

// SplitKey takes a key in the form `/$namespace/$path` and splits it into
// `$namespace` and `$path`.
func SplitKey(key string) (string, string, error) {
	fmt.Println("#### SplitKey key: ", key)
	if len(key) == 0 || key[0] != '/' {
		return "", "", ErrInvalidRecordType
	}

	key = key[1:]

	fmt.Println("#### SplitKey key2: ", key)
	i := strings.IndexByte(key, '/')
	if i <= 0 {
		fmt.Println("#### SplitKey key err: ")
		return "", "", ErrInvalidRecordType
	}

	return key[:i], key[i+1:], nil
}
