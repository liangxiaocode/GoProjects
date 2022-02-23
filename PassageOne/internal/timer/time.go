package timer

import "time"

func GetNowTime() time.Time {
	return time.Now()
}

// 用parseDuration处理传入的未知字符串
func GetCalTime(currentTime time.Time, d string) (time.Time, error) {
	duration, err := time.ParseDuration(d)
	if err != nil {
		return time.Time{}, err
	}
	return currentTime.Add(duration), nil
}
