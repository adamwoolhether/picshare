package main

import "fmt"

type Cat struct {}

func (c Cat) Speak() {
	fmt.Println("meow ")
}

type Dog struct{}

func (d Dog) Speak() {
	fmt.Println("woof")
}

type Husky struct {
	// To use embedding, we replace the name with JUST the type.
	// instead of `dog Dog` we just use `Dog`
	//Dog
	Speaker
}

type SpeakerPrefix struct {
	Speaker
}

func (sp SpeakerPrefix) Speak() {
	fmt.Print("Prefix: ")
	sp.Speaker.Speak()
}

type Speaker interface {
	Speak()
}

func main() {
	h := Husky{Dog{}}
	h.Speak() // equal to h.Dog.Speak()
	h = Husky{Cat{}}
	h.Speak() // equal to h.Dog.Speak()
	h = Husky{SpeakerPrefix{Cat{}}}
	h.Speak()
}