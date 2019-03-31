// comparotor.go
package comparator

import (
	"fmt"
)

type Comparator interface {
	Compare(got interface{}, expect interface{}) bool
}

const (
	CompEqual = "eq"
)

func CreateComparator(compType string) Comparator {
	if compType == CompEqual {
		compEq := &ComparatorEq{}
		return compEq
	}

	fmt.Printf("compType = %v\n", compType)

	return nil
}
