package date

import "time"

// GetNowDay 返回当前时间的日期 固定格式 2006-01-02
func GetNowDay() string {
	return time.Now().UTC().Format("2006-01-02")
}

// GetBeforeDay 返回多少天之前的时间 固定格式 2006-01-02
func GetBeforeDay(i int) string {
	return time.Now().AddDate(0, 0, i).Format("2006-01-02")
}

// GetLastMonth1stDay 获取上个月的1号时间 固定格式 2006-01-02
func GetLastMonth1stDay() string {
	LastMonth := time.Now().AddDate(0, 0, -30).Format("2006-01")
	layout := "2006-01"
	t, _ := time.Parse(layout, LastMonth)
	return t.Format("2006-01-02")
}
