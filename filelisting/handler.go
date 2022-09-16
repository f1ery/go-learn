package filelisting

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
)

// HandleFileListen description
func HandleFileListen (writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path[len("/list/"):]
	file, err := os.Open(path)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	all, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	writer.Write(all)
}

func triangle() {
	var a, b int = 3, 4
	fmt.Println(calTriangle(a, b))
}

func calTriangle(a, b int) int {
	var c int
	c = int(math.Sqrt(float64(a * a + b * b)))
	return c
}
