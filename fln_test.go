package fln

import (
	"reflect"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

type DS struct {
	Name    string `fln:"myname"`
	MyAge   int    `fln:"age"`
	Address string `fln:"address"`
	Male    bool   `fln:"mymale"`
}

func TestNewFLN(t *testing.T) {
	head := "name,age,address,male"
	data := "crastom,10,home,true"
	want := DS{
		Name:    "crastom",
		MyAge:   10,
		Address: "home",
		Male:    true,
	}
	convey.Convey("test_new_fln", t, func() {
		f, err := NewFln(
			WithHeadLine(head),
			WithSpliter(","),
		)
		convey.So(f, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
		ds := DS{}
		err = f.Unmarshal([]byte(data), &ds)
		convey.So(err, convey.ShouldBeNil)
		convey.So(reflect.DeepEqual(want, ds), convey.ShouldBeTrue)
		t.Logf("%+v", ds)
	})
}
func TestType(t *testing.T) {
	type a struct{}
	t.Log(reflect.TypeOf(a{}).Kind())
	t.Log(reflect.TypeOf(a{}))
}
