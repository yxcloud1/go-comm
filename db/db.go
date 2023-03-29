package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type dataService struct {
	conn           *gorm.DB
	sqlDB          *sql.DB
	connectionType string
	connectionUrl  string
}

var (
	db = &dataService{
		conn:           nil,
		sqlDB:          nil,
	}
)

func SetOption(connectionType string, connectionUrl string){
	db.connectionType = connectionType
	db.connectionUrl = connectionUrl
}

func DB() *dataService {
	if db.conn == nil {
		db.Start()
	}
	return db
}

func (s *dataService) Start() error {
	err := s.open()
	if err != nil {
		log.Println("DB.open() error:", err)
	} else {
		log.Println(s.connectionUrl + " connected !!")
	}
	return err
}

func (s *dataService) Conn() *gorm.DB {
	return s.conn
}

func (s *dataService) open() error {
	if s.conn != nil {
		return nil
	}

	if s.connectionType == "" {
		s.connectionType = "pg"
	}
	if s.connectionUrl == "" {
		return errors.New("缺少连接字符串！！")
	}

	var conn *gorm.DB = nil
	var err error
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},

		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logger.Silent),
	}
	if s.connectionType == "mssql" {
		conn, err = gorm.Open(sqlserver.Open(s.connectionUrl), config)
	} else if s.connectionType == "pg" {
		conn, err = gorm.Open(postgres.Open(s.connectionUrl), config)

	} else if s.connectionType == "mysql" {
		conn, err = gorm.Open(mysql.Open(s.connectionUrl), config)
	} else {
		return errors.New("不支持的数据库驱动 " + s.connectionType)
	}


	if err != nil {
		s.conn = nil
		return err
	}

	s.conn = conn
	sqlDB, err := s.conn.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Ping(); err != nil {
		return err
	}
	s.sqlDB = sqlDB
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return nil
}

func (s *dataService) ExecuteSQL(command string, params ...interface{}) error {
	tx := s.conn.Exec(command, params...)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s *dataService) ExecuteBatchSQL(command string, params [][]interface{}) error {
	for _, pm := range params {
		tx := s.conn.Exec(command, pm...)
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}

func (s *dataService) ExecuteQuery(command string, params ...interface{}) ([]map[string]interface{}, error) {
	tx := s.conn.Raw(command, params...)
	if tx.Error != nil {
		return nil, tx.Error
	}
	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	result := make([]map[string]interface{}, 0)
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	row := make([]interface{}, len(cols))
	for i, _ := range cols {
		var a interface{}
		row[i] = &a
	}

	for rows.Next() {
		{
			err = rows.Scan(row...)
			if err != nil {
				fmt.Println(err)
			} else {
				rowmap := make(map[string]interface{})
				for k, c := range cols {
					rowmap[c] = *row[k].(*interface{})
				}
				result = append(result, rowmap)
			}
		}
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}
