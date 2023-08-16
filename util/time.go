package util

import (
	"math"
	"strconv"
	"time"
)

// GetDtByOffset 获取日期 GetDtByOffset(time.Now(), 0) return 20060102
func GetDtByOffset(tm time.Time, offset int) (r int) {
	r, _ = strconv.Atoi(tm.AddDate(0, 0, offset).Format("20060102"))
	return
}

// DtToTime 日期缩写转时间
func DtToTime(dt int) (r time.Time) {
	tim, _ := time.ParseInLocation("20060102", strconv.Itoa(dt), time.Local)
	return tim
}

// StartOfDay 获取指定日期零点时间
func StartOfDay(t time.Time) (r time.Time) {
	r = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return
}

// DaysBetween 计算两个时间之间的天数差异
func DaysBetween(start, end time.Time) int {
	startZero := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endZero := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	duration := endZero.Sub(startZero)
	return int(math.Abs(duration.Hours()) / 24)
}

type TimeRange struct {
	StartTime time.Time
	EndTime   time.Time
}

/**
 * GetTimeRangesByDay 按天划分时间范围,保留开始时间和结束时间
 *
 * Example:
 *
 * startTime := time.Unix(1688870348, 0) // 2023-7-9 10:39:08
 * endTime := time.Unix(1689059808, 0)   // 2023-7-11 15:16:48
 *
 * result := GetTimeRangesByDay(startTime, endTime)
 *
 * for _, v := range result {
 * 	fmt.Printf("startTime: %s, endTime: %s\n", v.StartTime.Format("2006-01-02 15:04:05"), v.EndTime.Format("2006-01-02 15:04:05"))
 * 	//[
 * 	//	["2023-07-09 10:39:08", "2023-07-10 00:00:00"],
 * 	//	["2023-07-10 00:00:00", "2023-07-11 00:00:00"],
 * 	//	["2023-07-11 00:00:00", "2023-07-11 15:16:48"],
 * 	//]
 * }
 */
func GetTimeRangesByDay(st, et time.Time) []*TimeRange {
	var (
		stZero       = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
		etZero       = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())
		intervalDays = int(etZero.Sub(stZero).Hours() / 24) // 间隔天数
		dtLayout     = "20060102"
		r            = []*TimeRange{}
	)

	for i := 0; i <= intervalDays; i++ {
		t := &TimeRange{
			StartTime: stZero.AddDate(0, 0, i),
			EndTime:   stZero.AddDate(0, 0, i+1),
		}
		if t.StartTime.Format(dtLayout) == st.Format(dtLayout) { // 比较年月日是否相同
			t.StartTime = st
		}
		if t.EndTime.Add(-time.Second).Format(dtLayout) == et.Format(dtLayout) {
			t.EndTime = et
		}
		r = append(r, t)
	}
	return r
}

/**
 * GetDateRange 获取日期范围,从零点开始
 *
 * Example:
 *
 * startTime := time.Unix(1688870348, 0) // 2023-7-9 10:39:08
 * endTime := time.Unix(1689059808, 0)   // 2023-7-11 15:16:48
 *
 * result := GetDateRange(startTime, endTime)
 *
 * for _, v := range result {
 * 	fmt.Println(v.Format("2006-01-02 15:04:05"))
 * 	//[
 * 	//	"2023-07-09 00:00:00",
 * 	//	"2023-07-10 00:00:00",
 * 	//	"2023-07-11 00:00:00",
 * 	//]
 * }
 */
func GetDateRange(st, et time.Time) []time.Time {
	var (
		stZero       = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
		etZero       = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())
		intervalDays = int(etZero.Sub(stZero).Hours() / 24) // 间隔天数
		r            = []time.Time{}
	)

	for i := 0; i <= intervalDays; i++ {
		r = append(r, stZero.AddDate(0, 0, i))
	}
	return r
}

// GetWeek 获取周次 GetWeek(time.Unix(1672588800, 0)) return 202301
func GetWeek(tm time.Time) int {
	year, week := tm.ISOWeek()
	return year*100 + week
}

// SplitYearWeek 分割年和周次 SplitYearWeek(202301) return 2023, 01
func SplitYearWeek(combined int) (year, week int) {
	year = combined / 100
	week = combined % 100
	return year, week
}

// GetWeekRange 获取时间范围内的周次 GetWeekRange(time.Unix(1671790242, 0), time.Unix(1672588800, 0)) return []int{202251, 202252, 202301}
func GetWeekRange(st, et time.Time) []int {
	var (
		stZero = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
		etZero = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())
		r      = []int{}

		intervalDays = int(etZero.Sub(stZero).Hours() / 24) // 间隔天数
		seen         = map[int]struct{}{}                   // 用于记录已经处理过的周次
	)
	for i := 0; i <= intervalDays; i++ {
		week := GetWeek(stZero.AddDate(0, 0, i))
		if _, ok := seen[week]; !ok {
			seen[week] = struct{}{}
			r = append(r, week)
		}
	}
	return r
}

// GetDataHour 获取日期到小时 GetDataHour(time.Now()) return 2022122609
func GetDataHour(tm time.Time) (r int) {
	r, _ = strconv.Atoi(tm.Format("2006010215"))
	return
}

// GetDataHourRange 获取日期范围 GetDataHourRange(time.Unix(1672012800, 0), time.Unix(1672016400, 0)) ruturn []int{2022122608 2022122609}
func GetDataHourRange(st, et time.Time) (r []int) {
	var (
		stZero   = time.Date(st.Year(), st.Month(), st.Day(), st.Hour(), 0, 0, 0, time.Local)
		etZero   = time.Date(et.Year(), et.Month(), et.Day(), et.Hour(), 0, 0, 0, time.Local)
		interval = etZero.Sub(stZero).Hours()
	)
	for i := 0; i <= int(interval); i++ {
		var (
			t  = stZero.Add(time.Hour * time.Duration(i))
			dh = GetDataHour(t)
		)
		r = append(r, dh)
	}
	return
}

// GetMonth 获取月份 GetMonth(time.Now()) return 202308
func GetMonth(tm time.Time) (r int) {
	r, _ = strconv.Atoi(tm.Format("200601"))
	return
}

// GetMonthRange 获取时间范围内的月份 GetMonthRange(time.Unix(1672016400, 0), time.Unix(1674694800, 0)) return [202212 202301]
func GetMonthRange(st, et time.Time) []int {
	var (
		stZero = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
		etZero = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())
		r      = []int{}

		intervalDays = int(etZero.Sub(stZero).Hours() / 24) // 间隔天数
		seen         = map[int]struct{}{}
	)
	for i := 0; i <= intervalDays; i++ {
		month := GetMonth(stZero.AddDate(0, 0, i))
		if _, ok := seen[month]; !ok {
			seen[month] = struct{}{}
			r = append(r, month)
		}
	}
	return r
}
