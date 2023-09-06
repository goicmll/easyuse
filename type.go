package easyuse

import "reflect"

// GetMethodSet 获取类型的方法集合名，可调用不一定属于方法集合
// 方法集合决定接口实现
// 参数 Kind 不是 Array, Chan, Map, Pointer, Slice时 返回无元素切片
func GetMethodSetName(a any) []string {
	t := reflect.TypeOf(a)
	ms := make([]string, 0, 4)
	k := t.Kind().String()
	if k != "Array" && k != "Chan" && k != "Map" && k != "Pointer" && k != "Slice" {
		return ms
	}
	elemType := t.Elem()
	n := elemType.NumMethod()
	for i := 0; i < n; i++ {
		ms = append(ms, elemType.Method(i).Name)
	}
	return ms
}
