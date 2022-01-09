package context

type Key string

type Context struct {
	values map[Key]interface{}
}

func New() *Context {
	return &Context{values: make(map[Key]interface{})}
}

func (c *Context) GetValue(key Key) interface{} {
	value, exist := c.values[key]
	if !exist {
		return nil
	}

	return value
}

func (c *Context) StoreValue(key Key, value interface{}) {
	c.values[key] = value
}
