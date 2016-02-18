dbtest

-------------
## 基本原理
- 收集general_log 中所有select的sql对数据库进行压力测试
- general_log_parse 格式化general_log
- dbtest 压力测试程序
- go1.5

## 安装
- git clone https://github.com/tangbinbin/dbtest.git
- make
- 执行文件在 bin 目录下

## 使用说明
    ./bin/general_log_parse -h
    Usage of ./bin/general_log_parse:
    -I string
        mysql general log file (default "/var/log/mysql/mysql.log")
    -O string
        output file (default "/tmp/out.log")

    使用示例：
    general_log_parse -I=/var/log/mysql/mysql.log -O=/tmp/output.log

    ./bin/dbtest -h
    flag needs an argument: -h
    Usage of ./bin/dbtest:
    -I string
        input file (default "aaa.log") sql日志文件(用general_log_parse 的输出)
    -c int
        max mysql connection (default 100) 最大的连接数
    -d string
        mysql database (default "test")
    -h string
        MySQL addr (default "127.0.0.1:3306")
    -n int
        max process (default 1) 压测线程数
    -p string
        passwd to connect mysql (default "test")
    -u string
        user to connect mysql (default "test")

    使用示例：
    dbtest -I=/tmp/output.log -c=100 -h=127.0.0.1:3307 -d=test -u=test -p=test -n=3 
