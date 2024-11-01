package models

type IPerson interface{}
type Broker struct {
	Name string
} //implements IPerson

type Trader struct {
	Name string
} //implements IPerson

type IHuman interface{}
