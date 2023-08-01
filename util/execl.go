package util

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"reflect"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

/**
 * ExportExcel 导出excel
 *
 * Example:
 *
 * type ReportExcel struct {
 * 	Dt         string `excel:"日期"`
 * 	NewUsers   int64  `excel:"注册人数"`
 * 	LoginUsers int64  `excel:"登录人数"`
 * 	Tmp        int64  `excel:"-"`
 * }
 * tbody := []ReportExcel{
 * 	{"2006-01-02", 1, 2},
 * }
 * buf, err := ExportExcel(ReportExcel{}, tbody)
 * if err != nil {
 * 	//
 * }
 * os.WriteFile("test.xlsx", buf.Bytes(), 0644)
 */
func ExportExcel(table interface{}, tbody interface{}) (*bytes.Buffer, error) {
	// make sure 'table' is a Struct
	tableVal := reflect.ValueOf(table)
	if tableVal.Kind() == reflect.Ptr {
		tableVal = tableVal.Elem()
	}
	if tableVal.Kind() != reflect.Struct {
		return nil, fmt.Errorf("table not struct")
	}
	// make sure 'tbody' is a Slice
	tbodyVal := reflect.ValueOf(tbody)
	if tbodyVal.Kind() != reflect.Slice {
		return nil, fmt.Errorf("tbody not slice")
	}
	var (
		tableFieldNum = tableVal.NumField()
		tbodyLen      = tbodyVal.Len()
		sheet         = "Sheet1"
		headCol       = 1
	)

	xlsx := excelize.NewFile()
	index, err := xlsx.NewSheet(sheet)
	if err != nil {
		return nil, fmt.Errorf("new sheet error: %s", err)
	}

	// write table head
	for i := 0; i < tableFieldNum; i++ {
		t := tableVal.Type().Field(i)

		if !ast.IsExported(t.Name) {
			continue
		}
		head := t.Tag.Get("excel")
		// is ignored field
		if head == "-" {
			continue
		}
		if len(head) == 0 {
			head = t.Name
		}
		axis, err := excelize.CoordinatesToCellName(headCol, 1)
		if err != nil {
			return nil, fmt.Errorf("head to cell name error: %s", err)
		}
		headCol++
		xlsx.SetCellValue(sheet, axis, head)
	}

	// write tbody
	for i := 0; i < tbodyLen; i++ {
		v := tbodyVal.Index(i)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			return nil, fmt.Errorf("tbody element not struct")
		}
		num_field := v.NumField() // number of fields in struct
		bodyCol := 1

		for j := 0; j < num_field; j++ {
			axis, err := excelize.CoordinatesToCellName(bodyCol, i+2)
			if err != nil {
				return nil, fmt.Errorf("tbody to cell name error: %s", err)
			}
			f := v.Field(j)
			t := v.Type().Field(j)

			if !ast.IsExported(t.Name) {
				continue
			}
			head := t.Tag.Get("excel")
			// is ignored field
			if head == "-" {
				continue
			}
			bodyCol++
			xlsx.SetCellValue(sheet, axis, f.Interface())
		}
	}
	xlsx.SetActiveSheet(index)

	buf, err := xlsx.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("write to buffer error: %s", err)
	}
	return buf, nil
}

/**
 * ReadExcelToStruct 读取excel到结构体
 *
 * Example:
 *
 * type ReportExcel struct {
 * 	Dt         string `excel:"日期"`
 * 	NewUsers   int64  `excel:"注册人数"`
 * 	LoginUsers int64  `excel:"登录人数"`
 * 	Tmp        int64  `excel:"-"`
 * }
 *
 * records, err := ReadExcelToStruct("test.xlsx", ReportExcel{})
 * if err != nil {
 * 	//
 * }
 * fmt.Println(records)
 */
func ReadExcelToStruct[T any](filename string, body T) ([]T, error) {
	var (
		fieldIndexMap = map[string]int{}

		bodyType = reflect.TypeOf(body)
		fieldNum = bodyType.NumField()
	)
	if t := bodyType.Kind(); t != reflect.Struct {
		return nil, fmt.Errorf("body must be a struct, currently %s", t.String())
	}
	// get field index
	for i := 0; i < fieldNum; i++ {
		f := bodyType.Field(i)
		head := f.Tag.Get("excel")
		if head == "-" {
			continue
		}
		if len(head) == 0 {
			head = f.Name
		}
		fieldIndexMap[head] = i
	}
	// open excel file
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("open file error: %s", err)
	}
	defer f.Close()

	// get sheet
	sheet := f.GetSheetList()
	if len(sheet) < 1 {
		return nil, errors.New("sheet is empty")
	}
	// get all the rows
	rows, err := f.GetRows(sheet[0])
	if err != nil {
		return nil, fmt.Errorf("get rows error: %s", err)
	}
	if len(rows) < 2 {
		return []T{}, nil
	}
	var (
		r    = []T{}
		head = rows[0]
	)
	for rowKey, row := range rows[1:] {
		var (
			bodyCopy = body
			t        = &bodyCopy
			rv       = reflect.ValueOf(t).Elem()
		)
		for colKey, colCell := range row {
			if len(head) <= colKey {
				continue
			}
			colCell = strings.TrimSpace(colCell)

			h := head[colKey]
			fieldIndex, ok := fieldIndexMap[h]
			if !ok {
				continue
			}
			fieldVal := rv.Field(fieldIndex)

			switch fieldVal.Kind() {
			case reflect.String:
				fieldVal.SetString(colCell)
			case reflect.Int, reflect.Int32, reflect.Int64:
				if len(colCell) == 0 {
					continue
				}
				v, err := strconv.Atoi(colCell)
				if err != nil {
					return nil, fmt.Errorf("row(%d) col(%d) head(%s) strconv.Atoi(%s) error: %s", rowKey+1, colKey+1, h, colCell, err)
				}
				fieldVal.SetInt(int64(v))
			}
		}
		r = append(r, *t)
	}
	return r, nil
}
