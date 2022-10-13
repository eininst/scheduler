package util

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func CurrentFile() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic(errors.New("Can not get current file info"))
	}
	return file
}

func Md5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
func FormatTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func FormatTimeMill() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
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

func Unicode2utf8(source string) string {
	var res = []string{""}
	sUnicode := strings.Split(source, "\\u")
	var context = ""
	for _, v := range sUnicode {
		var additional = ""
		if len(v) < 1 {
			continue
		}
		if len(v) > 4 {
			rs := []rune(v)
			v = string(rs[:4])
			additional = string(rs[4:])
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			context += v
		}
		context += fmt.Sprintf("%c", temp)
		context += additional
	}
	res = append(res, context)
	return strings.Join(res, "")
}
