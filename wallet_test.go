// description: keybox 
// 
// @author: xwc1125
// @date: 2020/8/18 0018
package keybox

import (
	"fmt"
	"testing"
)

func TestParseChildKeyPath(t *testing.T) {
	childKeyPath, _ := buildChildKeyPath(Purpose45, 0x80000200, 0x80000003, 0x80000000, 0, 0)
	fmt.Println("childKeyPath", childKeyPath)
	keyPath, err := parseChildKeyPath(childKeyPath)
	if err != nil {
		panic(err)
	}
	fmt.Println(keyPath)
}
