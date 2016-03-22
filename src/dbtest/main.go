package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	. "sync/atomic"
	"syscall"
	"time"
)

var (
	input           = flag.String("I", "aaa.log", "input file")
	addr            = flag.String("h", "127.0.0.1:3306", "MySQL addr")
	user            = flag.String("u", "test", "user to connect mysql")
	passwd          = flag.String("p", "test", "passwd to connect mysql")
	database        = flag.String("d", "test", "mysql database")
	conn            = flag.Int("c", 100, "max mysql connection")
	thread          = flag.Int("n", 1, "max process")
	num      uint64 = 0
	num2     uint64 = 0
	db       *sql.DB
)

func init() {
	flag.Parse()
	initDb()
}

func initDb() {
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&timeout=100ms",
		*user, *passwd, *addr, *database,
	)
	var err error
	db, err = sql.Open("mysql", connStr)
	db.SetMaxOpenConns(*conn)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go state()
	ch := make(chan string, 1000)

	for i := 0; i < *thread; i++ {
		go func() {
			for {
				sql := <-ch
				if exec(sql) {
					AddUint64(&num, 1)
					continue
				}
				AddUint64(&num2, 1)
			}
		}()
	}

	f, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	br := bufio.NewReader(f)
	go func() {
		for {
			line, err := br.ReadString('\n')
			if err != nil && err == io.EOF {
				log.Fatal(err)
			}
			ch <- line
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	db.Close()
}

func exec(sql string) bool {
	rows, err := db.Query(sql)
	if err != nil {
		return false
	}
	rows.Close()
	return true
}

func state() {
	var (
		i  uint64 = 0
		i2 uint64 = 0
	)
	for range time.NewTicker(time.Second).C {
		j := LoadUint64(&num)
		j2 := LoadUint64(&num2)
		log.Printf("execute sql qps:%d  err sql qps:%d",
			(j - i), (j2 - i2))
		i, i2 = j, j2
	}
}
