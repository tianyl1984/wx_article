package util

import "time"

const (
	timeFormat = "2006-01-02 15:04:05"
)

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	if err == nil {
		*t = Time(now)
	}
	return err
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormat)
}