package util

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"go/ast"
	"reflect"
)

/**
 * ExportCSV 导出csv
 *
 * Example:
 *
 * type Person struct {
 * 	Name  string `csv:"Name"`
 * 	Age   int    `csv:"Age"`
 * 	Email string `csv:"Email"`
 * }
 * tbody := []Person{
 * 	{Name: "Alice", Age: 25, Email: "alice@example.com"},
 * 	{Name: "Bob", Age: 30, Email: "bob@example.com"},
 * }
 * buf, err := ExportCSV(Person{}, tbody)
 * if err != nil {
 * 	//
 * }
 * os.WriteFile("test.csv", buf.Bytes(), 0644)
 */
func ExportCSV[T any](table T, tbody []T) (*bytes.Buffer, error) {
	// make sure 'table' is a Struct
	tableVal := reflect.ValueOf(table)
	if tableVal.Kind() == reflect.Ptr {
		tableVal = tableVal.Elem()
	}
	if tableVal.Kind() != reflect.Struct {
		return nil, fmt.Errorf("table not struct")
	}
	var (
		buffer = bytes.NewBuffer(nil)
		tag    = "csv"

		tableFieldNum = tableVal.NumField()
		tbodyLen      = len(tbody)
		thead         = []string{}
	)
	// get table head
	for i := 0; i < tableFieldNum; i++ {
		field := tableVal.Type().Field(i)

		if !ast.IsExported(field.Name) {
			continue
		}
		head := field.Tag.Get(tag)
		// is ignored field
		if head == "-" {
			continue
		}
		if len(head) == 0 {
			head = field.Name
		}
		thead = append(thead, head)
	}
	// new CSV writer
	writer := csv.NewWriter(buffer)
	// write table head
	if err := writer.Write(thead); err != nil {
		return nil, fmt.Errorf("write table head error: %s", err)
	}
	// write tbody
	for i := 0; i < tbodyLen; i++ {
		var (
			record = []string{}
			value  = reflect.ValueOf(tbody[i])
		)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		fieldNum := value.NumField()

		for j := 0; j < fieldNum; j++ {
			field := value.Field(j)
			fieldT := value.Type().Field(j)

			if !ast.IsExported(fieldT.Name) {
				continue
			}
			head := fieldT.Tag.Get(tag)
			// is ignored field
			if head == "-" {
				continue
			}
			record = append(record, fmt.Sprintf("%v", field.Interface()))
		}
		// write record head
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("write record error: %s, %v", err, record)
		}
	}
	// writer flush
	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("writer flush error: %s", err)
	}
	return buffer, nil
}
