package compakq

import (
	"time"
)

type TFruitName (string)

type MyFruit struct {
	QItem[MyFruit]
	T TFruitName
}

func (f TFruitName) String() string {
	const (
		Apple  TFruitName = "Apple"
		Banana TFruitName = "Banana"
		Cherry TFruitName = "Cherry"
		Grape  TFruitName = "Grape"
		Mango  TFruitName = "Mango"
		Orange TFruitName = "Orange"
	)

	return map[TFruitName]string{
		Apple:  "Apple",
		Banana: "Banana",
		Cherry: "Cherry",
		Grape:  "Grape",
		Mango:  "Mango",
		Orange: "Orange",
	}[f]
}

func (ctx MyFruit) Index() int {
	return int(time.Now().Unix())
}

func (ctx MyFruit) Key() string {
	return ctx.T.String()
}
