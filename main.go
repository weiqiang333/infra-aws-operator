// author: weiqiang; date: 2022-12
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"infra-aws-operator/internal/api"
	"infra-aws-operator/internal/databases"
	"infra-aws-operator/internal/databases/infrastructure"
	"infra-aws-operator/internal/metrics"
	"infra-aws-operator/pkg/job"
)

func init() {
	pflag.String("configFile", "configs/config.yaml", "go config file")
	pflag.String("listen_address", "0.0.0.0:8088", "server listen address.")
	pflag.String("databases_type", "mysql", "databases sqlite/mysql")
}

func main() {
	loadConfig()

	// conn db
	dbCoon := databases.NewDBConns()
	dbType := viper.GetString("databases_type")
	//if dbType == "sqlite" {
	//dbFile := viper.GetString("databases.sqlite.dbfile")
	//connmaxlifetime := time.Duration(viper.GetInt64("databases.sqlite.connmaxlifetime"))
	//maxopenconns := viper.GetInt("databases.sqlite.maxopenconns")
	//if err := dbCoon.ConnsSqlite(dbFile, connmaxlifetime, maxopenconns); err != nil {
	//	panic(err)
	//}
	//}
	if dbType == "mysql" {
		user := viper.GetString("databases.mysql.crud.user")
		password := viper.GetString("databases.mysql.crud.password")
		address := viper.GetString("databases.mysql.crud.address")
		dbName := viper.GetString("databases.mysql.crud.dbName")
		connmaxlifetime := time.Duration(viper.GetInt64("databases.sqlite.connmaxlifetime"))
		maxopenconns := viper.GetInt("databases.sqlite.maxopenconns")
		if err := dbCoon.ConnsMysql(user, password, address, dbName, connmaxlifetime, maxopenconns); err != nil {
			panic(err)
		}
	}
	if err := infrastructure.CreateTables(dbCoon.DBCRUD); err != nil {
		panic(err)
	}

	// job
	job.Job(dbCoon.DBCRUD)

	prometheus.MustRegister(metrics.NewExporter(dbCoon.DBCRUD))

	listenAddress := viper.GetString("listen_address")
	router := engine(dbCoon.DBCRUD)
	err := router.Run(listenAddress) // listen and serve on 0.0.0.0:8080
	if err != nil {
		panic(fmt.Errorf("Failed web server: %s ", err.Error()))
	}
}

// gin web run engine
func engine(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.LoadHTMLGlob("web/templates/*")
	router.Static("/static", "./web/static")

	router.GET("/check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})
	router.POST("/-/reload", reloadConfig)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.GET("/", api.Default)

	userAuth := loadAuthUsers()
	authorized := router.Group("/api/v1", gin.BasicAuth(userAuth))
	authorized.GET("/", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		c.String(200, "asd", user)
	})

	return router
}

// load config and flag config
func loadConfig() {
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		fmt.Println(err.Error())
		panic(fmt.Errorf("Fatal error BindPFlags: %w \n", err))
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(viper.GetString("configFile"))
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}

// reloadConfig 127.0.0.1:8080/-/reload
func reloadConfig(c *gin.Context) {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Println(fmt.Errorf("Fatal error config file: %w \n", err))
		c.String(http.StatusOK, fmt.Sprintf("Failed reload config file: %s, err: %s", viper.ConfigFileUsed(), err.Error()))
		return
	}
	fmt.Println("reload config file: ", viper.ConfigFileUsed())
	c.String(http.StatusOK, fmt.Sprintf("reload config file: %s", viper.ConfigFileUsed()))
}

func loadAuthUsers() map[string]string {
	return viper.GetStringMapString("auth.basic")
}
