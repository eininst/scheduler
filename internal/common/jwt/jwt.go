package jwt

import (
	"encoding/base64"
	"encoding/json"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/util"
	"time"
)

type Jwt struct {
	SecretKey string
}

type Token struct {
	Data interface{} `json:"data"`
	Exp  int64       `json:"exp"`
}

func New(secretKey string) *Jwt {
	return &Jwt{SecretKey: secretKey}
}

func (j *Jwt) CreateToken(data interface{}, expire time.Duration) string {
	b, _ := json.Marshal(&Token{
		Data: data,
		Exp:  time.Now().UnixNano() + int64(expire),
	})
	result, err := util.AesEncrypt(b, []byte(j.SecretKey))
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(result)
}

func (j *Jwt) ParseToken(token string) (interface{}, error) {
	b, _ := base64.StdEncoding.DecodeString(token)
	origData, err := util.AesDecrypt(b, []byte(j.SecretKey))
	if err != nil {
		return nil, err
	}
	var tk Token
	err = json.Unmarshal(origData, &tk)
	if err != nil {
		return nil, err
	}
	if time.Now().UnixNano() > tk.Exp {
		return nil, service.NewServiceError("token is expired")
	}
	return tk.Data, nil
}
