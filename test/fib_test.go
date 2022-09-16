package test

import "testing"

//func TestFibList(t *testing.T){
//	a := 1
//	b := 1
//	for i := 0; i < 5; i++ {
//		temp := a
//		a = b
//		b = temp + a
//		t.Log(b)
//	}
//
//}
//
//func TestExchange(t *testing.T){
//	a := 1
//	b := 2
//	t.Log(a, b)
//	a,b = b, a
//	t.Log(a, b)
//}
//
//const (
//	Monday = 1 + iota
//	Tuesday
//	wednesday
//)
//
//func TestConstant(t *testing.T)  {
//	t.Log(Monday, Tuesday)
//}
//
//func TestType(t *testing.T)  {
//	a := "aaaaa"
//	b := 1
//	c := func() {
//		fmt.Println("a")
//	}
//	t.Logf("%T %T %T", a, b, c)
//}

//func TestArr(t *testing.T){
//	b := [2][2]int{
//		{
//			1, 2,
//		},
//		{
//			2, 3,
//		},
//	}
//	t.Log(b)
//}
//
//func TestSlice(t *testing.T)  {
//	slia := []int{1, 2, 3}
//	slib := [][]int{
//		{
//			1, 2,
//		},
//		{
//			2, 3,
//		},
//	}
//	t.Log(slia, slib)
//}


func TestMap(t *testing.T){
	mapA := make(map[int]int, 0)
	t.Log(mapA[1])
}

func TestRaceLock(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name:"test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}