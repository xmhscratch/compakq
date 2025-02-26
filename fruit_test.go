package compakq

import (
	"time"
)

type TFruitName (string)

const (
	_                 = iota
	Apple  TFruitName = "Apple"
	Banana TFruitName = "Banana"
	Cherry TFruitName = "Cherry"
	Grape  TFruitName = "Grape"
	Mango  TFruitName = "Mango"
	Orange TFruitName = "Orange"
)

func (f TFruitName) String() string {
	return map[TFruitName]string{
		Apple:  "Apple",
		Banana: "Banana",
		Cherry: "Cherry",
		Grape:  "Grape",
		Mango:  "Mango",
		Orange: "Orange",
	}[f]
}

type Fruit struct {
	QItem[Fruit]
	T TFruitName
}

func (ctx Fruit) Index() int {
	return int(time.Now().Unix())
}

func (ctx Fruit) Key() string {
	return ctx.T.String()
}
