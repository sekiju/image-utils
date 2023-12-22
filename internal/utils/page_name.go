package utils

import (
	"fmt"
	"strings"
)

type PageName struct {
	pad string
}

func NewPageName(pagesCount int) *PageName {
	count := 999
	if pagesCount > 0 {
		count = pagesCount
	}

	pad := strings.Repeat("0", len(fmt.Sprint(count)))

	return &PageName{pad: pad}
}

func (pn *PageName) GetName(index int) string {
	str := fmt.Sprint(index)
	return pn.pad[:len(pn.pad)-len(str)] + str
}
