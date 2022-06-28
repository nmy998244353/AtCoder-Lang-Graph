package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	//"os/exec"
)



func countlangac(user []byte) (map[string][]int, bool){
	ret := make(map[string][]int)
	file, err := os.Open("csv/users/"+string(user[0:2])+".csv")
	defer file.Close()
	if err != nil{
		fmt.Println(string(user), err)
		return ret, true
	}
	rd := csv.NewReader(file)
	t := 0
	erret := true
	for {
		t++
		line, err_read := rd.Read()
		if err_read == io.EOF {
			break
		}
		user_id := line[0]
		if user_id != string(user){
			continue
		}
		erret = false
		lang := line[1]
		epoc_time, err := strconv.Atoi(line[2])
		if err != nil {
			continue
		}
		year := epoc_time/(60*60*24*365)+1970
		if year < 2010{
			continue
		}
		ret[lang] = append(ret[lang], epoc_time)
	}
	return ret,erret
}

func putjson(user []byte) []byte {
	data , err := countlangac(user)
	if(err){
		return []byte("error")
	}
	e := make(map[string][]int)
	for lang, _ := range data {
		e[lang] = make([]int, 365*100)
		for _, etime := range data[lang] {
			e[lang][etime/(60*60*24)]++
		}
	}
	ret, _ := json.Marshal(&e)
	return ret
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	r := newRoom()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.Handle("/", &templateHandler{filename: "index.html"})
	http.Handle("/room", r)
	go r.run()

	http.ListenAndServe(":"+port, nil)
	//output, _ := exec.Command("sh", "./update.sh").Output()
	//fmt.Println(output)

}
