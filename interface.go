package fln

type Unmarshaler interface {
	Unmarshal([]byte, interface{}) error
}
