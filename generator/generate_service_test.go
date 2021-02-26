package generator

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"testing"
)

type A struct {
}

func (a A) X() error {
	return nil
}
func TestFun(t *testing.T) {
	x := jen.Func().Params(jen.Id("a").Op("*").Id("A")).Id("X").Params().Params(jen.Id("err").Id("error")).Block(
		jen.Var().Id("foo1").Id("int"),
		jen.Var().Id("foo2").Id("int"),
		jen.Id("RETRY:"),
		jen.Id("bar").Op(":=").Lit(1111),
		jen.Var().Id("foo3").Id("int"),
		jen.Var().Id("foo3").Map(jen.String()).Index().String(),
		jen.List(jen.Id("a"), jen.Id("b")).Op(":=").Id("funcx").Call(),
		jen.Id("a").Op(",").Id("b").Op(":=").Id("funcx").Call(),
		jen.Goto().Id("RETRY"),

		jen.Return(),
	)

	fmt.Println()
	fmt.Printf("%#v", x)
	fmt.Println()
	xx := jen.NewFile("x")
	xx.Add(x)
	xx1 := xx.GoString()

	fmt.Printf("%#v", xx1)

}
