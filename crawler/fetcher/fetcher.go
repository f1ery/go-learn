package fetcher

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
)

func init()  {
	fmt.Println("fetcher")
}

func Fetch(url string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	//增加header选项
	request.Header.Add("Cookie", "ec=od8BJKEz-1658965656275-51eab52365d4d-957290223; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1658965671; FSSBBIl1UgzbN7NO=5193dfuhqUx1sKbBAziYsgshyBomnCbI6_3RoXhREWsjRQ37yf8QpPsKARFrocb6Ojm1KhhgcizOWftf2JkyrFa; sid=xOF57NV82A8QTeksGF3K; _exid=b8U%2FR%2FskeV6fqCjvzEQAlqWcIbKTzwK8WZwRi5j2trRnp9I1%2FllyQGWQ31Lar66wqDdjZYfbJT8xpoaJer%2FCYQ%3D%3D; _efmdata=U7RSWpAEfJfRQaSSrxpggXqaDuN3iAMwE04MSTzH89lxmKHd%2BDn2AKyfGqazLq87wCq9t7dm9Ml1aYgc7D1l28fdELcQ1mBtVO3QB1msfiU%3D; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1659661181; FSSBBIl1UgzbN7NP=53XS.GCh7F67qqqDcF4RYFAkPddwdU9lSDcuxMSnWomI5TOWnVun7GPjYDOpZVQ9bRyHvmbEmMJ22coCT805NvqVvQR5bu2gPn3usyvQkTBZMCJ8tRX_BtLlMklwrjqw1S.3L1JbysMMTfQ6Smrn3J7qT5iJc6TOWaqV7wpjYxlZ4qxb.O_403Fc0v5e9roCB5Qt_BYrdgSYanNnSPrWb.KyV01K0GxGpjrV0DqdQv7N6BfkNedIj7BStSAssBODusgP_b39Ru44HbVB8LHIP9J")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36")
	request.Header.Add("X-Requested-With", "xxxx")

	resp, err := client.Do(request)
	//resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//if resp.StatusCode != http.StatusOK {
	//	//fmt.Println("Error: status code", resp.StatusCode)
	//	return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	//}

	//utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	//all, err := ioutil.ReadAll(utf8Reader)

	//all, err := ioutil.ReadAll(resp.Body)

	return ioutil.ReadAll(utf8Reader)
}

func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetch error:%v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "text/html")
	//fmt.Println(e, name, certain)
	return e
}