package calc

import (
	"fmt"

	"gopkg.in/Knetic/govaluate.v2"
)

func CalcStr(str string) (float64, error) {
	fmt.Println(str)

	expression, _ := govaluate.NewEvaluableExpression(str)
	out, _ := expression.Evaluate(nil)
	fmt.Println(out)
	return out.(float64), nil
}

func FormatMathOut(out float64) string {
	//if is int
	if out == float64(int(out)) {
		return fmt.Sprintf("%d", int(out))
	}
	return fmt.Sprintf("%f", out)
}
