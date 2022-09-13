package configs

import (
	"encoding/json"
	"github.com/eininst/flog"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"sync"
)

var data map[string]any
var ret gjson.Result
var once sync.Once

func SetConfig(conf_path string) {
	once.Do(func() {
		profile := os.Getenv("ENV")
		if profile == "" {
			profile = "dev"
		}
		flog.Infof("profile is: %s", profile)

		file, err := os.Open(conf_path)
		defer func() { _ = file.Close() }()
		if err != nil {
			log.Fatal(err)
		}
		dec := yaml.NewDecoder(file)
		err = dec.Decode(&data)

		for {
			var t map[string]interface{}
			err = dec.Decode(&t)
			if err != nil {
				break
			}
			if p, ok := t["profile"]; ok {
				if p == profile {
					for k, v := range t {
						data[k] = v
					}
					break
				}
			}
		}
		v, er := json.Marshal(&data)
		if er != nil {
			log.Println(er)
		}
		ret = gjson.Parse(string(v))
	})
}

func Get(path ...string) gjson.Result {
	if len(path) == 0 {
		return ret
	}

	var r gjson.Result
	for _, p := range path {
		if r.Value() == nil {
			r = ret.Get(p)
		} else {
			r = r.Get(p)
		}
	}
	return r
}
