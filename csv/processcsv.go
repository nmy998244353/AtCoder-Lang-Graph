package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
)

type keys struct {
	user string
	lang string
	problem string
}

func main() {
	input, _ := os.Open("submissions.csv")
	rd := csv.NewReader(input)
	mp := make(map[keys]int)
	tmp := make([][]string, 0)
	t := 0
	for {
		t++
		if t%1000000==0 {
			fmt.Println(t)
		}
		line, err := rd.Read()
		if err == io.EOF {
			break
		}
		problem := line[2]
		user := line[4]
		if result := line[8]; result != "AC" {
			continue
		}
		lang_ful := line[5]
		epoc_time := line[1]
		lang := ""
		for i := 0; i < len(lang_ful); i++ {
			if lang_ful[i] == ' ' && lang_ful[i+1] == '('{
				break
			}
			lang += string(lang_ful[i])
		}
		if len(lang) > 3 && (lang[:3] == "C++" || lang[:3] == "IOI") {
			lang = "C++"
		} else if len(lang) > 6 && lang[:6] == "Python" {
			lang = "Python"
		}
		if _, exist := mp[keys{user, lang, problem}]; exist{
			continue
		}
		tmp = append(tmp, []string{user,lang,epoc_time})
		mp[keys{user, lang, problem}] = 1
	}
	input.Close()
	fmt.Println("")
	sort.Slice(tmp, func(i,j int) bool {return tmp[i][0] < tmp[j][0]})
	os.RemoveAll("users")
	os.Mkdir("users", 0777)
	pf := "@@"
	var output *os.File
	var wr *csv.Writer = nil
	for _, vec := range tmp{
		if len(vec[0])<2 {
			continue
		}
		if vec[0][0:2] != pf {
			if(pf != "@@"){
				wr.Flush()
			}
			output.Close()
			pf = vec[0][0:2]
			fmt.Println(pf)
			output, _ = os.Create("users/"+pf+".csv")
			wr = csv.NewWriter(output)
		}

		wr.Write([]string{vec[0], vec[1], vec[2]})
	}
	wr.Flush()
	output.Close()
	return
}