package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// create a data type
type data struct {
	ip    []string
	alias []string
}

// inititalize map of struct
var u = make(map[string]data)

func main() {
	r := read(os.Args[1])
	consolidate(r)
	writeCSV()

}

func read(f string) [][]string {

	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(file)

	r, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return r
}

func consolidate(r [][]string) {

	// check no. of columns by Headers
	col := len(r[0])
	fmt.Printf("Found %d Columns: %s\n", col, r[0])

	// u[r[1][0]] = r[1][]
	// bar := prog
	// Each Row
	for i := 1; i < len(r); i++ {
		// Col of each row
		// for j := 0; j < len(r[0]); j++ {

		_, exists := u[r[i][0]]
		if !exists {
			// Add new entry
			fmt.Println(r[i][0], "Doesn't Exist")
			u[r[i][0]] = data{
				ip:    regexIP(r[i][1]),
				alias: add(r[i][2]),
			}
		} else {
			// append all data
			fmt.Println(r[i][0], "Already Exist")
			u[r[i][0]] = data{
				ip:    append(u[r[i][0]].ip, regexIP(r[i][1])...),
				alias: append(u[r[i][0]].alias, add(r[i][2])...),
			}
		}

		// printing
		fmt.Println(u[r[i][0]])
		// time.Sleep(2 * time.Second)

		// }
	}

}

// 8.8.8.8, 4.4.4.4
/* 
	8.8.8.8,	4.4.4.4
*/

func strip(d string) []string {
	// replace newline char with comma
	d = strings.ReplaceAll(d, "\n", ",")
	// trim leading and trailing spaces
	d = strings.TrimSpace(d)
	// replace all blank spaces with comma
	d = strings.ReplaceAll(d, " ", ",")
	// return []string
	return strings.Split(d, ",")
}

func add(d string) []string {

	// append(u[r[i][0]].ip,newIP
	res := strip(d)
	return res
}

func regexIP(s string) []string {
	IPv6 := `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	IPv4 := `(\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`

	r4 := regexp.MustCompile(IPv4)
	r6 := regexp.MustCompile(IPv6)

	res := r4.FindAllString(s, -1)
	res = append(res, r6.FindAllString(s, -1)...)
	res = removeDuplicateValues(res)
	return res
}

func removeDuplicateValues(Slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range Slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func writeCSV() error {
	f, err := os.OpenFile("consolidated.csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)

	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.UseCRLF = true
	defer w.Flush()
	w.Write([]string{"Domain", "IP", "Alias"})
	for domain, data := range u {
		if err := w.Write([]string{domain, strings.Join(data.ip, "\r\n"), strings.Join(data.alias, "\r\n")}); err != nil {
			log.Fatalln("error writing record to file", err)
		}

		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
