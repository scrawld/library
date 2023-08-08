package util

import (
	"os"
	"testing"
)

// TestExportExcel
func TestExportExcel(t *testing.T) {
	type ReportExcel struct {
		Dt         string `excel:"日期"`
		NewUsers   int64  `excel:"注册人数"`
		LoginUsers int64  `excel:"登录人数"`
		Tmp        int64  `excel:"-"`
	}
	tbody := []ReportExcel{
		{Dt: "2006-01-02", NewUsers: 1, LoginUsers: 2, Tmp: 1},
	}
	buf, err := ExportExcel(ReportExcel{}, tbody, true)
	if err != nil {
		t.Errorf("export error: %s", err)
		return
	}
	os.WriteFile("test.xlsx", buf.Bytes(), 0644)
}
