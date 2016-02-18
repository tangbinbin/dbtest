package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	. "sync/atomic"
	"time"
)

var (
	input         = flag.String("I", "/var/log/mysql/mysql.log", "mysql general log file")
	output        = flag.String("O", "/tmp/out.log", "output file")
	num    uint64 = 0
)

func init() {
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go state()
	file, err := os.Open(*input)
	if err != nil {
		log.Fatalf("open file %s error, reason:%v", *input, err)
		return
	}
	defer file.Close()

	wfile, err := os.Create(*output)
	if err != nil {
		log.Fatalf("failed to create output file:%s reason:%v", *output, err)
		return
	}
	defer wfile.Close()

	br := bufio.NewReader(file)
	buf := ""
	for {
		line, err := br.ReadString('\n')
		if err != nil && err == io.EOF {
			log.Println("read file over")
			break
		}
		ws, tmp := parseLine(line, buf)
		if ws == "" {
			buf = tmp
			continue
		}
		if strings.Contains(ws, "select") ||
			strings.Contains(ws, "SELECT") {
			wfile.WriteString(ws)
			AddUint64(&num, 1)
		}
		buf = tmp
	}
}

func parseLine(line string, buf string) (ret, tmp string) {
	if strings.Contains(line, "\t") {
		if strings.Contains(line, "Quit") ||
			strings.Contains(line, "Connect") ||
			strings.Contains(line, "Prepare") ||
			strings.Contains(line, "Close stmt") ||
			strings.Contains(line, "Init DB") ||
			strings.Contains(line, "Id\tCommand") {
			return buf, ""
		}
		if strings.Contains(line, "Query") ||
			strings.Contains(line, "Execute") {
			l := strings.Split(line, "\t")
			switch len(l) {
			case 4:
				if strings.Contains(l[2], "Query") ||
					strings.Contains(l[2], "Execute") {
					return buf, l[3]
				}
				if strings.Contains(l[1], "Query") ||
					strings.Contains(l[1], "Execute") {
					return buf, strings.Join(l[2:], ",")
				}
			case 3:
				return buf, l[2]
			default:
				if strings.Contains(l[2], "Query") ||
					strings.Contains(l[2], "Execute") {
					return buf, strings.Join(l[3:], ",")
				}
				if strings.Contains(l[1], "Query") ||
					strings.Contains(l[1], "Execute") {
					return buf, strings.Join(l[2:], ",")
				}
			}
		}
		return "", strings.Replace(buf, "\n", " ", -1) + line
	}
	if buf == "" {
		return "", ""
	}
	return "", strings.Replace(buf, "\n", " ", -1) + line
}

func state() {
	var i uint64 = 0
	for range time.NewTicker(time.Second).C {
		j := LoadUint64(&num)
		log.Printf("parse sql qps:%d", (j - i))
		i = j
	}
}
