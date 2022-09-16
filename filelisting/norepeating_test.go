package filelisting

import (
	"fmt"
	"testing"
)

func BenchmarkLenghtOfNoRepeatingSubStr(b *testing.B) {
	s := "黑化肥挥发发灰会花飞灰化肥挥发发黑会飞花"
	for i := 0; i < 13; i++ {
		s += s
	}
	ans := 8
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if actual := LenghtOfNoRepeatingSubStr(s); actual != ans {
			b.Errorf("lenghtOfNoRepeatingSubStr() = %v, want %v", actual, ans)
		}
	}
	b.StopTimer()

}

func TestLenghtOfNoRepeatingSubStr(t *testing.T) {
	type args struct{
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "lenC1",
			args: args{
				s: "abcde123sdf23a231",
			},
			want: 8,
		},
		{
			name: "lenC2",
			args: args{
				s: "aceerofwe",
			},
			want: 5,
		},
		{
			name: "lenC3",
			args: args{
				s: "我们大家都是好朋友呀我们的呢",
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LenghtOfNoRepeatingSubStr(tt.args.s); got != tt.want {
				t.Errorf("lenghtOfNoRepeatingSubStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleLenghtOfNoRepeatingSubStr() {
	s := "adfsfw32134123"
	fmt.Println(LenghtOfNoRepeatingSubStr(s))
	s1 := "123451234"
	fmt.Println(LenghtOfNoRepeatingSubStr(s1))

	// Output:
	// 2
	// 5
}