/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

//var Str string

type v1 struct {
	int
	string
}

type user22 struct {
	name    string
	address string
	person
	Interest     string
	big          int
	LikeBook     bool `json:"like_book"`
	DontLikePlay bool
}

type student struct {
	id   int
	name string
	age  int
}

func demo(ce []student) {
	//切片是引用传递，是可以改变值的
	ce[1].age = 999
	// ce = append(ce, student{3, "xiaowang", 56})
	// return ce
}

func add1(x, y int) (z int) {
	{
		z := x + y
		return z
	}
}

// 外部引用函数参数局部变量
func add(base int) func(int) int {
	return func(i int) int {
		base += i
		return base
	}
}

// 返回2个函数类型的返回值
func test01(base int) (func(int) int, func(int) int) {
	// 定义2个函数，并返回
	// 相加
	add := func(i int) int {
		base += i
		return base
	}
	// 相减
	sub := func(i int) int {
		base -= i
		return base
	}
	// 返回
	return add, sub
}

func test() {
	x, y := 10, 20

	defer func(i int) {
		println("defer:", i, y) // y 闭包引用
	}(x) // x 被复制

	x += 10
	y += 100
	println("x =", x, "y =", y)
}

func test1(ch chan string) {
	time.Sleep(time.Second * 5)
	ch <- "test1"
}
func test2(ch chan string) {
	time.Sleep(time.Second * 2)
	ch <- "test2"
}

var x int64
var wg sync.WaitGroup
var lock sync.Mutex
var rwlock sync.RWMutex

func add11() {
	for i := 0; i < 5000; i++ {
		//lock.Lock()
		//x = x + 1
		//lock.Unlock()
		atomic.AddInt64(&x, 1)
		//time.Sleep(30 * time.Millisecond)
	}
	wg.Done()
}

func write() {
	// lock.Lock()   // 加互斥锁
	rwlock.Lock() // 加写锁
	x = x + 1
	time.Sleep(10 * time.Millisecond) // 假设读操作耗时10毫秒
	rwlock.Unlock()                   // 解写锁
	// lock.Unlock()                     // 解互斥锁
	wg.Done()
}

func read() {
	// lock.Lock()                  // 加互斥锁
	rwlock.RLock()               // 加读锁
	time.Sleep(time.Millisecond) // 假设读操作耗时1毫秒
	rwlock.RUnlock()             // 解读锁
	// lock.Unlock()                // 解互斥锁
	wg.Done()
}

var m = make(map[string]int)

func get(key string) int {
	return m[key]
}

func set(key string, value int) {
	m[key] = value
}

var emailRe = `\w+@\w+\.com`
var phoneRe = `1[3456789]\d{9}`

func crawlerEmail() {
	resp, err := http.Get("https://tieba.baidu.com/p/6051076813?red_tag=1573533731")
	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	pageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	contents := string(pageBytes)

	re := regexp.MustCompile(emailRe)
	matches := re.FindAllStringSubmatch(contents, -1)
	for k, v := range matches {
		fmt.Println(k, v[0])
	}
}

func crawlerPhone() {
	resp, err := http.Get("http://www.zhaohaowang.com/ ")
	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	pageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	contents := string(pageBytes)

	re := regexp.MustCompile(phoneRe)
	matches := re.FindAllStringSubmatch(contents, -1)
	for k, v := range matches {
		fmt.Println(k, v[0])
	}
}

func HandleError(err error, why string) {
	if err != nil {
		fmt.Println(why, err)
	}
}

// 下载图片，传入的是图片叫什么
func DownloadFile(url string, filename string) (ok bool) {
	resp, err := http.Get(url)
	HandleError(err, "http.get.url")
	defer resp.Body.Close()
	fileBytes, err := io.ReadAll(resp.Body)
	HandleError(err, "resp.body")
	filename = "/Users/wangwulve/qimao/qimao/go-learn/data/img/" + filename
	// 写出数据
	err = os.WriteFile(filename, fileBytes, 0666)
	if err != nil {
		return false
	} else {
		return true
	}
}

// 并发爬思路：
// 1.初始化数据管道
// 2.爬虫写出：26个协程向管道中添加图片链接
// 3.任务统计协程：检查26个任务是否都完成，完成则关闭数据管道
// 4.下载协程：从管道里读取链接并下载

