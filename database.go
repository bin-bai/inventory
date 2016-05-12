// database
package inventory

type Database interface {
	Load() error
	Save() error
	Close()

	Use(col string, elemtype interface{}) Collection

	ColSize() int
}

type Collection interface {
	Set(key string, elem interface{}) error

	Get(key string) interface{}
	Del(key string)

	GetIndex(i int) (key string, elem interface{})

	Size() int
}
