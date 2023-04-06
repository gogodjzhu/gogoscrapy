package entity

type IResultItems interface {
	Get(key string) interface{}
	Put(key string, value interface{})
	IsSkip() bool
	All() map[string]interface{}
}

type ResultItems struct {
	Items map[string]interface{}
	Skip  bool
}

func NewResultItems(skip bool) IResultItems {
	return &ResultItems{
		Items: make(map[string]interface{}),
		Skip:  skip,
	}
}

func (ri *ResultItems) Get(key string) interface{} {
	return ri.Items[key]
}

func (ri *ResultItems) Put(key string, value interface{}) {
	ri.Items[key] = value
}

func (ri *ResultItems) IsSkip() bool {
	return ri.Skip
}

func (ri *ResultItems) All() map[string]interface{} {
	return ri.Items
}
