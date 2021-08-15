#  实现一个简单的文件反序列化器

##  1.  需求

现在有一个文件，文件中的第一行是"name,address,phone,country,male,age"表示这个文件的后续内容类型，可以视为列名。之后的每一行都是这几部分数据，使用","分割。例如，从第二行开始后续的每一行的内容大致为："crastom,hone,111111111,china,true,20"。

如果想要提取这些内容，是不是很简单，只需要使用：

```go
strings.Split(line, ",")
```

就可以获得每一行中的各部分内容。然后把每部分数据赋值给一个结构体，例如：

```go
type Person struct {
  Name string
  Address string
  Phone string
  Country string
  Fale bool
  Age int
}
```

这样就完成了，但是如果之后需要解析更多的字段呢，或者需要解析的字段类型出现变化呢。

因此，本文就用go实现一种简单的Unmarshaler，它可以从文件中Unmarshal出所需要的数据，并且不需要写冗长的赋值语句；可以适用于不同的文件内容。

##  2.  思路

实现的思路也比较简单：

1. 使用第一行headLine，来初始化一个Unmarshaler，分析headLine中每个name对应的位置。例如，headLine为"name,age,address,male"，那么name对应idx为0，age为1，依次类推。
2. 实现Unmarshal函数时，传入需要反序列化的一行数据line以及存放数据的结构体ds。结构体中通过字段的tag或者字段名获取该字段的数据在一行中对应的位置。例如，line为"crastom,20,home,true"，那么crastom就对应与headLine的name，以此类推。
3. 在Unmarshaler中，找到数据的位置，从line中取出数据，然后通过反射设置ds对应字段的内容即可。line中获取到的内容都是string，而ds的字段中可能存在多种类型：int、bool、string、float64等。针对不同的类型，需要设计成为可注册的处理方式，这样遇到对应的类型直接取出对应的parseFunc即可处理。

##  3.  golang实现

###  3.1.  Unmarshaler数据结构

Unmarshaler的数据结构定义为fln，options是fln的相关配置；total是headLine中的列数；headToIdx是一个map，将headLine中列名与它的位置idx对应起来。

```go
type fln struct {
	options   *Options
	total     int
	headToIdx map[string]int
}
```

###  3.2.  初始化Unmarshaler

```go
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
```

这个函数用来新建一个Unmarshaler，传入的参数oos是用来配置Options的，Options和Option的定义如下：

```go
// Options 解析参数
type Options struct {
	// 文件的第一行，
	// 例如:"name,age,country"这些声明字段
	HeadLine string
	// 文件中每一行各部分内容
	// 的分割符，默认使用"\t"
	Spliter string
}

// Option 用来设置options
type Option func(*Options)

// WithHeadLine 向Options中添加headLine
func WithHeadLine(line string) Option {
	return func(o *Options) {
		o.HeadLine = line
	}
}

// WithSpliter 设置options中的spliter
func WithSpliter(spliter string) Option {
	return func(o *Options) {
		o.Spliter = spliter
	}
}
```

Options就是fln的相关参数配置，而Option就是用来处理Options的函数，目前有WithHeadLine以及WithSpliter这两个函数。

而上面的parseHeadLine实现很简单，就是把headLine通过split分割成string数据，然后映射到headToIdx中。

###  3.3.  Unmarshal实现

方法签名如下：

```go
func (f *fln) Unmarshal(data []byte, ptr interface{}) error
```

data即每一行需要反序列化的数据，ptr则是一个结构体指针，用来存放数据。

接下来是简略实现思路：

1. 将data转化为string然后分割成string数组：datas
2. 对ptr指向的结构体中字段遍历，跳过无法设置值的字段。
3. 通过字段名或tag获取该字段对应的数据在datas中的位置idx，然后设置该字段的值为datas[idx]

```go
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
```

在上面的setValue函数中，首先将content转换为fieldt的类型，然后通过fieldv.SetXxx进行设置。

```go
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
```

setValue函数中，在parseValueFuncs找到对应的转换函数parseFunc：pf，然后用执行pf。

那parseValeFuncs中的parseFunc是如何设置的呢？

```go
type parseFunc func(reflect.Value, string) error

var parseValueFuncs map[reflect.Kind]parseFunc

func RegisteParseFunc(pf parseFunc, ks ...reflect.Kind) {
	for _, k := range ks {
		parseValueFuncs[k] = pf
	}
}
```

在init函数中，已经实现了int、float64、string、bool等类型的parseFunc，对于其他的类型，可以自己实现，然后注册到fln中。

###  3.4.  测试

附加一个简单的例子，帮助理解。

```go
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
```

##  4.  总结

这个反序列化的工具比较简单，主要内容就是使用反射设置字段数据。但是fln的配置Options以及parseFunc的注册还是值得一看的，方便后续新功能的添加。
