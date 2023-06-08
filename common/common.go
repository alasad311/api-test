package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

type RequestBody struct {
	reqMethod  string
	reqUri     string
	statusCode int
	clientIP   string
	hostname   string
	protocol   string
	reqAgent   string
}

var logger *zap.Logger

func LoadEnv() {
	//err := godotenv.Load("./usr/local/sbin/api.scd.edu.om/.env")
	err := godotenv.Load("./.env")
	if err != nil {
		WriteLogWithoutContext("Error couldnt load .env file "+err.Error(), zapcore.ErrorLevel, "ERROR")
	}
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
func NewRequestBody() *RequestBody {
	return &(RequestBody{})
}
func init() {
	//logFilePath := "./usr/local/sbin/api.scd.edu.om/"
	logFilePath := "./"
	logFileName := time.Now().Format("02-01-2006") + ".log"
	// log file
	fileName := path.Join(logFilePath, logFileName)
	// Write file
	_, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("err", err)
	}
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     1, // days
	})
	rawJSON := []byte(`{}`)
	var cfg zapcore.EncoderConfig
	json.Unmarshal(rawJSON, &cfg)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		w,
		zap.InfoLevel,
	)
	logger = zap.New(core)
	defer logger.Sync()
}

func LoggerInitializer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request RequestBody
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter
		// Processing request
		c.Next()
		if c.Request.Method == "POST" {
			c.Request.ParseForm()
		}
		// Request mode
		request.reqMethod = c.Request.Method
		// Request routing
		request.reqUri = c.Request.RequestURI
		// Status code
		request.statusCode = c.Writer.Status()
		// Request Agent
		request.reqAgent = c.Request.UserAgent()
		// Request IP
		request.clientIP = c.ClientIP()
		// Hostname
		request.hostname, _ = os.Hostname()
	}
}

func WriteLogWithoutContext(message string, zapCoreLevel zapcore.Level, level string) {
	logger.Log(zapCoreLevel, "",
		zap.String("datetime", time.Now().UTC().Format("2006-01-02T15:04:05-0700")),
		zap.String("app_id", "api.scd.edu.om"),
		zap.String("event", "HTTP_TRANSCATION_LOG"),
		zap.String("level", level),
		zap.String("description", message),
	)
}

func WriteToLog(c *gin.Context, message string, zapCoreLevel zapcore.Level, level string) {
	bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = bodyLogWriter
	// Processing request
	c.Next()
	if c.Request.Method == "POST" {
		c.Request.ParseForm()
	}
	// Request mode
	reqMethod := c.Request.Method
	// Request routing
	reqUri := c.Request.RequestURI
	// Status code
	statusCode := c.Writer.Status()
	// Request Agent
	reqAgent := c.Request.UserAgent()
	// Request IP
	clientIP := c.ClientIP()
	// Hostname
	hostname, error := os.Hostname()
	if error != nil {
		panic(error)
	}
	// Protocol
	protocol := c.Request.Proto

	logger.Log(zapCoreLevel, "",
		zap.String("datetime", time.Now().UTC().Format("2006-01-02T15:04:05-0700")),
		zap.String("app_id", "api.scd.edu.om"),
		zap.String("event", "HTTP_TRANSCATION_LOG"),
		zap.String("level", level),
		zap.String("description", message),
		zap.String("useragent", reqAgent),
		zap.String("source_ip", clientIP),
		zap.String("host_ip", "192.168.11.154"),
		zap.String("host_name", hostname),
		zap.String("protocol", protocol),
		zap.String("port", "80"),
		zap.String("request_uri", reqUri),
		zap.String("request_method", reqMethod),
		zap.Int("status_code", statusCode),
	)

}
func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.Query("start"))
		//if page == 0 {
		//	page = 1
		//}

		pageSize, _ := strconv.Atoi(c.Query("length"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		colNum, _ := strconv.Atoi(c.Query("order[0][column]"))
		var colName = c.Query("columns[0][data]")
		var colOrder = ""
		if colNum != 0 {
			colOrder = c.Query("order[0][dir]")
			colName = c.Query("columns[" + strconv.Itoa(colNum) + "][data]")
		}
		//offset := (page - 1) * pageSize
		return db.Offset(page).Limit(pageSize).Order(colName + " " + colOrder)
	}
}
