package util

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"time"
)

func Md5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
func FormatTime() string {
	return time.Now().Format("2006.01.02 15:04:05")
}

func ConvertIds[T any](tlist []*T, key string) []int64 {
	ids := []int64{}
	temp := make(map[int64]bool, len(tlist))
	for _, item := range tlist {
		vobj := reflect.ValueOf(*item)
		v := vobj.FieldByName(key)
		if v.Kind() != reflect.Invalid {
			_, ok := temp[v.Int()]
			if !ok {
				ids = append(ids, v.Int())
				temp[v.Int()] = true
			}

		}
	}
	return ids
}

func ConvertIdMap[T any](tlist []*T) map[int64]*T {
	ids := []int64{}
	idMap := make(map[int64]*T)
	for _, item := range tlist {
		vobj := reflect.ValueOf(*item)
		v := vobj.FieldByName("Id")
		if v.Kind() == reflect.Invalid {
			continue
		}
		id := v.Int()
		ids = append(ids, id)
		idMap[id] = item
	}
	return idMap
}
