package models

type SimpleCounterUint struct {
	Counter uint
}

func (this *SimpleCounterUint) Add(number uint) {
	this.Counter += number
}

func (this *SimpleCounterUint) GetCount() uint {
	return this.Counter
}
