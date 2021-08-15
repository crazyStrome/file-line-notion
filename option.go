package fln

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
