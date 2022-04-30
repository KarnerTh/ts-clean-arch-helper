package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/KarnerTh/ts-clean-arch-helper/analyze"
	"github.com/KarnerTh/ts-clean-arch-helper/generation"
)

func main() {
	// to test locally
	input := `
export enum EnumWithValuesEntity {
  Draft = "draft",
  Published = "published",
}

export enum EnumWithoutValuesEntity {
  Monthly,
  Yearly,
}

`
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

	out, err := json.Marshal(output.ModelConverter)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(strings.Replace(string(out), `\n`, "\n", -1))
}
