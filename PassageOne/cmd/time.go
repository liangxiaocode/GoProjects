package cmd

import (
	"PassageOne/internal/timer"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var calculateTime string
var duration string

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "时间格式处理",
	Long:  "时间格式处理",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var nowTimeCmd = &cobra.Command{
	Use:   "now",
	Short: "获取当前时间",
	Long:  "获取当前时间",
	Run: func(cmd *cobra.Command, args []string) {
		nowTime := timer.GetNowTime()
		log.Printf("输出结果: %s, %d", nowTime.Format("2006-01-02 15:04:05"), nowTime.Unix())
	},
}

var calculateTimeCmd = &cobra.Command{
	Use:   "calc",
	Short: "计算所需时间",
	Long:  "计算所需时间",
	Run: func(cmd *cobra.Command, args []string) {
		var currentTimer time.Time
		var layout = "2006-01-02 15:04:05"
		if calculateTime == "" {
			currentTimer = timer.GetNowTime()
		} else {
			var err error
			spaceNum := strings.Count(calculateTime, " ")
			if spaceNum == 0 {
				layout = "2006-01-02"
			}
			if spaceNum == 1 {
				layout = "2006-01-02 15:04:05"
			}
			// 解析字符串
			currentTimer, err = time.Parse(layout, calculateTime)
			// 若出现异常，直接按时间戳格式处理
			if err != nil {
				// 转化为数字
				t, _ := strconv.Atoi(calculateTime)
				currentTimer = time.Unix(int64(t), 0)
			}
		}
		t, err := timer.GetCalTime(currentTimer, duration)
		if err != nil {
			log.Fatalf("time.GetCalTime failed,err:%v\n", err)
		}
		log.Printf("输出结果: %s, %d", t.Format(layout), t.Unix())
	},
}

func init() {
	timeCmd.AddCommand(nowTimeCmd)
	timeCmd.AddCommand(calculateTimeCmd)

	calculateTimeCmd.Flags().StringVarP(&calculateTime, "calculate", "c", "", "需要计算的时间,单位为时间戳或者格式化后的时间")
	calculateTimeCmd.Flags().StringVarP(&duration, "dur", "d", "", "持续的时间,支持的有效单位是ns,us,ms,s,m和h")

}
