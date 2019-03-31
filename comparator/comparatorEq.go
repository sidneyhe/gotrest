// comparatorEq.go
package comparator

import (
	"fmt"
	"reflect"
)

type ComparatorEq struct {
}

func (comp *ComparatorEq) Compare(got interface{}, expect interface{}) bool {
	fmt.Println("gotType =", reflect.TypeOf(got), "; expectType=", reflect.TypeOf(expect))
	gotStr := fmt.Sprintf("%v", got)
	expStr := fmt.Sprintf("%v", expect)
	return gotStr == expStr
}
