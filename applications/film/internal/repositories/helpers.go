package repositories

import "fmt"

func generateKey(prefix string, key interface{}) string {
	return fmt.Sprintf("%v.%v", prefix, key)
}