var (
	// 存放图片链接的数据管道
	chanImageUrls chan string
	waitGroup     sync.WaitGroup
	// 用于监控协程
	chanTask chan string
	reImg    = `https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

func crawlerImg() {
	// myTest()
	// DownloadFile("http://i1.shaodiyejin.com/uploads/tu/201909/10242/e5794daf58_4.jpg", "1.jpg")

	// 1.初始化管道
	chanImageUrls = make(chan string, 1000000)
	chanTask = make(chan string, 26)
	// 2.爬虫协程
	for i := 1; i < 27; i++ {
		waitGroup.Add(1)
		go getImgUrls("https://www.bizhizu.cn/shouji/tag-%E5%8F%AF%E7%88%B1/" + strconv.Itoa(i) + ".html")
	}
	// 3.任务统计协程，统计26个任务是否都完成，完成则关闭管道
	waitGroup.Add(1)
	go CheckOK()
	// 4.下载协程：从管道中读取链接并下载
	for i := 0; i < 5; i++ {
		waitGroup.Add(1)
		go DownloadImg()
	}
	waitGroup.Wait()
}

// 下载图片
func DownloadImg() {
	for url := range chanImageUrls {
		filename := GetFilenameFromUrl(url)
		ok := DownloadFile(url, filename)
		if ok {
			fmt.Printf("%s 下载成功\n", filename)
		} else {
			fmt.Printf("%s 下载失败\n", filename)
		}
	}
	waitGroup.Done()
}

// 截取url名字
func GetFilenameFromUrl(url string) (filename string) {
	// 返回最后一个/的位置
	lastIndex := strings.LastIndex(url, "/")
	// 切出来
	filename = url[lastIndex+1:]
	// 时间戳解决重名
	timePrefix := strconv.Itoa(int(time.Now().UnixNano()))
	filename = timePrefix + "_" + filename
	return
}

// 任务统计协程
func CheckOK() {
	var count int
	for {
		url := <-chanTask
		fmt.Printf("%s 完成了爬取任务\n", url)
		count++
		if count == 26 {
			close(chanImageUrls)
			break
		}
	}
	waitGroup.Done()
}

// 爬图片链接到管道
// url是传的整页链接
func getImgUrls(url string) {
	urls := getImgs(url)
	// 遍历切片里所有链接，存入数据管道
	for _, url := range urls {
		chanImageUrls <- url
	}
	// 标识当前协程完成
	// 每完成一个任务，写一条数据
	// 用于监控协程知道已经完成了几个任务
	chanTask <- url
	waitGroup.Done()
}

// 获取当前页图片链接
func getImgs(url string) (urls []string) {
	pageStr := GetPageStr(url)
	re := regexp.MustCompile(reImg)
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Printf("共找到%d条结果\n", len(results))
	for _, result := range results {
		url := result[0]
		urls = append(urls, url)
	}
	return
}

// 抽取根据url获取内容
func GetPageStr(url string) (pageStr string) {
	resp, err := http.Get(url)
	HandleError(err, "http.Get url")
	defer resp.Body.Close()
	// 2.读取页面内容
	pageBytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "ioutil.ReadAll")
	// 字节转字符串
	pageStr = string(pageBytes)
	return pageStr
}

var (
	server = "127.0.0.1:11211"
)

func bufioDemo() {
	reader := bufio.NewReader(os.Stdin) // 从标准输入生成读对象
	fmt.Print("请输入内容：")
	text, _ := reader.ReadString('\n') // 读到换行
	text = strings.TrimSpace(text)
	fmt.Printf("%#v\n", text)
}

func (p person) like() *User {
	fmt.Println("aa")
	return nil
}

type student1 struct {
	Name string
	Age  int
}

func pase_student() {
	m := make(map[string]*student1)
	stus := []student1{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}
	fmt.Printf("%p\n", &stus)
	fmt.Printf("%p\n", &stus[0])
	fmt.Printf("%p\n", &stus[1])
	for _, stu := range stus {
		fmt.Printf("%p\n", &stu)
		m[stu.Name] = &stu
		//stu.Age = stu.Age + 10
	}
	fmt.Println(m)
}

type UserAges struct {
	ages map[string]int
	sync.Mutex
}

func (ua *UserAges) Add(name string, age int) {
	ua.Lock()
	defer ua.Unlock()
	ua.ages[name] = age
}
func (ua *UserAges) Get(name string) int {
	ua.Lock()
	defer ua.Unlock()
	if age, ok := ua.ages[name]; ok {
		return age
	}
	return -1
}

type People11 interface {
	Show()
}

type Student11 struct{}

func (stu *Student11) Show() {
}
func live() People11 {
	var stu *Student11
	return stu
}

func GetValue() int {
	return 1
}

func DeferFunc1(i int) (t int) {
	t = i
	defer func() {
		t += 3
	}()
	return t
}
func DeferFunc2(i int) int {
	t := i
	fmt.Println(t)
	defer func() {
		fmt.Println(t)
		t += 3
		fmt.Println(t)
	}()
	fmt.Println(t)
	return t
}
func DeferFunc3(i int) (t int) {
	fmt.Println(t)
	defer func() {
		fmt.Println(t)
		t += i
		fmt.Println(t)
	}()
	fmt.Println(t)
	return 2
}

func Foo(x interface{}) {
	if x == nil {
		fmt.Println("emptyinterface")
		return
	}
	fmt.Println("non-emptyinterface")
}

//func GetValue1(m map[int]string, id int) (string, bool) {
//	if _, exist := m[id]; exist {
//		return "存在数据", true
//	}
//	return nil, false
//}

const (
	x1 = iota
	y
	z = "zz"
	k
	p = iota
)

//var (
//	size := 1024
//	max_size = size * 2
//)

const cl = 100

var bl = 123

type User1 struct {
}
type MyUser1 User1
type MyUser2 = User1

func (i MyUser1) m1() {
	fmt.Println("MyUser1.m1")
}
func (i User1) m2() {
	fmt.Println("User.m2")
}

type T1 struct {
}

func (t T1) m1() {
	fmt.Println("T1.m1")
}

type T2 = T1

type MyStruct struct {
	T1
	T2
}

var ErrDidNotWork = errors.New("did not work")

func DoTheThing(reallyDoIt bool) (err error) {
	//var result string
	if reallyDoIt {
		//result, err = tryTheThing()
		result, err := tryTheThing()
		if err != nil || result != "it worked" {
			err = ErrDidNotWork
		}
	}
	return err
}
func tryTheThing() (string, error) {
	return "", ErrDidNotWork

}

func test11() []func() {
	var funs []func()
	for i := 0; i < 2; i++ {
		funs = append(funs, func() {
			println(&i, i)
		})
	}
	return funs
}

func test12(x int) (func(), func()) {
	return func() {
			println(x)
			x += 10
		}, func() {
			println(x)
		}
}

func main1() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("fatal")
		}
	}()
	defer func() {
		panic("defer panic")
	}()
	panic("panic")
}

func main() {
	fmt.Println(errors.New("this is a error"))
	err := errors.New("this is a error")
	fmt.Println(errors.WithStack(err))
	//defer func() {
	//	if err := recover();err != nil {
	//		fmt.Println("++++")
	//		f := err.(func() string)
	//		fmt.Println(err)
	//		fmt.Println(f())
	//		fmt.Println(reflect.TypeOf(err).Kind().String())
	//	} else {
	//		fmt.Println("fatal")
	//	}
	//}()
	//defer func() {
	//	panic(func() string {
	//		fmt.Println(222)
	//		return "defer panic"
	//	})
	//}()
	//defer func() {
	//	fmt.Println(111)
	//	panic("err111")
	//}()
	//defer func() {
	//	fmt.Println(222)
	//	panic("err222")
	//}()
	//panic("panic")
	//a, b := test12(100)
	//a()
	//b()
	//funs := test11()
	//for _, f := range funs {
	//	f()
	//}
	//fmt.Println(DoTheThing(true))
	//fmt.Println(DoTheThing(false))
	//s := MyStruct{}
	//s.m1()
	//var i1 MyUser1
	//var i2 MyUser2
	//i1.m1()
	//i2.m2()
	//i1.m2()

	//type MyInt1 int
	//type MyInt2 = int
	//var i int = 9
	//var i1 MyInt1 = i
	//var i2 MyInt2 = i
	//fmt.Println(i1, i2)
	//for i := 0; i < 10; i++ {
	//	loop:
	//		println(i)
	//}
	//goto loop
	//println(&bl, bl)
	//println(&cl, cl)
	//fmt.Println(x1, y, z, k, p)
	//intmap := map[int]string{
	//	1: "a",
	//	2: "bb",
	//	3: "ccc",
	//}
	//v, err := GetValue1(intmap, 3)
	//fmt.Println(v, err)
	//var x *int = nil
	//Foo(x)
	//sn1 := struct {
	//	age  int
	//	name string
	//}{age: 11, name: "qq"}
	//sn2 := struct {
	//	age  int
	//	name string
	//}{age: 11, name: "qq"}
	//if sn1 == sn2 {
	//	fmt.Println("sn1== sn2")
	//}
	//sm1 := struct {
	//	age int
	//	m   map[string]string
	//}{age: 11, m: map[string]string{"a": "1"}}
	//sm2 := struct {
	//	age int
	//	m   map[string]string
	//}{age: 11, m: map[string]string{"a": "1"}}
	//if sm1 == sm2 {
	//	fmt.Println("sm1== sm2")
	//}

	//list := new([]int)
	//list = append(list,1)
	//fmt.Println(list)

	//println(DeferFunc1(1))
	//println(DeferFunc2(1))
	//println(DeferFunc3(1))

	//i := GetValue()
	//switch i.(type) {
	//case int:
	//	println("int")
	//case string:
	//	println("string")
	//case interface{}:
	//	println("interface")
	//default:
	//	println("unknown")
	//}

	//if live() == nil {
	//	fmt.Println("AAAAAAA")
	//} else {
	//	fmt.Println("BBBBBBB")
	//}

	//pase_student()
	//var p uintptr
	//fmt.Println(p)
	//var e error
	//fmt.Println(e)
	//f1 := func(){fmt.Println(22)}
	//fmt.Println(f1)
	//var buf [16]byte
	//os.Stdin.Read(buf[:])
	//fmt.Println(buf)
	//fmt.Println(string(buf[:]))
	//os.Stdin.WriteString(string(buf[:]))

	//log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	//log.Printf("err is:%s\n", "a")
	//logger := log.New(os.Stdout, "<log>", log.Llongfile | log.Lmicroseconds)
	//logger.Println("print a log")
	//log.Fatal("bbbb")
	//log.Panic("dddd")
	//var name string
	//var delay time.Duration
	//age := flag.Int("age", 18, "名字")
	//flag.StringVar(&name, "name", "xiaohong", "名字")
	//flag.DurationVar(&delay, "d", 0, "延迟时间间隔")
	//fmt.Println(age)
	//fmt.Println(*age)
	//fmt.Println(name)
	////flag解析
	//flag.Parse()
	//fmt.Println(*age)
	//fmt.Println(name)
	//fmt.Println(flag.Args())
	//fmt.Println(flag.NArg())
	//fmt.Println(flag.NFlag())
	//fmt.Println(delay)
	//for k, v := range os.Args {
	//	fmt.Println(k, v)
	//}
	//now := time.Now()
	//fmt.Println(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	//fmt.Println(now.Date())
	//fmt.Println(now.Clock())
	//g, _ := now.GobEncode()
	//fmt.Println(string(g))
	//fmt.Println(now.Local())
	//fmt.Println(now.Location())
	//fmt.Println(now.String(), "\n", now.Unix(), "\n", now.UnixNano(), "\n",  now.UnixMilli(), "\n",  now.UnixMicro(), "\n",  now.UTC())
	//fmt.Println(time.Unix(1662712299, 567827))
	//afterHourTime := time.Now().Add(3600 * time.Second)
	//now := time.Now()
	////dura := now.Sub(afterHourTime)
	////fmt.Println(dura.Seconds())
	////fmt.Println(now.Before(afterHourTime))
	//fmt.Println(now.Format("20060102 15:04:05"))
	//fmt.Println(now.Format("2006-01-02 15:04:05"))
	//// 格式化的模板为Go的出生时间2006年1月2号15点04分 Mon Jan
	//// 24小时制
	//fmt.Println(now.Format("2006-01-02 15:04:05.000 Mon Jan"))
	//// 12小时制
	//fmt.Println(now.Format("2006-01-02 03:04:05.000 PM Mon Jan"))
	//fmt.Println(now.Format("2006/01/02 15:04"))
	//fmt.Println(now.Format("15:04 2006/01/02"))
	//fmt.Println(now.Format("2006/01/02"))
	//fmt.Println(time.Now().UTC())
	//loc, _ := time.LoadLocation("Asia/Shanghai")
	//time1, err :=time.ParseInLocation("2006/01/02 15:04:05", "2019/08/04 14:15:20", loc)
	//fmt.Println(time1,err)
	//bufioDemo()

	//fmt.Println(fmt.Errorf("this is a %s", "错误啊！"))
	//fmt.Printf("%v\n", 100)
	//fmt.Printf("%v\n", false)
	//o := struct{ name string }{"枯藤"}
	//fmt.Printf("%v\n", o)
	//fmt.Printf("%#v\n", o)
	//fmt.Printf("%T\n", o)
	//fmt.Printf("100%%\n")
	//var b bool
	//b = true
	//fmt.Printf("%t\n", b)
	//crawlerEmail()
	//crawlerPhone()
	//crawlerImg()

	//for i:=0;i<8;i++{
	//	go func() {
	//		time.Sleep(3 * time.Second)
	//		fmt.Println("aaaa")
	//	}()
	//}
	//
	//fmt.Println(runtime.NumGoroutine())
	//time.Sleep(4 * time.Second)
	//fmt.Println(runtime.NumGoroutine())

	//wg := sync.WaitGroup{}
	//for i := 0; i < 2000; i++ {
	//	wg.Add(1)
	//	go func(n int) {
	//		key := strconv.Itoa(n)
	//		set(key, n)
	//		fmt.Printf("k=:%v,v:=%v\n", key, get(key))
	//		wg.Done()
	//	}(i)
	//}
	//wg.Wait()

	//start := time.Now()
	//for i := 0; i < 10; i++ {
	//	wg.Add(1)
	//	go write()
	//}
	//
	//for i := 0; i < 1000; i++ {
	//	wg.Add(1)
	//	go read()
	//}
	//
	//wg.Wait()
	//end := time.Now()
	//fmt.Println(end.Sub(start))
	//fmt.Println(time.Since(start))

	//f, err := os.Create("trace.out")
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//
	//err = trace.Start(f)
	//if err != nil {
	//	panic(err)
	//}
	//defer trace.Stop()
	//
	//
	//wg.Add(2)
	//go add11()
	//go add11()
	//wg.Wait()
	//fmt.Println(x)
	//fmt.Println(atomic.LoadInt64(&x))
	//fmt.Println("end")

	//// 2个管道
	//output1 := make(chan string)
	//output2 := make(chan string)
	//// 跑2个子协程，写数据
	//go test1(output1)
	//go test2(output2)
	//// 用select监控
	//select {
	//case s1 := <-output1:
	//	fmt.Println("s1=", s1)
	//case s2 := <-output2:
	//	fmt.Println("s2=", s2)
	//}

	//ticker := time.NewTicker(time.Second * 2)
	//for {
	//	fmt.Println(<-ticker.C)
	//	fmt.Println("ticker")
	//}

	//timer := time.NewTimer(time.Second * 2)
	//t1 := time.Now()
	//fmt.Printf("t1:%v\n", t1)
	//t2 := <-timer.C
	//fmt.Printf("t1:%v\n", t2)

	// 3.timer实现延时的功能
	//(1)
	//time.Sleep(time.Second)
	////(2)
	//timer3 := time.NewTimer(2 * time.Second)
	//<-timer3.C
	//fmt.Println("2秒到")
	////(3)
	//<-time.After(2*time.Second)
	//fmt.Println("2秒到")

	//timer4 := time.NewTimer(time.Second * 3)
	//go func() {
	//	<-timer4.C
	//	fmt.Println("定时器执行了")
	//}()
	//
	//b := timer4.Stop()
	//if b {
	//	fmt.Println("定时器关闭了")
	//}
	//for  {
	//
	//}

	// 5.重置定时器
	//timer5 := time.NewTimer(10 * time.Second)
	//timer5.Reset(1 * time.Second)
	//fmt.Println(time.Now())
	//fmt.Println(<-timer5.C)
	//
	//for {
	//}

	//fmt.Println(23%10)
	//fmt.Println(27/10)

	//t1 := time.Now()
	//for i := 0; i < 10000; i++ {
	//
	//}
	//fmt.Println(time.Since(t1))
	//fmt.Println(time.Now().Sub(t1))

	//test()

	//tmp1 := add(10)
	//fmt.Println(tmp1(1), tmp1(2))
	//// 此时tmp1和tmp2不是一个实体了
	//tmp2 := add(100)
	//fmt.Println(tmp2(1), tmp2(2))
	//
	//f1, f2 := test01(10)
	//// base一直是没有消
	//fmt.Println(f1(1), f2(2))
	//// 此时base是9
	//fmt.Println(f1(3), f2(4))

	//sqrtfun := func(x float64) float64 {
	//	return math.Sqrt(x)
	//}
	//
	//fmt.Println(sqrtfun(4))

	//fmt.Println(add(1, 2))

	//for i:=0; i< 10; i++ {
	//	fmt.Println(i)
	//	if i == 3 {
	//		goto label1
	//	}
	//}
	//
	//label1:
	//	fmt.Println("label1")
	//
	//fmt.Println("aa")

	//a := [3]int{0, 1, 2}
	//
	//for i, v := range a { // index、value 都是从复制品中取出。
	//
	//	if i == 0 { // 在修改前，我们先修改原数组。
	//		a[1], a[2] = 999, 999
	//		fmt.Println(a) // 确认修改有效，输出 [0, 999, 999]。
	//	}
	//
	//	a[i] = v + 100 // 使用复制品中取出的 value 修改原数组。
	//
	//}
	//
	//fmt.Println(a)

	//s := []int{1, 2, 3, 4, 5}
	//
	//for i, v := range s { // 复制 struct slice { pointer, len, cap }。
	//
	//	if i == 0 {
	//		s = s[:3]  // 对 slice 的修改，不会影响 range。
	//		s[2] = 100 // 对底层数据的修改。
	//	}
	//
	//	fmt.Println(i, v)
	//}
	//fmt.Println(s)
	//c1 := make(chan int, 20)
	////c2 := make(chan int, 20)
	//c3 := make(chan int, 20)
	//go func() {
	//	for i := 0; i < 20; i++ {
	//		//time.Sleep(time.Second * 1)
	//		c1 <-1
	//	}
	//}()
	//time.Sleep(time.Second * 1)
	//fmt.Println("aaaaa")
	//var i1, i2, i3 int
	//select {
	//case i1 = <-c1:
	//	fmt.Println("c1", i1)
	//case i2 = <-c1:
	//	fmt.Println("c2", i2)
	//case i3 = <-c3:
	//	fmt.Println("c3", i3)
	//	//default:
	//	//	fmt.Println("other")
	//}
	//
	//fmt.Println("end")

	//m1 := make(map[int]int, 0)
	//m2 := make(map[int]struct{}, 0)
	//for i := 0; i < 10; i++ {
	//	m1[i] = i
	//	m2[i] = struct{}{}
	//}
	//
	//fmt.Println(unsafe.Sizeof(m1))
	//fmt.Println(unsafe.Sizeof(m2))
	//a := 1
	//switch a {
	//case 1:
	//	fmt.Println(1)
	//	fallthrough
	//case 2:
	//	fmt.Println(2)
	//}
	//
	//fmt.Printf("a type is %T\n", a)
	//switch interface{}(a).(type) {
	//case int:
	//	fmt.Println("int")
	//case string:
	//	fmt.Println("string")
	//}
	//
	//fmt.Println(1111)
	//var ce []student  //定义一个切片类型的结构体
	//ce = []student{
	//	student{1, "xiaoming", 22},
	//	student{2, "xiaozhang", 33},
	//}
	//fmt.Println(ce)
	//demo(ce)
	//fmt.Println(ce)

	//	ce := make(map[int]student)
	//	ce[1] = student{1, "xiaolizi", 22}
	//	ce[2] = student{2, "wang", 23}
	//	fmt.Println(ce)
	//	delete(ce, 2)
	//	fmt.Println(ce)
	//
	//	u1 := user22{
	//		name: "a",
	//		address: "shanghai",
	//		person: person{
	//			name: "b",
	//			age: 10,
	//			height: 188,
	//		},
	//		Interest: "girl",
	//		big: 33,
	//		LikeBook: true,
	//		DontLikePlay: false,
	//	}
	//	//var u1 user22
	//	//u1.name = "bb"
	//	data, err := json.Marshal(u1)
	//
	//fmt.Println(u1)
	//fmt.Printf("data is %s", data)
	//fmt.Println(err)

	//t := v1{
	//	18,
	//	"aa",
	//}
	//t.string = "ab"
	//t.int = 19
	//fmt.Println(t)

	//m1 := make(map[string]string, 9)
	//fmt.Println(m1)
	//for i := 0; i < 20; i++ {
	//	str := fmt.Sprintf("%d", i)
	//	m1[str] = str
	//}
	//a := m1["0"]
	//b, ok := m1["0"]
	//fmt.Println(a, b, ok)
	//fmt.Println(uintptr(1))
	//str := "\nabcd\n\n\n\ni2342sdf0923 \newr\n\n\nwerio\n\nfiwoerwe\n'we0sdfs\nsdfew\n"
	//strRe := `\n{2,}`
	//re := regexp.MustCompile(strRe)
	//str = re.ReplaceAllString(str, "\n")
	//str = strings.TrimSpace(str)
	//fmt.Println(str)
	//url1 := "https://audio-api.mengchang.com/qi-mao/detail?key=20000812&albumid={album_id}&timestamp={timestamp}&sign={sign}"
	////urlParse, _ := url.Parse(url1)
	////fmt.Println(urlParse)
	////values, _ := url.ParseQuery(urlParse.RawQuery)
	////fmt.Println(values)
	////fmt.Println(values["key"][0])
	//
	//urlParse, err := url.Parse(url1)
	//if err != nil {
	//	zap.L().Error("url parse fail,", zap.Error(err))
	//}
	//urlQuerys, err := url.ParseQuery(urlParse.RawQuery)
	//if err != nil {
	//	zap.L().Error("url ParseQuery fail,", zap.Error(err))
	//}
	//
	//var urlKeys []string
	//for key, _ := range urlQuerys {
	//	if key != "sign" {
	//		urlKeys = append(urlKeys, key)
	//	}
	//}
	//
	//sort.Strings(urlKeys)
	//var querySli []string
	//for _, v := range urlKeys {
	//	valueSli := urlQuerys[v]
	//	if len(valueSli) > 1 {
	//		//同名有多个
	//		for key1, val1 := range valueSli {
	//			querySli = append(querySli, fmt.Sprintf("%s[%d]=%s", v, key1, val1))
	//		}
	//	} else {
	//		querySli = append(querySli, fmt.Sprintf("%s=%s", v, valueSli[0]))
	//	}
	//}
	//signParams := strings.Join(querySli, "")
	//h := md5.New()
	//_, _ = io.WriteString(h, signParams)
	//sign := fmt.Sprintf("%x", h.Sum(nil))
	//fmt.Println(sign)

	//var mp = make(map[int]int, 1024)
	//fmt.Println(mp)
	//var url string
	//clientId := "NGY3MWFhNjYtMjI5My0xMWVkLWJkZGUtZmExNjNlNWMyMDI1"
	//clientSecret := "NjQ4YWJhZjItMmM2YS0zNjEyLTkwYjctYTRlNzVkZWM2MmQ2"
	//domain := "https://open.staging.qtfm.cn"
	////domain := "https://api.open.qtfm.cn"
	//
	//client := http.Client{}
	//url = domain + "/media/v7/categories"

	//url = fmt.Sprintf(domain + "/media/v7/categories" +
	//	"?device_id=%s" +
	//	"&access_token=%s&" +
	//	"user_id=%s" +
	//	"&coop_open_id=%s" +
	//	"&device_os=%s" +
	//	"&device_os_version=%s" +
	//	"&app_version=%s" +
	//	"&device_model=%s",
	//	"",
	//	"",
	//	"",
	//	"",
	//	"",
	//	"",
	//	"",
	//	"",
	//)

	//req, _ := http.NewRequest("GET", url, nil)
	////req.Header.Set("QT-Sign", )
	//
	//resp, err := client.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(resp)
	//arr1 := []int{1,2,3,4,5,6,7,8}
	//sli1 := arr1[1:5]
	//fmt.Println(sli1, len(sli1), cap(sli1))
	//sli2 := arr1[1:3:3]
	//fmt.Println(sli2, len(sli2), cap(sli2))

	//sli1 := []int{1, 2, 3}
	//fmt.Println(len(sli1), cap(sli1))
	//sli1 = append(sli1, []int{4,5,6,7}...)
	//fmt.Println(sli1, len(sli1), cap(sli1))
	//a := 2
	//var b int32 = 3
	//var c int64 = 4
	//fmt.Println(unsafe.Sizeof(a), unsafe.Alignof(a))
	//fmt.Println(unsafe.Sizeof(b), unsafe.Alignof(b))
	//fmt.Println(unsafe.Sizeof(c), unsafe.Alignof(c))
	//
	//fmt.Println(unsafe.Sizeof(1))
	//fmt.Println(unsafe.Alignof(int(1))) // 8 -- min(8,8)
	//fmt.Println(unsafe.Alignof(int32(1))) // 4 -- min (8,4)
	//fmt.Println(unsafe.Alignof(int64(1))) // 8 -- min (8,8)
	//fmt.Println(unsafe.Alignof(complex128(1)))

	//切片扩容原理
	//sli1 := make([]int64, 10)
	////fmt.Println(unsafe.Sizeof(sli1), unsafe.Alignof(sli1))
	////sli1 = append(sli1, make([]int, 11)...
	//a := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	//sli1 = append(sli1, a...)
	////
	//fmt.Println(sli1, len(sli1), cap(sli1))
	//
	//newcap := 10
	//doublecap := 20
	//cap := 21
	//newcap = 21
	//10 * 8
	//21 * 8
	//21 * 8
	//21 * 8 = 168
	//21   12    176  22
	//
	//23
	//newcap := 10
	//doublecap := 20
	//cap := 33
	//newcap = 33
	//10 * 8
	//33 * 8
	//33 * 8  / 8  = 33
	//33 => 19 => 288 / 8 = 36

	//map1 := make(map[int]int, 3)
	//map1[1] = 1
	//map1[2] = 1
	//map1[3] = 1
	//fmt.Println(len(map1))
	//map1[4] = 1
	//fmt.Println(map1, len(map1))

	//var mapSlice = make([]map[string]string, 3)
	//fmt.Println(mapSlice)
	//for index, value := range mapSlice {
	//	fmt.Printf("index:%d value:%v\n", index, value)
	//}
	//fmt.Println("after init")
	//// 对切片中的map元素进行初始化
	//mapSlice[0] = make(map[string]string, 10)
	//mapSlice[0]["name"] = "王五"
	//mapSlice[0]["password"] = "123456"
	//mapSlice[0]["address"] = "红旗大街"
	//for index, value := range mapSlice {
	//	fmt.Printf("index:%d value:%v\n", index, value)
	//}

	//var sliceMap = make(map[string][]string, 3)
	//fmt.Println(sliceMap)
	//fmt.Println("after init")
	//key := "中国"
	//value, ok := sliceMap[key]
	//if !ok {
	//	value = make([]string, 0, 2)
	//}
	//value = append(value, "北京", "上海")
	//sliceMap[key] = value
	//fmt.Println(sliceMap)

	//rand.Seed(time.Now().UnixNano())
	//
	//scoreMap := make(map[string]int, 100)
	//
	//for i := 0; i < 100; i++ {
	//	key := fmt.Sprintf("%02d", i)
	//	//value := rand.Int()
	//	value := rand.Intn(100)
	//	scoreMap[key] = value
	//}
	//
	////取出map中所有的key然后存入切片keys
	//keys := make([]string, 0, 100)
	//for k := range scoreMap {
	//	keys = append(keys, k)
	//}
	//fmt.Println(keys)
	//sort.Strings(keys)
	//isSort :=sort.StringsAreSorted(keys)
	//fmt.Println(keys, isSort)
	//
	//for k, v := range keys {
	//	fmt.Println(k, v, scoreMap[v])
	//}

	//fmt.Println(scoreMap)
	//for k, v := range scoreMap {
	//	fmt.Println(k, v)
	//}

	//ptr := new(map[int]int)
	//fmt.Println(ptr)
	//fmt.Printf("%p", ptr)
	//指针变量仅声明不初始化无法赋值，因为初始化才会分配内存空间
	//var a *int
	//fmt.Println(a)
	//*a = 100
	//fmt.Println(a)

	//arr1 := []int{1, 2, 3, 4}
	//arr2 := []int{5, 6, 7, 8, 9}
	//arr3 := make([]int, 6)
	//copy(arr2, arr1)
	//copy(arr3, arr1)
	//fmt.Println(arr1, arr2)
	//fmt.Println(arr1, arr3)

	//byteSli := make([]byte, 3)
	//n := copy(byteSli, "abcdefg")
	//fmt.Println(n, byteSli, string(byteSli))

	//改变切片影响原数组
	//arr1 := []int{1, 2, 3, 4}
	//sli1 := arr1[1:3]
	//sli1[0] += 1
	//sli1[1] = 33
	////sli1[2] = 0
	//sli1 = append(sli1,0)
	//fmt.Println(arr1)
	//fmt.Println(sli1)

	//nil切片和空切片
	//var a []int
	//silce1 := make([]int , 0)
	//slice2 := []int{}
	//if slice2 == nil {
	//	fmt.Println("nil slice")
	//} else {
	//	fmt.Println("not nil slice")
	//}

	//n := initValue()
	//fmt.Println(*n/2)
	//fmt.Println(*n)
	//arr1 := []int{1,2,3,4,5,6}
	//fmt.Println(&arr1)
	//fmt.Println(&arr1[1])
	//fmt.Printf("%p\n", &arr1)
	//sli1 := arr1[1:]
	//fmt.Printf("%p\n", &sli1)
	////fmt.Println(&sli1)
	//fmt.Println(&sli1[0])
	//ptr := unsafe.Pointer(&sli1[0])
	//fmt.Println(ptr)
	//fmt.Println(&sli1[1])
	//fmt.Println(sli1[0])

	//a := 1
	//fmt.Println(&a)
	//fmt.Printf("%x",&a)
	//f, _ := os.OpenFile("cpu.profile", os.O_CREATE|os.O_RDWR, 0755)
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()
	//f1, _ := os.OpenFile("mem.profile", os.O_CREATE|os.O_RDWR, 0755)
	//pprof.WriteHeapProfile(f1)

	//var s Student
	//t := reflect.TypeOf(s)
	//for i := 0; i < t.NumField(); i++{
	//	fmt.Println(t.Field(i).Tag.Get("json"))
	//	fmt.Println(t.Field(i).Tag.Get("xorm"))
	//}

	//u := User{1, "5lmh.com", 20, "", ""}
	//v := reflect.ValueOf(u)
	//// 获取方法
	//m := v.MethodByName("Hello")
	//// 构建一些参数
	//args := []reflect.Value{reflect.ValueOf("6666")}
	//// 没参数的情况下：var args2 []reflect.Value
	//// 调用方法，需要传入方法的参数
	//m.Call(args)

	//var v Boy
	//fmt.Printf("%v\n", v)
	//fmt.Printf("%+v\n", v)
	//fmt.Printf("%#v\n", v)

	//m := Boy{User{1, "zs", 20, "aaa", "bbb"}, "bj"}
	//t := reflect.TypeOf(m)
	//fmt.Println(t)
	//// Anonymous：匿名
	//fmt.Printf("%#v\n", t.Field(0))
	//// 值信息
	//fmt.Printf("%#v\n", reflect.ValueOf(m).Field(0))

	//u := User{1, "zs", 20, "beijing", "yellow"}
	//Poni(u)
	//u := User{1, " xiaohong", 20, " shanghai", " black"}
	//SetValue(&u)
	//fmt.Println(u)

	//urls := "https://www.baidu.com/user/list/index.php?user_ids=1,2,3,4&timestamp=16043290234&source=aaa&where=list.json"
	//parseUlr, err := url.ParseRequestURI(urls)
	//fmt.Println(parseUlr.Path, err)
	//fmt.Println(testDefer())
	//fmt.Println(fmt.Sprintf("%04s", strconv.FormatInt(2, 2)))

	//	var albumInfo AlbumInfo
	//	albumInfo.Data.AlbumId = 81
	//	albumInfo.Data.Title = " 测试 专辑1 "
	//	albumInfo.Data.Author = "作者1,作者2"
	//	albumInfo.Data.Anchor = "主播1,主播2,主播3"
	//	albumInfo.Data.Category1Name = "游戏"
	//	albumInfo.Data.Category2Name = "娱乐明星"
	//	albumInfo.Data.Intro = "这里是测试专辑1的简介哦！"
	//	albumInfo.Data.OverStatus = 1
	//	albumInfo.Data.ImageLink = "https://contenthub-test.oss-cn-shanghai.aliyuncs.com/bookimg/public/images/album_cover/1g9h9dps4zXHfQTe3KZs3SKXxZeZD.png"
	//	albumInfo.Data.RecordCompany = "幻想七猫有声"
	//	albumInfo.Data.BookTitle = "测试专辑1文学作品名"
	//	albumInfo.Data.Tag = "游戏,热血,竞技"
	//	albumInfo.Data.AlbumType = 1
	////fmt.Println(reflect.ValueOf(albumInfo.Data.Intro))
	////fmt.Println(reflect.TypeOf(albumInfo.Data.Intro))
	//	//字符串去除前后空格 @todo
	//	trimStruct(albumInfo.Data)
	//	fmt.Println(albumInfo)
	//fmt.Println(time.Now().Nanosecond())
	//fmt.Println(FindDomiantColor("https://contenthub-test.oss-cn-shanghai.aliyuncs.com/bookimg/public/images/album_cover/1g9h9dps4zXHfQTe3KZs3SKXxZeZD.png"))
	//fmt.Println(FindDomiantColor("1g9h9dps4zXHfQTe3KZs3SKXxZeZD.png"))
	// Step 1: Load the image
	//img, err := loadImage("1g9h9dps4zXHfQTe3KZs3SKXxZeZD.png")
	//if err != nil {
	//	log.Fatal("Failed to load image", err)
	//}
	//
	//// Step 2: Process it
	//colours, err := prominentcolor.Kmeans(img)
	//if err != nil {
	//	log.Fatal("Failed to process image", err)
	//}
	//
	//fmt.Println("Dominant colours:")
	//for _, colour := range colours {
	//	fmt.Println("#" + colour.AsString())
	//}

	//ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancelFn()
	//
	//fileReader, err := os.Open("1.m4a")
	//if err != nil {
	//	log.Panicf("Error opening test file: %v", err)
	//}
	//
	//data, err := ffprobe.ProbeReader(ctx, fileReader)
	//if err != nil {
	//	log.Panicf("Error getting data: %v", err)
	//}
	//fmt.Println(data)
	//ExampleShowProgress("1.m4a", "005康震讲诗仙李白-03身世家世—李白为何做倒插门女婿？.mp3")
	//	data, err := ffmpeg.Probe("/Users/wangwulve/qimao/qimao/go-learn/1.m4a", nil)
	//	duration, err := probeOutputDuration(data)
	//fmt.Println(err,duration)
	//	err := ffmpeg.Input("1.m4a")
	//	fmt.Println(err)
	//tag, err := id3v2.Open("各个格式/1.m4a", id3v2.Options{Parse: true})
	//if err != nil {
	//	log.Fatal("Error while opening mp3 file: ", err)
	//}
	//defer tag.Close()
	//
	//// Read tags
	//fmt.Println(tag.Artist())
	//fmt.Println(tag.Title())
	//
	//// Set tags
	//tag.SetArtist("Aphex Twin")
	//tag.SetTitle("Xtal")
	//
	//comment := id3v2.CommentFrame{
	//	Encoding:    id3v2.EncodingUTF8,
	//	Language:    "eng",
	//	Description: "My opinion",
	//	Text:        "I like this song!",
	//}
	//tag.AddCommentFrame(comment)
	//
	//// Write tag to file.mp3
	//if err = tag.Save(); err != nil {
	//	log.Fatal("Error while saving a tag: ", err)
	//}
	//测试init加载顺序
	//test.Echo()
	//maze.Echo()

	//var arr1 [5]int
	//printArr(&arr1)
	//fmt.Println(arr1)
	//arr2 := [...]int{2, 4, 6, 8, 10}
	//printArr(&arr2)
	//fmt.Println(arr2)

	//username := flag.String("name", "f1ery", "the user name.")
	//flag.Parse()
	//fmt.Printf("param is %s!\n", *username)
	//test := new([3]int)
	//fmt.Println(test)
	//flag.Parse()
	//lib.Hello(name)
	//lib.Hello(name)
	//Hello(*username)
	//fmt.Println(Str)

	//container := map[int]string{0: "zero", 1: "one", 2:"two"}
	//value, ok := interface{}(container).(map[int]string)
	//fmt.Println(value)
	//fmt.Println(ok)
	//switch t := interface{}(container).(type) {
	//case []string:
	//	fmt.Println(t)
	//case map[string]string:
	//	fmt.Println(t)
	//default:
	//	fmt.Println("unsupported container type: %T", container)
	//}

	//fmt.Println(float32(1))
	//fmt.Println(string(-1))

	//a := []byte{'\xe4','\xbd','\xe5'}
	//fmt.Println(a)

	//for i := 0; i < 100; i++ {
	//	fmt.Printf("main route num %d\n", i)
	//	go func() {
	//		fmt.Printf("goroutine num %d\n", runtime.NumGoroutine())
	//		fmt.Println(i)
	//
	//		fmt.Printf("sleep %d 秒后\n", i)
	//		fmt.Println(i+1)
	//	}()
	//}
	//time.Sleep(2)
	//cnt := strings.Count("abcde", "")
	//fmt.Println(cnt)

	//sli := []int{1, 2, 3, 4, 5, 6, 7, 8}
	//fmt.Println(len(sli))
	//fmt.Println(cap(sli))
	//sli1 := sli[5:6]
	//fmt.Println(sli1)
	//fmt.Println(len(sli1))
	//fmt.Println(cap(sli1))
	//sli := make([]int, 0)
	//i := 1
	//for ; i < 1024; i++ {
	//	sli = append(sli, i)
	//	if i % 2 == 0 {
	//		fmt.Println(i)
	//		fmt.Println(len(sli))
	//		fmt.Println(cap(sli))
	//	}
	//}
	//fmt.Println(len(sli))
	//fmt.Println(cap(sli))
	//i = i + 1
	//sli = append(sli, i)
	//fmt.Println(len(sli))
	//fmt.Println(cap(sli))
	//i = i + 1
	//sli = append(sli, i)
	//fmt.Println(len(sli))
	//fmt.Println(cap(sli))
	//i = i + 1
	//sli = append(sli, i)
	//fmt.Println(len(sli))
	//fmt.Println(cap(sli))

	//var testMap = map[interface{}]string{
	//	1: "a",
	//	"bb": "bb1",
	//	3: "",
	//	[]int{1,2}: "33",
	//}
	//fmt.Println(testMap)
	//var a = [2]int{1,2}
	//b := &a
	////var b = [2]int{1,2}
	////if a == b {
	////	fmt.Println("equal")
	////} else {
	////	fmt.Println("not equal")
	////}
	//
	//fmt.Println(&a)
	//fmt.Println(&b)
	//fmt.Println(&a[0])
	//fmt.Println(&a[1])
	//fmt.Println(&b[0])
	//fmt.Println(&b[1])

	//var a1 = []int{1,2}
	//var b1 = []int{1,2}
	//if a1 == b1 {
	//	fmt.Println("equal")
	//} else {
	//	fmt.Println("not equal")
	//}

	//testMap := map[int]int{1:1,2:2,3:3,4:4,5:5}
	//for k, v := range testMap{
	//	fmt.Println(k)
	//	fmt.Println(v)
	//	fmt.Println(&k)
	//	fmt.Println(&v)
	//}

	//testMap1 := make(map[int]int, 0)
	//go func() {
	//	for i := 0; i < 5; i++ {
	//		testMap1[i] = i
	//	}
	//}()
	//go func() {
	//	for i := 0; i < 5; i++ {
	//		testMap1[i] = i
	//	}
	//}()

	//for i := 0; i < 50; i++ {
	//	go func() {
	//		testMap1[0] = 1
	//	}()
	//}
	//fmt.Println(testMap1)

	//data := &testSendParam{
	//	a: 1,
	//	b: 2,
	//}
	//data := [2]int{1,2}
	//data = updateData(data)
	//fmt.Println(data)

	//
	//array := [3]int{7,8,9}
	//fmt.Printf("main ap brfore: len: %d cap:%d data:%+v\n", len(array), cap(array), array)
	//ap(array)
	//fmt.Printf("main ap after: len: %d cap:%d data:%+v\n", len(array), cap(array), array)

	//b := make([]int, 4)
	//n := copy(b, a)
	//fmt.Println(n, b, &a[0], &b[0])
	//a := []int{1, 2, 3}
	//b := a
	//b[0] = 4
	//c := a
	//c = append(c, 5)
	//d := a
	//d = append(d, 6)
	//d[1] = 7
	//fmt.Println(a, b, c, d)

	//chan1 := make(chan int, 3)
	//
	//go func() {
	//	for i := 0; i < 10; i++ {
	//		chan1 <- i
	//	}
	//	close(chan1)
	//}()
	//
	//for v := range chan1 {
	//	fmt.Println(v)
	//}

	//	ch := make(chan int)
	//	c := 0
	//	stopCh := make(chan bool)
	//
	//	go Chann(ch, stopCh)
	//
	//	for {
	//		select {
	//		case c = <-ch:
	//			fmt.Println("Receive", c)
	//			fmt.Println("channel")
	//		case s := <-ch:
	//			fmt.Println("Receive", s)
	//		case _ = <-stopCh:
	//			goto end
	//		}
	//	}
	//end:
	//	fmt.Println("aaaa")

	//x, y := 0, 1
	//c := make(chan int, 2)
	//sli := make([]int, 0)
	//quit := make(chan int, 1)
	//
	//go func() {
	//
	//}()
	//
	//for  {
	//	select {
	//	case c <- x:
	//		x, y = y, x+y
	//		sli = append(sli, x, y)
	//	case <- quit:
	//		break
	//	}
	//	fmt.Println(sli)
	//}
	//go func() {
	//	t := time.NewTicker(2 * time.Second)
	//	for range t.C {
	//		fmt.Println(2)
	//	}
	//}()
	//
	//for range time.Tick(2 * time.Second) {
	//	fmt.Println(1)
	//}

	//bufChan := make(chan int)
	//
	//go func() {
	//	for {
	//		bufChan <- 1
	//		time.Sleep(time.Second)
	//	}
	//}()
	//
	//go func() {
	//	for {
	//		fmt.Println(<-bufChan)
	//	}
	//}()
	//
	////select {}
	//for  {
	//
	//}
	//c := make(chan int, 3)
	//go func() {
	//	for i := 0; i < 10; i++ {
	//		c <-i
	//		//time.Sleep(2 * time.Second)
	//	}
	//	close(c)
	//}()
	//
	//chanFor:
	//for  {
	//	select {
	//	case a, ok := <-c:
	//		if ok {
	//			fmt.Println(a)
	//		} else {
	//			fmt.Println("close")
	//			//goto closeing
	//			break chanFor
	//		}
	//	}
	//}
	////closeing:
	//	fmt.Println("end")

	//rand.Seed(time.Now().Unix())
	//fmt.Println(rand.Int())
	//fmt.Println(rand.Intn(4))

	//	intChan := make(chan int, 1)
	//	// 一秒后关闭通道。
	//	//time.AfterFunc(time.Second, func() {
	//	//	close(intChan)
	//	//})
	//	intChan <-1
	//	select {
	//	case _, ok := <-intChan:
	//		if !ok {
	//			fmt.Println("The candidate case is closed.")
	//			break
	//		}
	//		fmt.Println("The candidate case is selected.")
	//
	//		break
	//	}
	//close(intChan)
	//	_, ok := <-intChan
	//	fmt.Println(ok)

	//var p Printer
	//p = printStd
	//n, err := p("123")
	//fmt.Println(n, err)

	//complexArray1 := [3][]string{
	//	[]string{"d", "e", "f"},
	//	[]string{"g", "h", "i"},
	//	[]string{"j", "k", "l"},
	//}
	//fmt.Printf("The array: %v\n", complexArray1)
	//array2 := modifyArray(complexArray1)
	//fmt.Printf("The modified array: %v\n", array2)
	//fmt.Printf("The original array: %v\n", complexArray1)

	//cat := New("little pig", "American Shorthair", "cat")
	//cat.SetName("monster") // (&cat).SetName("monster")
	//fmt.Printf("The cat: %s\n", cat)
	//
	//cat.SetNameOfCopy("little pig")
	//fmt.Printf("The cat: %s\n", cat)
	//
	//type Pet interface {
	//	SetName(name string)
	//	Name() string
	//	Category() string
	//	ScientificName() string
	//}
	//
	//_, ok := interface{}(cat).(Pet)
	//fmt.Printf("Cat implements interface Pet: %v\n", ok)
	//_, ok = interface{}(&cat).(Pet)
	//fmt.Printf("*Cat implements interface Pet: %v\n", ok)

	//a := animal{
	//	atype: "a",
	//}
	//a.String1()
	//fmt.Println(a)
	//a.String2()
	//fmt.Println(a)

	//var a int = 1
	//b, c := interface{}(a).(int32)
	//fmt.Println(b, c)
	//sli := []int{1,2,3,4}
	//fmt.Printf("%p\n", sli)
	//fmt.Println( &sli)
	//fmt.Printf("%p\n", sli[0])
	//fmt.Printf("%p\n", &sli[0])

	//var wg sync.WaitGroup
	//wg.Add(3)
	//var a,b,c int64
	//sli :=[]int{1,2,3}
	//for _, v := range sli{
	//	go func(v int) {
	//		defer wg.Done()
	//		fmt.Println(v)
	//		count := int64(33)
	//		if v == 1 {
	//			a = count
	//		} else if v == 2 {
	//			b = count
	//		} else if v == 3 {
	//			c = count
	//		}
	//		//fmt.Println(&a,&b,&c)
	//		//time.Sleep(5 * time.Second)
	//	}(v)
	//}
	//wg.Wait()
	//fmt.Println(a, b,c)
	//fmt.Println(&a,&b,&c)
	//local, _ := time.LoadLocation("Local")
	//t, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02")+" 00:00:00", local)
	//monthStr := t.AddDate(0, -1, 0).Format("01")
	//fmt.Println(monthStr)

	//timeNow := time.Now()
	//thisMonth := time.Date(timeNow.Year(), timeNow.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, -1, 0).Format("01")
	//time.Now().With(timeNow).BeginningOfMonth().AddDate(0, -1, 0)
	//fmt.Println(thisMonth)
	//var name string
	//flag.StringVar(&name, "name", "xiaoming", "xiaoming is a boy")
	//flag.Parse()
	//fmt.Println(name, &name)
	//var name string
	//name = "1"
	//fmt.Println(&name)
	//GetString()
	//fmt.Printf(string(-1))

	//type myInt int
	//var myIntVal myInt
	//var intVal int
	//myIntVal = 1
	//intVal = 1
	//exIntVal := int(myIntVal)
	//if exIntVal == intVal {
	//	fmt.Println("equal")
	//} else {
	//	fmt.Println("not equal")
	//}

	//slia := make([]int, 0)
	//for i := 0; i < 1028; i++ {
	//	fmt.Println(cap(slia))
	//	//fmt.Println(&slia[i])
	//	slia = append(slia, i)
	//}
	//fmt.Println(&slia)
	//
	//var l list.List

	//intChan := make(chan int, 1)
	//// 一秒后关闭通道。
	//time.AfterFunc(time.Second, func() {
	//	intChan <- 1
	//})
	//time.AfterFunc(5*time.Second, func() {
	//	close(intChan)
	//})
	//for  {
	//	select {
	//	case _, ok := <-intChan:
	//		if !ok {
	//			fmt.Println("The candidate case is closed.")
	//			break
	//		}
	//		fmt.Println("The candidate case is selected.")
	//	}
	//}
	//

	//for i := 0; i < 10; i++ {
	//	fmt.Println(&i)
	//	i := i
	//	fmt.Println(&i)
	//	//go func() {
	//	//	fmt.Println(i)
	//	//}()
	//}
	//for {
	//
	//}

	//var count uint32
	//trigger := func(i uint32, fn func()) {
	//	for {
	//		if n := atomic.LoadUint32(&count); n == i {
	//			fn()
	//			atomic.AddUint32(&count, 1)
	//			break
	//		}
	//		time.Sleep(time.Nanosecond)
	//	}
	//}
	//
	//for i := uint32(0); i < 10; i++ {
	//	go func(i uint32) {
	//		fn := func() {
	//			fmt.Println(i)
	//		}
	//		trigger(i, fn)
	//	}(i)
	//}
	//trigger(10, func() {})

	//numbers2 := []int{1,2,3,4,5,6}
	//numbers3 := numbers2
	//maxIndex2 := len(numbers2) - 1
	//for i, e := range numbers2 {
	//	if i == maxIndex2 {
	//		numbers3[0] += e
	//	} else {
	//		numbers3[i+1] += e
	//	}
	//}
	//fmt.Println(numbers2)
	//fmt.Println(numbers3)

	//slia := []int{1,2,3}
	//slib := slia
	//for key, val := range slia {
	//	slib[key] = val + 1
	//}
	//fmt.Println(slia)
	//fmt.Println(slib)

	//slia := []int{3,4,5,6}
	//switch slia[0] {
	//case 0, 1:
	//	fmt.Println(1)
	//case 11, 2:
	//	fmt.Println(2)
	//default:
	//	fmt.Println(10)
	//}

	//var a int
	//a = 1
	//switch interface{}(a).(type) {
	//case int:
	//	fmt.Println(1)
	//case string:
	//	fmt.Println(2)
	//}

	//var err error
	//err = errors.New("empty 11")
	//fmt.Printf("err is %s", err)
	//defer func() {
	//	//recover()
	//	fmt.Println("func defer1")
	//}()
	//defer func() {
	//	//recover()
	//	fmt.Println("func defer2")
	//}()
	//defer func() {
	//	//recover()
	//	fmt.Println("func defer3")
	//}()

	//defer fmt.Println("first")
	//
	//for i := 0; i < 3; i++ {
	//	defer func(i int) {
	//		fmt.Println(i)
	//	}(i)
	//}
	//
	//defer fmt.Println("last")
	//defer func() {
	//	//if p := recover(); p != nil {
	//	//	fmt.Println("incorrect")
	//	//	fmt.Println(p)
	//	//} else {
	//	//	fmt.Println("correct")
	//	//	fmt.Println(p)
	//	//}
	//	panic("bbb")
	//}()
	////slia := []int{1,2,3,4}
	//panic("aaa")
	////fmt.Println(slia[4])
	//fmt.Println("handle panic")

	//slia := []int{1, 2, 3, 4}
	//var m sync.Mutex
	////m.Lock()
	//for k, v := range slia {
	//	slia[k] = v + 1
	//}
	////m.Lock()
	//fmt.Println(slia)
	//m.Unlock()

	//mapa := make(map[int]string, 0)
	//mapa[0] = "a"
	//mapa[1] = "b"
	//mapa[2] = "c"
	//fmt.Println(mapa)
	//delete(mapa, 1)
	//fmt.Println(mapa)

	//var mailbox uint8
	//var lock sync.RWMutex
	//sendCond := sync.NewCond(&lock)
	//recvCond := sync.NewCond(lock.RLocker())
	//
	//lock.Lock()
	//for mailbox == 1 {
	//	sendCond.Wait()
	//}
	//mailbox = 1
	//lock.Unlock()
	//recvCond.Signal()

	//var wg sync.WaitGroup
	//wg.Add(1)
	//for i := 0; i < 2; i++ {
	//	go func(i int) {
	//		defer wg.Done()
	//		fmt.Println(i)
	//	}(i)
	//}
	//wg.Wait()
	//fmt.Println("end")

	//var once sync.Once
	//var ctx context.Context
	//ctx, cancel := context.WithCancel(context.Background())
	//cancel()

	//ch1 := make(chan int, 2)
	//
	//go func() {
	//	a := <-ch1
	//	fmt.Println(a)
	//}()
	//
	//ch1 <- 1
	//ch1 <- 2

	//var c context.Context

	//fmt.Println(reflect.TypeOf("1111"))

	//str := "Go 爱好者 "
	//fmt.Printf("The string: %q\n", str)
	//fmt.Printf("  => runes(char): %q\n", []rune(str))
	//fmt.Printf("  => runes(hex): %x\n", []rune(str))
	//fmt.Printf("  => bytes(hex): [% x]\n", []byte(str))

	//	str1 := "we 好孩子"
	//	fmt.Println(len(str1))
	//	fmt.Println(&str1)
	//	for k, v := range str1 {
	//		fmt.Printf("%d => %q\n", k, v)
	//	}
	//
	//fmt.Println("---------------------")
	//
	//	str2 := []rune(str1)
	//	fmt.Println(len(str2))
	//	fmt.Println(&str2)
	//	for k2, v2 := range str2 {
	//		fmt.Printf("%d => %q\n", k2, v2)
	//	}
	//
	//	str3 := utf8.RuneCountInString(str1)
	//	fmt.Println(str3)
	//
	//	fmt.Println(reflect.TypeOf(str1))
	//	fmt.Println(reflect.TypeOf(str2))

	//fmt.Println(interface{}(str1).(string))
	//fmt.Println(interface{}(str2).([]int32))

	//var str1 string
	//str1 = "h我们都是好孩子"
	//fmt.Println(&str1)
	//str2 := []rune(str1)
	//str3 := []rune(str1)
	//fmt.Println(&str2)
	//str2[1] = '你'
	//str3[1] = "你"
	//
	//
	//
	//fmt.Println(str2)
	//str1 += ",我也是"
	//fmt.Println(&str1)
	//fmt.Println(str1)
	//str1 = str1[:5]
	//fmt.Println(str1)

	//var body strings.Builder
	//var body1 strings.Builder
	//cpBody1 := body1
	//body.WriteString("a")
	//cpBody := body
	//fmt.Println(cpBody)
	//fmt.Println(cpBody1)
	//fmt.Println(len(body.String()))
	//cap1 := body.Cap()
	//fmt.Println(cap1)
	//body.Grow(10)
	//cap2 := body.Cap()
	//len2 := body.Len()
	//fmt.Println(cap2)
	//fmt.Println(len2)
	//body.WriteString("bcdefghijk")
	//fmt.Println(len(body.String()))
	//cap1 = body.Cap()
	//fmt.Println(cap1)
	////body.Write([]byte(str))
	//fmt.Println(body.String())
	//body.Reset()
	//fmt.Println(body.String())

	//slia := []int{1, 2, 3, 4, 5, 6}
	//slib := make([]int, 12)
	//rs := copy(slib, slia)
	//fmt.Println(rs)
	//fmt.Println(slia)
	//fmt.Println(slib)
	//
	//var builder1 strings.Builder
	//builder1.Grow(1)
	//builder3 := builder1
	//builder3.Grow(1) // 这里会引发 panic。
	//_ = builder3

	//var builder strings.Builder
	//builder.WriteString("abcdefg hiiadfowerdsfiower")
	//var reader strings.Reader
	//fmt.Println(reader.Len(), reader.Size())

	//reader1 := strings.NewReader("NewReader returns a new Reader reading from s. " +
	//	"It is similar to bytes.NewBufferString but more efficient and read-only.")
	//fmt.Printf("The size of reader: %d\n", reader1.Size())
	//fmt.Printf("The reading index in reader: %d\n", reader1.Size()-int64(reader1.Len()))
	//
	//buf1 := make([]byte, 47)
	//n, _ := reader1.Read(buf1)
	//fmt.Println(n)
	//fmt.Printf("The size of reader: %d\n", reader1.Size())
	//fmt.Printf("The reading index in reader: %d\n", reader1.Size()-int64(reader1.Len()))
	////
	////buf11 := make([]byte, 10)
	////n11, _ := reader1.Read(buf11)
	////fmt.Println(n11)
	////fmt.Printf("The size of reader: %d\n", reader1.Size())
	////fmt.Printf("The reading index in reader: %d\n", reader1.Size()-int64(reader1.Len()))
	////fmt.Println(reader1)
	//buf2 := make([]byte, 21)
	//offset1 := int64(64)
	//n, _ = reader1.ReadAt(buf2, offset1)
	//fmt.Printf("%d bytes were read. (call ReadAt, offset: %d)\n", n, offset1)
	//fmt.Printf("The reading index in reader: %d\n",
	//	reader1.Size()-int64(reader1.Len()))
	//fmt.Println()
	//
	//offset2 := int64(18)
	//expectedIndex := reader1.Size() - int64(reader1.Len()) + offset2
	//fmt.Printf("Seek with offset %d and whence %d ...\n", offset2, io.SeekCurrent)
	//readingIndex, _ := reader1.Seek(offset2, io.SeekCurrent)
	//fmt.Printf("The reading index in reader: %d (returned by Seek)\n", readingIndex)
	//fmt.Printf("The reading index in reader: %d (computed by me)\n", expectedIndex)
	//
	//n, _ = reader1.Read(buf2)
	//fmt.Printf("%d bytes were read. (call Read)\n", n)
	//fmt.Printf("The reading index in reader: %d\n",
	//	reader1.Size()-int64(reader1.Len()))

	//str1 := "ab_cd d b dwq iow,Ate,Booqw,ioo,我123,13,21,12.,"
	//fmt.Println(strings.HasPrefix(str1, "ab"))
	//fmt.Println(strings.HasSuffix(str1, " ."))
	//fmt.Println(strings.Contains(str1, "13"))
	//fmt.Println(strings.ToUpper(str1))
	//fmt.Println(strings.ToLower(str1))
	//fmt.Println(strings.Title(str1))
	//fmt.Println(strings.Count(str1, "12"))
	//fmt.Println(strings.ContainsAny(str1, "879"))
	//fmt.Println(strings.LastIndex(str1, "2"))
	//fmt.Println(strings.Index(str1, "我"))
	////fmt.Println(strings.IndexByte(str1, '我'))
	//fmt.Println(strings.IndexRune(str1, '我'))
	//fmt.Println(strings.IndexFunc(str1, func(r rune) bool {
	//	return string(r) == "1"
	//}))
	//fmt.Println(strings.Fields(str1))
	//fmt.Println(strings.Split(str1, ","))
	//fmt.Println(strings.SplitAfter(str1, ","))
	//fmt.Println(strings.SplitAfterN(str1, ",", 2))
	//fmt.Println(strings.SplitN(str1, ",", 18))
	//
	//fmt.Println(strings.ToTitle(str1))
	//fmt.Println(strings.TrimLeft(str1, "acd"))
	//fmt.Println(strings.Repeat("a", 5))
	//slia := []string{"a", "b", "c"}
	//fmt.Println(strings.Join(slia, "="))
	//fmt.Println(strings.Replace(str1, "1", "9", 2))
	//fmt.Println(strings.Compare("a", "D"))
	//fmt.Println(strings.EqualFold("ABC", "abc"))
	//fmt.Println(strings.EqualFold("ABC", "ab"))

	//var buffer1 bytes.Buffer
	//contents := "Simple byte buffer for marshaling data."
	//fmt.Printf("Writing contents %q ...\n", contents)
	//buffer1.WriteString(contents)
	//fmt.Printf("The length of buffer: %d\n", buffer1.Len())
	//fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap())
	//
	//p1 := make([]byte, 7)
	//n, _ := buffer1.Read(p1)
	//fmt.Printf("%d bytes were read. (call Read)\n", n)
	//fmt.Printf("The length of buffer: %d\n", buffer1.Len())
	//fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap())
	//
	//buffer1.WriteString("a")
	//fmt.Printf("The length of buffer: %d\n", buffer1.Len())
	//fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap())

	//var str1 bytes.Buffer
	//str1.WriteString("我说：")
	//r, _ := os.Open("./test.txt")
	//len1, _ := str1.ReadFrom(r)
	//fmt.Println(len1, str1.Len(), str1.Cap(), str1.String())
	//str1.Write([]byte{'a','b'})
	//str1.WriteByte('a')
	//fmt.Println(str1.String())
	//str1.WriteString("a")

	//a := 'b'
	//val, ok := interface{}(a).(rune)
	//fmt.Println(val, ok)
	//switch interface{}(a).(type) {
	//case byte:
	//	fmt.Println("type byte uint8")
	//case rune:
	//	fmt.Println("type rune int32")
	//default:
	//	fmt.Println("untype")
	//}

	//src := strings.NewReader(
	//	"CopyN copies n bytes (or until an error) from src to dst. " +
	//		"It returns the number of bytes copied and " +
	//		"the earliest error encountered while copying.")
	//dst := new(strings.Builder)
	//written, err := io.CopyN(dst, src, 58)
	//if err != nil {
	//	fmt.Printf("error: %v\n", err)
	//} else {
	//	fmt.Printf("Written(%d): %q\n", written, dst.String())
	//}

	//io.Writer()
	//io.Reader()
	//io.ReadWriter()
	//io.Pipe()
	//strings.Reader{}

	//url := "https://www.baidu.com"
	//res, err := http.Get(url)
	//fmt.Println(res.Proto, err)

	//并发操作
	//count := 0
	//var wg sync.WaitGroup
	//wg.Add(10)
	////var lock sync.Mutex
	////var count atomic2.Int32
	////count.Store(0)
	//for i:= 0; i < 10; i++ {
	//	go func() {
	//		defer wg.Done()
	//
	//		for i := 0; i < 10000; i++ {
	//			//lock.Lock()
	//			//count.Add(1)
	//			count++
	//			//lock.Unlock()
	//		}
	//	}()
	//}
	//wg.Wait()
	//fmt.Println(count)
	//fmt.Println(count.Load())

	//const (
	//	firstVal = iota
	//	secondVal
	//	thirdVal
	//	fourVal
	//)
	//fmt.Println(firstVal, secondVal,thirdVal, fourVal)
	//
	//var mu sync.Mutex
	//mu.Lock()

	//var c Counter
	//c.Lock()
	//defer c.Unlock()
	//c.Count++
	//foo(c)

	//   这种声明  mapA[1] = 1会报错
	//var mapA map[int]int
	////mapA := map[int]int{}
	////mapA := make(map[int]int, 0)
	//fmt.Println(mapA[100])
	//mapA[1] = 1
	//fmt.Println(mapA)

	//if v, ok := mapA[1]; ok {
	//	fmt.Printf("exist, v is %d\n", v)
	//} else {
	//	fmt.Println("not exist")
	//}
	//m := map[int]func(op int)int{}
	//m[1] = func(op int) int {
	//	return op
	//}
	//m[2] = func(op int) int {
	//	return op * op
	//}
	//fmt.Println(m[1](2), m[2](2))
	//fmt.Println(m)
	//delete(m, 1)
	//fmt.Println(m)

	//m := "abcdef"
	//fmt.Printf("%[1]s %[1]v", m)
	//rand.Seed(time.Now().Unix())
	//fmt.Println(rand.Intn(10), rand.Intn(20))

	//tsF := timeSpent(slowFun)
	//fmt.Println(tsF(3))

	//a := []int{1, 2, 3}
	//fmt.Println(&a[0], &a[1], &a[2])
	//fmt.Println(a)
	////a = append(a[:0], a[1:]...)
	////copyLen := copy(a, a[1:])
	////a = a[:copy(a, a[1:])]
	//a = a[1:]
	//fmt.Println(&a[0], &a[1])
	////a = a[0:2]
	//fmt.Println(a)

	//fmt.Println(Sum(1, 2, 3))
	//fmt.Println(Sum(1, 2, 3, 4))

	//defer func() {
	//	fmt.Println("defer")
	//	//recover()
	//	panic("defer panic")
	//	fmt.Println("end2")
	//}()
	//fmt.Println("start")
	//panic("panic")
	//fmt.Println("end1")
	//fmt.Println(fib(5))

	//defer func() {
	//	fmt.Println("finally")
	//}()
	//fmt.Println("start")
	////os.Exit(2)
	//panic("panic")

	//ch1 := make(chan int, 0)
	//defer func() {
	//	close(ch1)
	//	time.Sleep(time.Second * 2)
	//}()
	//
	//go func() {
	//	ticker := time.NewTicker(time.Second * 1)
	//	defer ticker.Stop()
	//	for  {
	//		select {
	//		//case r, ok := <-ticker.C:
	//		//	fmt.Println("timer", r, ok)
	//		case err, ok2 := <-time.After(time.Second * 1):
	//			fmt.Println("time out", err, ok2)
	//		case b, ok1 := <-ch1:
	//			fmt.Println("ch1 close", b, ok1)
	//			return
	//			//default:
	//			//	fmt.Println("default")
	//		}
	//	}
	//}()
	//time.Sleep(time.Second * 5)
	//fmt.Println("finish")

	//var wg sync.WaitGroup
	//for i := 0; i < 5; i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		o := GetInstance()
	//		fmt.Println(unsafe.Pointer(o))
	//	}()
	//}
	//wg.Wait()
	//fmt.Println("end")
	//fmt.Printf("before: %d\n", runtime.NumGoroutine())
	//num := testFunc()
	//time.Sleep(10 * time.Millisecond)
	//fmt.Println(num)
	//fmt.Printf("after: %d\n", runtime.NumGoroutine())
	//fmt.Println("end")
	//
	//sync.Pool{}
	//
	//buffer.NewPool()

	//map1 := map[int]int{1:1, 2:2}
	//fmt.Println(len(map1))
	//
	//var map2 map[int]int
	//fmt.Println(len(map2))
	//
	//if map1 != nil {
	//	fmt.Println("a")
	//} else {
	//	fmt.Println("b")
	//}
	//
	//fun1 := func() {
	//	fmt.Println("1")
	//}
	//if fun1 == nil {
	//	fmt.Println("ccccc")
	//}
	//
	//fun2 := func() {}
	//
	//if fun2 == nil {
	//	fmt.Println("dddd")
	//}
	//
	//
	//var map1 map[int]int
	//if map1 == nil {
	//	fmt.Println("map1 nil")
	//} else {
	//	fmt.Println("map1 not nil")
	//}
	//
	//map2 := map[int]int{}
	//if map2 == nil {
	//	fmt.Println("map2 nil")
	//} else {
	//	fmt.Println("map2 not nil")
	//}
	//
	//map3 := make(map[int]int)
	//if map3 == nil {
	//	fmt.Println("map3 nil")
	//} else {
	//	fmt.Println("map3 not nil")
	//}
	//
	//var t *testing.T
	//assert.Equal(t, map4, map3)

	//http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	//	fmt.Println("hello")
	//	fmt.Println(writer)
	//	fmt.Println(request)
	//	writer.Write([]byte("hello a"))
	//})

	//http.HandleFunc("/time", func(writer http.ResponseWriter, request *http.Request) {
	//	t := time.Now()
	//	timeStr := fmt.Sprintf("{\"time\": \"%s\"}", t)
	//	num, err := writer.Write([]byte(timeStr))
	//	fmt.Println(num, err)
	//})
	//
	//http.ListenAndServe(":8080", nil)

	//float1 := 4.4999999999999999
	//fmt.Println(math.Floor(float1 + 0.5))
	//fmt.Println(math.Ceil(float1))
	//float2 := 4.499999999999999
	//fmt.Println(math.Floor(float2 + 0.5))
	//fmt.Println(math.Ceil(float2))
	//
	//fmt.Println(math.Round(float2))

	//var float33 float32 = 4.999999
	//float3 := float64(float33)
	//fmt.Println(math.Floor(float3))

	//sli1 := []int{1, 2, 3, 4, 5, 6}
	//rand.Seed(time.Now().Unix())
	//rand.Shuffle(len(sli1), func(i, j int) {
	//	sli1[i], sli1[j] = sli1[j], sli1[i]
	//})
	//fmt.Println(sli1)

	//ctx := context.Background()
	//ctx, cancel := context.WithCancel(context.Background())
	//var wg sync.WaitGroup
	//for i := 0; i < 10; i++{
	//	wg.Add(1)
	//	go func(ctx context.Context, cancel func(), i int) {
	//		defer wg.Done()
	//		for  {
	//			select {
	//			case <-ctx.Done():
	//				fmt.Println("done")
	//				return
	//			default:
	//				if i == 6 {
	//					fmt.Println("cancel")
	//					cancel()
	//				}
	//			}
	//		}
	//	}(ctx, cancel, i)
	//}
	//wg.Wait()
	//fmt.Println(runtime.NumGoroutine())

	//s1 := "big"
	//s11 := []byte(s1)
	//s11[0] = 'a'
	//fmt.Println(s1, string(s11))

	//for i := 0; i < 100; i++{
	//	if i == 60 {
	//		goto res
	//	}
	//	fmt.Println(i)
	//}
	//res:
	//	fmt.Println("end")
	//tag2:
	//for i := 0 ; i < 10; i++{
	//	tag1:
	//	for j := 0; j < 5; j++ {
	//		if j == 2 {
	//			break tag1
	//		}
	//		fmt.Println(j)
	//	}
	//	if i == 9 {
	//		break tag2
	//	}
	//	fmt.Println(i)
	//}
	//fmt.Println("end")

	//a := 1
	//b := float32(a)
	//
	//var c map[int]int
	//fmt.Println(len(c))
	//fmt.F

	//a := 1
	//pi := &a
	//nameP := new(string)
	//fmt.Println(pi, *pi, nameP)

	//fmt.Println(f1())
	//fmt.Println(f2())
	//fmt.Println(f3())
	//var a int
	//fmt.Println(a)

	//c := 3 + 4i
	//fmt.Println(cmplx.Abs(c))
	//fmt.Println(cmplx.Pow(math.E, 1i * math.Pi) + 1)
	//fmt.Println(cmplx.Exp(1i * math.Pi) + 1)
	//fmt.Printf("%.3f",cmplx.Exp(1i * math.Pi) + 1)

	//fmt.Println(cpp, python,golang, javascript)
	//fmt.Println(1 << 10)
	//fmt.Println(math.Pow(2, 10))
	//content, _ := ioutil.ReadFile("test.txt")
	//contentStr := string(content)
	//fmt.Println(content,contentStr)
	//fmt.Printf("%s\n", content)

	//printFile("test.txt")

	//
	//fmt.Println(converToBin(5),
	//converToBin(13),converToBin(89234234234))

	//fmt.Println(apply(pow, 3, 4))
	//fmt.Println(apply(func(a int, b int) int {
	//	return int(math.Pow(float64(a), float64(b)))
	//}, 3, 4))

	//fmt.Println(sum(1,2,3,4,5))
	//var p1 person
	//p1.age = 20
	//p1.name = "xiaoming"
	//p1.height = 168
	//updatePerson(&p1)
	//fmt.Println(p1)

	//var grid [4][5]int
	//fmt.Println(grid)
	//slice := []int{1, 2, 3}
	//printSlice(slice)
	//fmt.Println(slice)

	//arr := [...]int{0,1,2,3,4,5,6}
	//sli1 := arr[2:5]
	////slice本身没有数据，是对底层array的一个view
	//fmt.Println(reflect.TypeOf(sli1))
	//fmt.Println(reflect.TypeOf(arr))
	//fmt.Println(reflect.ValueOf(arr))
	//fmt.Println(reflect.ValueOf(sli1))

	//arr := [...]int{0,1,2,3,4,5,6,7}
	//arr2 := arr[2:6]
	//s2 := arr2[3:5]
	//fmt.Println(s2) // [5,6]
	////slice可以向后扩展，不可以向前扩展
	////s[i]不可以超越len(s),向后扩展不可以超越底层数组cap(s)
	//
	//s3 := append(s2, 10)
	//s4 := append(s3, 11)
	//s5 := append(s4, 12)
	//fmt.Println(s3,s4,s5,arr)
	//s2[0] = 77
	//fmt.Println(s3,s4,s5,arr)
	//添加元素时如果超越cap，系统会重新分配更大的底层数组
	//由于值传递的关系，必须接收append的返回值(因为slice的扩容可能导致新的底层数组产生，其实是可能导致slice内部的ptr和cap的变动)

	//s1 := []int{1,2,3,4}
	//s2 := make([]int, 10, 22)
	//fmt.Println(s1, s2)
	//copy(s2, s1)
	//fmt.Println(s1, s2)
	////删除元素
	//s2 = append(s2[:3], s2[4:]...)
	//fmt.Println(s2)

	//map for range取出是无序的，hash无序，想要有序可以取出所有的key，排序后取出
	//map map[key]，key不存在，返回零值，所以需要使用 value, ok := map[key]，判断ok的值
	//delete(map, key)来删除
	//len(m)获取元素个数

	//map的key
	//map使用哈希表，所以map的key必须可以比较相等
	//除了slice，map,function的内件类型都可以作为key
	//struct类型不包含上述类型字段，也可以作为key

	//字符和字符串处理
	//使用range遍历pos,rune对
	//使用utf8.RuneCountInString获得字符数量
	//使用len获得字节长度
	//使用[]byte获得字节
	//var f1 float32
	//fmt.Println(f1)
	//var s1 []int
	//fmt.Println(s1)
	//fmt.Println(s1 == nil)
	//
	//i1 := 11
	//i2 := strconv.FormatInt(int64(i1), 2)
	//i4 := strconv.FormatInt(int64(i1), 4)
	//i6 := strconv.FormatInt(int64(i1), 6)
	//i8 := strconv.FormatInt(int64(i1), 8)
	//fmt.Println(i2, i4, i6, i8)

	//var f1 = 3.12345670891234567890123
	//fmt.Println(strconv.FormatFloat(f1, 'f', -1, 32))
	//fmt.Println(strconv.FormatFloat(f1, 'f', -1, 64))
	//fmt.Println(strconv.FormatFloat(f1, 'f', 4, 64))
	//fmt.Println(strconv.FormatFloat(f1, 'e', 4, 64))
	//fmt.Println(strconv.FormatFloat(f1, 'e', 8, 32))
	//f32 := float32(f1)
	//f64 := float64(f1)
	//fmt.Println(f32, f64)
	//
	//fmt.Println(strconv.ParseInt("190", 10, 32))

	//num := 15
	//fmt.Printf("十进制 = %d\n", num)
	//fmt.Printf("八进制 = %o\n", num)
	//fmt.Printf("十六进制 = %x\n", num)
	//fmt.Printf("二进制 = %b\n", num)
	//fmt.Printf("类型 = %T\n", num)
	//
	//num1 := 16
	//fmt.Println("num1 =", 1, "num2 =", 2)
	//fmt.Println(1, 2)
	//fmt.Println(num, num1)

	//n, err := fmt.Fprintf(os.Stdout, "age is %v, name is %v\n", 22, "小红")
	//fmt.Println(n, err)
	//s1 := "age is 22, name is 小红\n"
	//fmt.Println(len(s1))

	//var num1 int
	//var num2 int
	//fmt.Scanf("%d%d", &num1, &num2)
	//fmt.Println(num1, num2)

	/*	var num1 int
		var num2 int
		fmt.Scan(&num1, &num2)
		fmt.Println(num1, num2)*/
	//
	//var num3 float32
	//var num4 float32
	//fmt.Scan(&num3, &num4)
	//fmt.Println(num3, num4)
	//fmt.Println(runtime.NumCPU())

	//a := addr()
	//for i := 0; i < 10; i++ {
	//	fmt.Println(a(i))
	//}
	//func1 := fibonacci()
	//for i := 0; i < 10; i++ {
	//	fmt.Println(func1())
	//}

	//file, err := os.Create("test.txt")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer file.Close()
	//
	//writer := bufio.NewWriter(file)
	//defer writer.Flush()
	//f := fibonacci()
	//for i := 0; i < 20; i++ {
	//	fmt.Fprintln(writer, f())
	//}

	//for i := 0 ; i < 100; i++ {
	//	defer fmt.Println(i)
	//}

	//file, err := os.OpenFile("a.txt", os.O_EXCL|os.O_TRUNC, 0666)
	//defer file.Close()
	//
	//err = errors.New("this is custom err")
	//if err != nil {
	//	if pathError, ok := err.(*os.PathError); !ok {
	//		panic(err)
	//	} else {
	//		fmt.Println(pathError.Path, pathError.Op, pathError.Err)
	//	}
	//}

	//http.HandleFunc("/list/", filelisting.HandleFileListen)
	//
	//err := http.ListenAndServe(":8888", nil)
	//if err != nil {
	//	fmt.Println(err)
	//}

	//fmt.Println(filelisting.LenghtOfNoRepeatingSubStr("abcdabdcewadfewoidsfqwow"))

	//for i := 0; i < 1000; i++ {
	//	go func() {
	//		for  {
	//			fmt.Println("hello kitty")
	//		}
	//	}()
	//}
	//time.Sleep(time.Minute)
	//m := make(map[int]int, 10)
	//for i := 1; i<= 10; i++ {
	//	m[i] = i
	//}
	//
	//for k, v := range(m) {
	//	go func() {
	//		fmt.Println("k ->", k, "v ->", v)
	//	}()
	//}
	//time.Sleep(time.Second * 6)

	//a := 11 //1011
	//fmt.Println(a<<2) // 101100

}

type intGen func() int
type person1 struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p person1) String() string {
	panic("implement me")
}

func fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func addr() func(int) int {
	num := 0
	return func(i int) int {
		num += i
		return num
	}
}

func printSlice(sli1 []int) {
	sli1[0] = 11
}

func updatePerson(person2 *person) {
	person2.age = 25
	person2.name = "xiaobai"
	person2.height = 180
}

type person struct {
	age    int
	name   string
	height int
}

func sum(numbers ...int) int {
	sum1 := 0
	for i := range numbers {
		sum1 += numbers[i]
	}
	return sum1
}

func pow(a, b int) int {
	return int(math.Pow(float64(a), float64(b)))
}

func apply(op func(int, int) int, a, b int) int {
	p := reflect.ValueOf(op).Pointer()
	opName := runtime.FuncForPC(p).Name()
	fmt.Printf("Calling function %s with args "+"(%d,%d) ", opName, a, b)
	return op(a, b)
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println("-------")
		fmt.Println(scanner.Text())
	}
}

func converToBin(n int) string {
	result := ""
	for ; n > 0; n /= 2 {
		lsb := n % 2
		result = strconv.Itoa(lsb) + result
	}
	return result
}

const (
	cpp = iota + 1
	_
	python
	golang
	javascript
)

func f1() (r int) {
	defer func() {
		r++
		fmt.Printf("r is %d\n", r)
	}()

	return 0
}

func f2() (r int) {
	t := 5
	defer func() {
		t = t + 5
		fmt.Printf("t is %d\n", t)
	}()

	return t
}

func f3() (r int) {
	defer func() {
		r = r + 5
		fmt.Printf("r is %d\n", r)
	}()
	return 1
}

func testFunc() string {
	ch1 := make(chan string, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			ret := doTask(i)
			ch1 <- ret
		}(i)
	}
	return <-ch1
}

func doTask(i int) string {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("task id is %d\n", i)
}

var once sync.Once
var obj *testOnce

func GetInstance() *testOnce {
	once.Do(func() {
		fmt.Println("create obj")
		obj = new(testOnce)
	})
	return obj
}

type testOnce struct {
}

func fib(n int) []int {
	slia := []int{1, 1}
	for i := 2; i <= n; i++ {
		slia = append(slia, slia[i-2]+slia[i-1])
	}
	return slia
}

func Sum(args ...int) int {
	i := 0
	for _, v := range args {
		i += v
	}
	return i
}

func timeSpent(inner func(op int) int) func(op int) int {
	return func(n int) int {
		start := time.Now()
		ret := inner(n)
		fmt.Println("time Spent:", time.Since(start).Seconds())
		return ret
	}
}

func slowFun(op int) int {
	opDura := time.Duration(op)
	time.Sleep(time.Second * opDura)
	return op
}

func foo(c Counter) {
	c.Lock()
	defer c.Unlock()
	fmt.Println("in foo")
}

type Counter struct {
	sync.Mutex
	Count uint64
}

func GetString() {
	name := "aaa"
	fmt.Println(&name)

	{
		name := 1
		fmt.Println(&name)
	}
}

func update(a *int) {
	a = interface{}(a).(*int)
}

func (a animal) String1() {
	a.atype = "c"
}

func (a *animal) String2() {
	a.atype = "d"
}

type Cat struct {
	name           string // 名字。
	scientificName string // 学名。
	category       string // 动物学基本分类。
}

func New(name, scientificName, category string) Cat {
	return Cat{
		name:           name,
		scientificName: scientificName,
		category:       category,
	}
}

func (cat *Cat) SetName(name string) {
	cat.name = name
}

func (cat Cat) SetNameOfCopy(name string) {
	cat.name = name
}

func (cat Cat) Name() string {
	return cat.name
}

func (cat Cat) ScientificName() string {
	return cat.scientificName
}

func (cat Cat) Category() string {
	return cat.category
}

func (cat Cat) String() string {
	return fmt.Sprintf("%s (category: %s, name: %q)",
		cat.scientificName, cat.category, cat.name)
}

//func (a animal) String() string {
//	b := a.atype + "edf"
//	return fmt.Sprintln(b)
//}

type animal struct {
	atype    string
	category string
}

func modifyArray(a [3][]string) [3][]string {
	a[0][0] = "x"
	return a
}

type operate func(a, b int) int

func calculate(a, b int, c operate) int {
	return c(a, b)
}

type Printer func(contents string) (n int, err error)

func printStd(content string) (num int, err1 error) {
	return fmt.Println(content)
}

//func fibonacci(c, quit chan int) {
//	x, y := 0, 1
//	for {
//		select {
//		case c <- x:
//			x, y = y, x+y
//		case <-quit:
//			fmt.Println("quit")
//			return
//		}
//	}
//}

func Chann(ch chan int, stopCh chan bool) {
	for j := 0; j < 10; j++ {
		ch <- j
		time.Sleep(time.Second)
	}
	stopCh <- true
}

func ap(array [3]int) {
	fmt.Printf("ap brfore:  len: %d cap:%d data:%+v\n", len(array), cap(array), array)
	array[0] = 1
	//array = append(array, 10)
	fmt.Printf("ap after:   len: %d cap:%d data:%+v\n", len(array), cap(array), array)
}

type testSendParam struct {
	a int
	b int
}

func updateData(data [2]int) [2]int {
	//data.a = 11
	//data = 44
	data[0] = 3
	return data
}

func Printf() {
	return
}

//func init()  {
//	fmt.Println("main")
//}
//
//func init()  {
//	fmt.Println("main1")
//}

func printArr(arr *[5]int) {
	arr[0] = 10
	for i, v := range arr {
		fmt.Println(i, v)
	}
}

func TempSock(totalDuration float64) string {
	// serve
	rand.Seed(time.Now().Unix())
	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock", rand.Int()))
	l, err := net.Listen("unix", sockFileName)
	if err != nil {
		panic(err)
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		buf := make([]byte, 16)
		data := ""
		progress := ""
		for {
			_, err := fd.Read(buf)
			if err != nil {
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)
			cp := ""
			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				cp = fmt.Sprintf("%.2f", float64(c)/totalDuration/1000000)
			}
			if strings.Contains(data, "progress=end") {
				cp = "done"
			}
			if cp == "" {
				cp = ".0"
			}
			if cp != progress {
				progress = cp
				fmt.Println("progress: ", progress)
			}
		}
	}()

	return sockFileName
}

//func ExampleShowProgress(inFileName, outFileName string) {
//	a, err := ffmpeg.Probe(inFileName)
//	if err != nil {
//		panic(err)
//	}
//	totalDuration, err := probeDuration(a)
//	if err != nil {
//		panic(err)
//	}
//
//	err = ffmpeg.Input(inFileName).
//		Output(outFileName, ffmpeg.KwArgs{"c:v": "libx264", "preset": "veryslow"}).
//		GlobalArgs("-progress", "unix://"+TempSock(totalDuration)).
//		OverWriteOutput().
//		Run()
//	if err != nil {
//		panic(err)
//	}
//}

type probeFormat struct {
	Duration string `json:"duration"`
}

type probeData struct {
	Format probeFormat `json:"format"`
}

func probeDuration(a string) (float64, error) {
	pd := probeData{}
	err := json.Unmarshal([]byte(a), &pd)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(pd.Format.Duration, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func probeOutputDuration(a string) (float64, error) {
	pd := probeData{}
	err := json.Unmarshal([]byte(a), &pd)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(pd.Format.Duration, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func loadImage(fileInput string) (image.Image, error) {
	f, err := os.Open(fileInput)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func GetImageByUrl(url string) (image.Image, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, nil
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(res.Body)
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil
	}
	reader := bytes.NewReader(data)

	urlSli := strings.Split(url, ".")
	var img image.Image
	switch urlSli[len(urlSli)-1] {
	case "png":
		img, err = png.Decode(reader)
	case "jpeg":
		img, err = jpeg.Decode(reader)
	case "gif":
		img, err = gif.Decode(reader)
	}
	//img, _, err := image.Decode(reader)
	return img, err
}

//func FindDomiantColor(url string) (c string, err error) {
//	img, err := GetImageByUrl(url)
//	if err != nil {
//		return
//	}
//	return dominantcolor.Hex(dominantcolor.Find(img)), nil
//}

//func FindDomiantColor(fileInput string) (string, error) {
//	f, err := os.Open(fileInput)
//	defer f.Close()
//	if err != nil {
//		fmt.Println("File not found:", fileInput)
//		return "", err
//	}
//	img, _, err := image.Decode(f)
//	if err != nil {
//		return "", err
//	}
//
//	return dominantcolor.Hex(dominantcolor.Find(img)), nil
//}

//func getDominantHue(url string) (string, error) {
//	res, err := http.Get(url)
//	if err != nil {
//		return "", err
//	}
//
//	defer res.Body.Close()
//	data, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return "", err
//	}
//	reader := bytes.NewReader(data)
//	img, _, err := image.Decode(reader)
//	if err != nil {
//		return "", err
//	}
//	return dominantcolor.Hex(dominantcolor.Find(img)), nil
//}

func testDefer() int {
	i := 1

	//for ; i < 10; i++ {
	//	if i > 5 {
	//		i = 9
	//		return i
	//	}
	//}

	defer func(i *int) {
		fmt.Println(*i)
	}(&i)

	i++

	time.Sleep(time.Second * 2)
	return i
}

type AlbumInfo struct {
	Data struct {
		AlbumId       int    `json:"album_id"`
		Title         string `json:"title"`
		Author        string `json:"author"`
		Anchor        string `json:"anchor"`
		Category1Name string `json:"category1_name"`
		Category2Name string `json:"category2_name"`
		Intro         string `json:"intro"`
		OverStatus    int    `json:"over_status"`
		ImageLink     string `json:"image_link"`
		RecordCompany string `json:"record_company"`
		BookTitle     string `json:"book_title"`
		Tag           string `json:"tag"`
		AlbumType     int    `json:"album_type"`
	} `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func trimStruct(stru1 interface{}) {
	stru2 := stru1
	//v := reflect.ValueOf(stru1)
	v1 := reflect.ValueOf(stru2)
	//fmt.Println(v)
	//fmt.Println(v1.NumField())
	for i := 0; i < v1.NumField(); i++ {
		//fmt.Println(v.Field(i))
		//fmt.Println(v.Type().Field(i))
		//if v1.Type().Field(i).Type.Kind() == reflect.String {
		name := v1.Type().Field(i).Name
		f := v1.Elem().FieldByName(name)
		if f.Kind() == reflect.String {
			fmt.Println(f)
			//f.SetString(strings.TrimSpace(v1.Field(i).String()))
		}
		//fmt.Println(v1.Type().Field(i).Name)
		//fmt.Println(v.Field(i).Elem().String())
		//v.Field(i).Elem().SetString(strings.TrimSpace(v.Field(i).Elem().String()))
		//fmt.Println(v.Field(i))
		//}
	}

	//t := reflect.TypeOf(stru1)
	//for i := 0; i < t.NumField(); i++ {
	//	//fmt.Println(t.Field(i))
	//	//fmt.Println(t.Field(i).Type)
	//	//fmt.Println(t.Field(i).Name)
	//	if t.Field(i).Type.Kind() == 24 {
	//		v.Field(i).Elem().SetString(strings.TrimSpace(v.Field(i).String()))
	//	}
	//}
}

type User struct {
	Id    int
	Name  string
	Age   int
	House string
	Color string
}

func SetValue(o interface{}) {
	v := reflect.ValueOf(o)
	v = v.Elem()

	for i := 0; i < v.NumField(); i++ {
		//fmt.Println(v.Field(i).Type().Name())
		//fmt.Println(v.Type().Field(i).Name)
		name := v.Type().Field(i).Name
		f := v.FieldByName(name)
		if f.Kind() == reflect.String {
			//f.SetString(f.String())
			f.SetString(strings.TrimSpace(f.String()))
		}
	}

	//f := v.FieldByName("Name")
	//if f.Kind() == reflect.String {
	//	f.SetString("XIAOMING")
	//}
	//{1  xiaohong 20  shanghai  black}
	//{1 xiaohong 20 shanghai black}

}

// 绑方法
//func (u User) Hello() {
//	fmt.Println("Hello")
//}

// 传入interface{}
func Poni(o interface{}) {
	t := reflect.TypeOf(o)
	fmt.Println("类型：", t)
	fmt.Println("字符串类型：", t.Name())
	// 获取值
	v := reflect.ValueOf(o)
	fmt.Println(v)
	// 可以获取所有属性
	// 获取结构体字段个数：t.NumField()
	for i := 0; i < t.NumField(); i++ {
		// 取每个字段
		f := t.Field(i)
		fmt.Printf("%s : %v\n", f.Name, f.Type)
		// 获取字段的值信息
		// Interface()：获取字段对应的值
		val := v.Field(i).Interface()
		fmt.Println("val :", val)
	}
	fmt.Println("=================方法====================")
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fmt.Println(m.Name)
		fmt.Println(m.Type)
	}

}

// 匿名字段
type Boy struct {
	User
	Addr string
}

func (u User) Hello(name string) {
	fmt.Println("Hello：", name)
}

type Student struct {
	Name string `json:"name1" db:"name2"`
	Age  int    `json:"age" xorm:"age_int,omitempty"`
}

func initValue() *int {
	i := 4
	return &i
}
