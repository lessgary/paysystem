package models

/**
* 针对项目采用模块化和面向对象
 */
import (
	//这包一定要引用，是底层的sql驱动
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"public/common"
	"public/redisClient"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

/****************数据库的相关配置*****************/
type DatabaseConf struct {
	Mysqlhost    string `json:"mysqlhost"`
	Mysqlport    string `json:"mysqlport"`
	Mysqluser    string `json:"mysqluser"`
	Mysqlpass    string `json:"mysqlpass"`
	Mysqldb      string `json:"mysqldb"`
	Mysqlcharset string `json:"mysqlcharset"`
}

//
var mysql_connStr string

/********************日期的相关格式******************/
var M_format_time string = "2006-01-02 15:04:05"
var M_format_month string = "2006-01"
var M_format_day string = "2006-01-02"
var M_format_start string = "2006-01-02 00:00:00"

/**************************************************/

//redis的超时时间
var M_redis_long_outime int = 1800

var M_redis_time int = 300

var M_redis_short_time int = 10

type DbOrm struct {
	GDb *gorm.DB
	SDb *sql.DB
}

var Gorm *DbOrm

/**
* 开始连接数据库
 */
func init() {
	db_conf := loadConfig()
	mysql_connStr = db_conf.Mysqluser + ":" + db_conf.Mysqlpass + "@tcp(" + db_conf.Mysqlhost + ":" + db_conf.Mysqlport + ")/" + db_conf.Mysqldb + "?loc=Asia%2FShanghai"
	Gorm = ConnectMysql()
}

/**
* 读取数据库配置
 */
func loadConfig() *DatabaseConf {
	dcobj := DatabaseConf{}
	file, err := os.Open("conf/database.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fd, err := ioutil.ReadAll(file)
	//json解析到结构体里面
	err = json.Unmarshal(fd, &dcobj)
	if err != nil {
		fmt.Println("database json err->", err)
		return nil
	}
	return &dcobj
}
func ConnectMysql() *DbOrm {
	db, err := gorm.Open("mysql", mysql_connStr)

	if err != nil {
		for i := 0; i < 100; i++ {
			db, err = gorm.Open("mysql", mysql_connStr)
			if err == nil && db != nil {
				break
			}
		}
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(200)
	//输出所有的sql日志
	//db.LogMode(true)
	db.DB().SetConnMaxLifetime(30)
	return &DbOrm{db, db.DB()}
}

func (d_o *DbOrm) getSqlDb() *sql.DB {
	if d_o.SDb == nil {
		d_o.SDb = d_o.GDb.DB()
	}
	return d_o.SDb
}

func GetRadomRemoval(length int, param string) string {
	//生成随机数，并且写入redis，成功即跳出循环
	var res string

	return res
}

/**
* 生成一个表的主健id = 10位时间戳+6个随机吗
 */
func GetKeyId() string {
	key_id := strconv.FormatInt(time.Now().Unix(), 10)
	rand_str := ""
	for i := 0; i < 100; i++ {
		rand_str = common.GetRadomString(6, "number")
		//将数据写入redis的set
		key := "radom_" + time.Now().Format(M_format_time)
		row := redisClient.Redis.SetAddString(key, rand_str, M_redis_short_time)

		if row == 1 {
			break
		}
	}

	key_id = key_id + rand_str
	return key_id
}
