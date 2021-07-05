package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// JSONDate custom json data format
type JSONDate struct {
	time.Time
}

// // MarshalJSON trans time to josnstr
// func (this JSONDate) MarshalJSON() ([]byte, error) {
// 	var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("2006-01-02"))
// 	return []byte(stamp), nil
// }

// MarshalJSON trans time to josnstr
func (j JSONDate) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", j.Format("2006-01-02"))
	return []byte(stamp), nil
}

// UnmarshalJSON trans string to time.Time
func (j *JSONDate) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		t := time.Now()
		*j = JSONDate{Time: t}
	} else {
		t, err := time.Parse("\"2006-01-02\"", string(data))
		fmt.Println(err)
		*j = JSONDate{Time: t}
	}

	return nil
}

// String print custom tpye
// func (j JSONDate) String() string {
// 	return time.Time(j).Format("2006-01-02T15:04:05Z07:00")
// }

//Value insert timestamp into mysql need this function.
func (j JSONDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	if j.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return j.Time, nil
}

//Scan valueof time.Time
func (j *JSONDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*j = JSONDate{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
