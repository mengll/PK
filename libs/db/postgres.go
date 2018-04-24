package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	"PK/conf"
)

type (
	Pg struct {
		Db      *sql.DB
	}
)

type Pginterface interface{
	PgConnect() error
	Pgclose()
	Prepure(str string) (*sql.Stmt,error)
	Ping() error
}

//pgconnect 处理当前的数据库路链接
func (self *Pg) PgConnect() error {

	pg := conf.PgConfAdt
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("pg connect error", err)
		}
	}()

	var err error
	port, era := strconv.Atoi(pg.Port)

	if era != nil {
		fmt.Println("端口转化错误")
		return era
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pg.Host, port, pg.User, pg.PassWord, pg.DataBase)

	self.Db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		fmt.Print("pg connect err")
		return err
	}

	erra := self.Db.Ping()

	if erra != nil {
		fmt.Println("pg connect error")
		return erra
	}
	return nil
}

//关闭当前链接
func (self *Pg) Pgclose() {
	self.Db.Close()
}

//创建预处理语句  Prepare("insert into user(name, sex)values($1,$2)")
func (self *Pg) Prepure(str string)(*sql.Stmt,error){
	Pgstmt, err := self.Db.Prepare(str)
	return Pgstmt,err
}

//检查当前是链接
func (self *Pg) Ping() error {
	err := self.Db.Ping()
	if err != nil {
		return err
	}
	return nil
}

//创建新的pg对象
func NewPg() Pginterface{
	return &Pg{}
}
