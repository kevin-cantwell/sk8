package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
)

var (
	fileTypeFlag = flag.String("type", "csv", "File type (csv, json).")
	fieldsFlag   = flag.String("fields", "", "Comma-separated fields to obfuscate.")
)

func main() {
	flag.Parse()

	switch *fileTypeFlag {
	case "csv":
		obfuscateCSV(os.Stdin, os.Stdout)
	default:
		log.Fatalf("unsupported file type: %s", *fileTypeFlag)
	}
}

func obfuscateCSV(in io.Reader, out io.Writer) {
	fields := strings.Split(*fieldsFlag, ",")
	fieldsToObfuscate := make(map[string]bool, len(fields))
	for _, field := range fields {
		fieldsToObfuscate[field] = true
	}

	header := []string{}

	r := csv.NewReader(in)
	w := csv.NewWriter(out)

	for {
		record, err := r.Read()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("error: %+v", err)
		}
		if len(header) == 0 {
			header = record
			w.Write(header)
			w.Flush()
			continue
		}

		for i, val := range record {
			if !fieldsToObfuscate[header[i]] {
				continue
			}
			record[i] = hashValue(val)
		}

		w.Write(record)
		w.Flush()
	}
}

var (
	runeHash = map[rune]rune{}
)

func init() {
	for r := 'a'; r <= 'z'; r++ {
		runeHash[r] = rune('a' + rand.Intn('z'-'a'))
	}
	for r := 'A'; r <= 'Z'; r++ {
		runeHash[r] = rune('A' + rand.Intn('Z'-'A'))
	}
	for r := '0'; r <= '9'; r++ {
		runeHash[r] = rune('0' + rand.Intn('9'-'0'))
	}
}

func hashValue(val string) string {
	out := make([]rune, len([]rune(val)))

	for i, r := range val {
		h := runeHash[r]
		if h == 0 {
			out[i] = r
		} else {
			out[i] = h
		}
	}

	return string(out)

	// valRR := []rune(val)

	// h := md5.New()

	// if _, err := io.WriteString(h, val); err != nil {
	// 	log.Fatalf("error: %+v", err)
	// }

	// hash := fmt.Sprintf("%x", h.Sum(nil))
	// hashRR := []rune(hash)

	// // append the hash unto itself until it is at least as long as the original value
	// for diff := len(valRR) - len(hashRR); diff > 0; diff = len(valRR) - len(hashRR) {
	// 	hashRR = append(hashRR, hashRR...)
	// }

	// // keep original non-word characters, such as dashes or punctuation
	// for i, r := range valRR {
	// 	switch {
	// 	case r >= 'a' && r <= 'z':
	// 	case r >= 'A' && r <= 'Z':
	// 	case r >= '0' && r <= '9':
	// 	default:
	// 		hashRR[i] = r
	// 	}
	// }

	// return string(hashRR[0:len(valRR)])
}
