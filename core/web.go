//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/KarnerTh/ts-clean-arch-helper/analyze"
	"github.com/KarnerTh/ts-clean-arch-helper/generation"
)

var c chan bool

func init() {
	c = make(chan bool)
}

func main() {
	js.Global().Set("convert", js.FuncOf(convert))
	<-c
}

func convert(_ js.Value, inputs []js.Value) any {
	input := inputs[0].String()
	analyzer := analyze.NewAnalyzer()
	result, err := analyzer.Analyze(input)
	if err != nil {
		fmt.Println(err)
	}

	generator := generation.NewGenerator()
	output, err := generator.Generate(result)
	if err != nil {
		fmt.Println(err)
	}

	out, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
	}

	return string(out)
}
