package tools

import (
	"bytes"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/pborman/uuid"
)

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769")

// NewID is a globally unique identifier.  It is a [A-Z0-9] string 26
// characters long.  It is a UUID version 4 Guid that is zbased32 encoded
// with the padding stripped off.
func NewID() string {
	var b bytes.Buffer
	encoder := base32.NewEncoder(encoding, &b)
	encoder.Write(uuid.NewRandom())
	encoder.Close()
	b.Truncate(26) // removes the '==' padding
	return b.String()
}

// GetMillis return now time( millisecond )
func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// GetSeconds return now time( second )
func GetSeconds() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

// PrintTime 한국 기준으로 시간을 출력 ( YYYY-MM-DD HH:mm:ss )
func PrintTimeFromSec(epoch int64) string {
	t, _ := GetLocalTime(epoch*1000, "Korea")
	strTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return strTime
}

// PrintTime 한국 기준으로 시간을 출력 ( YYYY-MM-DD HH:mm:ss )
func PrintTimeFromMilli(epoch int64) string {
	t, _ := GetLocalTime(epoch, "Korea")
	strTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return strTime
}

var countryTz = map[string]string{
	"Korea": "Asia/Seoul",
}

func msToTime(ms int64) (time.Time, error) {
	return time.Unix(0, ms*int64(time.Millisecond)), nil
}

// GetLocalTime UTC 타임을 변환
func GetLocalTime(epochMilli int64, location string) (*time.Time, error) {
	loc, err := time.LoadLocation(countryTz[location])
	if err != nil {
		return nil, fmt.Errorf("error. failed to get location[ %s ]", err.Error())
	}
	inTime, err := msToTime(epochMilli)
	if err != nil {
		return nil, fmt.Errorf("error. failed to convert epoch millisecond to time[ %s ]", err.Error())
	}
	retTime := inTime.In(loc)
	return &retTime, nil
}
