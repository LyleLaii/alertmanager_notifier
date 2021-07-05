package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"reflect"
	"strings"
	"time"

	"alertmanager_notifier/config"
	"github.com/spf13/viper"
)

// TransTimeZone use custom timezone infor to tanse time
func TransTimeZone(t time.Time, tz string) string {
	ttz, _ := time.LoadLocation(tz)
	return t.In(ttz).Format(config.TimeLayout)
}

// TransTimeZoneAuto use machine timezone to trans time
func TransTimeZoneAuto(t time.Time) string {
	return t.In(time.Local).Format(config.TimeLayoutTZ)
}

// TransTimeZoneAutoCustom use machine timezone to trans time
func TransTimeZoneAutoCustom(t time.Time, formatStr string) string {
	return t.In(time.Local).Format(formatStr)
}

// StructToMap trans struct to map
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

// MergeStringMap merge map
func MergeStringMap(m1 map[string]string, m2 map[string]string) map[string]string {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func GetUseRote() bool {
	return viper.GetBool("userota")
}

func Check(content, encrypted string) bool {
	return strings.EqualFold(Encode(content), encrypted)
}
func Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

