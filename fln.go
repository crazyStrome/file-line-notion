package fln

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
}

const (
	TAG_NAME = "fln"
)

type fln struct {
	options   *Options
	total     int
	headToIdx map[string]int
}

func (f *fln) Unmarshal(data []byte, ptr interface{}) error {
	err := f.parsePtrStruct(string(data), ptr)
	if err != nil {
		return err
	}
	return nil
}

func NewFln(oos ...Option) (Unmarshaler, error) {
	options := Options{}
	for _, o := range oos {
		o(&options)
	}
	err := checkOptions(&options)
	if err != nil {
		return nil, err
	}

	f := &fln{
		options:   &options,
		headToIdx: make(map[string]int),
	}
	err = f.parseHeadLine()
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (f *fln) parseHeadLine() error {
	heads := strings.Split(f.options.HeadLine, f.options.Spliter)
	for i, head := range heads {
		f.headToIdx[head] = i
	}
	f.total = len(heads)
	return nil
}
func checkOptions(o *Options) error {
	if o.HeadLine == "" {
		return fmt.Errorf("HeadLine is null: %+v", o)
	}
	if o.Spliter == "" {
		o.Spliter = "\t"
	}
	return nil
}
func (f *fln) parsePtrStruct(line string, ptr interface{}) error {
	data := strings.Split(line, f.options.Spliter)
	if len(data) != f.total {
		return fmt.Errorf("The len of data[%d] is not total[%d]", len(data), f.total)
	}
	ptrt := reflect.TypeOf(ptr)
	if ptrt.Kind() != reflect.Ptr {
		return fmt.Errorf("The input interface is not ptr: %+v", ptr)
	}
	elet := ptrt.Elem()
	if elet.Kind() != reflect.Struct {
		return fmt.Errorf("The input interface is not struct: %+v", ptr)
	}
	elev := reflect.ValueOf(ptr).Elem()

	for i := 0; i < elev.NumField(); i++ {
		fieldt := elet.Field(i)
		fieldv := elev.Field(i)
		if !fieldv.CanSet() {
			continue
		}
		tagName := fieldt.Tag.Get(TAG_NAME)
		fieldName := fieldt.Name
		idx := f.getIdxFromName(tagName, fieldName)
		if idx == -1 {
			continue
		}
		content := data[idx]
		setValue(fieldv, fieldt, content)
	}
	return nil
}
func setValue(fieldv reflect.Value, fieldt reflect.StructField, value string) error {
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("error from setValue: %+v", err)
		}
	}()
	pf, ok := parseValueFuncs[fieldt.Type.Kind()]
	if !ok {
		return fmt.Errorf("not suppoted type: %+v", fieldt.Type)
	}
	err = pf(fieldv, value)
	return err
}
func (f *fln) getIdxFromName(tagName, fieldName string) int {
	if idx, ok := f.headToIdx[tagName]; ok {
		return idx
	}
	if idx, ok := f.headToIdx[fieldName]; ok {
		return idx
	}
	if idx, ok := f.headToIdx[getSmallCamel(fieldName)]; ok {
		return idx
	}
	return -1
}
func getSmallCamel(name string) string {
	if len(name) == 0 {
		return name
	}
	bs := []byte(name)
	if bs[0] >= 'A' && bs[0] <= 'Z' {
		bs[0] = bs[0] - 'A' + 'a'
	}
	return string(bs)
}
