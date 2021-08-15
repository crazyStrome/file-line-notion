package fln

import (
	"reflect"
	"strconv"
)

type parseFunc func(reflect.Value, string) error

var parseValueFuncs map[reflect.Kind]parseFunc

func RegisteParseFunc(pf parseFunc, ks ...reflect.Kind) {
	for _, k := range ks {
		parseValueFuncs[k] = pf
	}
}

func init() {
	parseValueFuncs = make(map[reflect.Kind]parseFunc)
	RegisteParseFunc(func(v reflect.Value, s string) error {
		data, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		v.SetInt(int64(data))
		return nil
	}, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8)
	RegisteParseFunc(func(v reflect.Value, s string) error {
		data, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(data)
		return nil
	}, reflect.Bool)
	RegisteParseFunc(func(v reflect.Value, s string) error {
		data, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(data)
		return nil
	}, reflect.Uint)
	RegisteParseFunc(func(v reflect.Value, s string) error {
		v.SetString(s)
		return nil
	}, reflect.String)
	RegisteParseFunc(func(v reflect.Value, s string) error {
		data, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(data)
		return nil
	}, reflect.Float64)
	RegisteParseFunc(func(v reflect.Value, s string) error {
		data, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return err
		}
		v.SetFloat(data)
		return nil
	}, reflect.Float32)
}
