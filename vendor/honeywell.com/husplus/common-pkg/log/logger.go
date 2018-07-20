package log

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"github.com/sohlich/elogrus"
)

var defaultLogger *Logger

const (
	defaultServiceName = "Husplus"
	formatterJSON      = "json"
	formatterText      = "text"
)

func init() {
	localIP, err := GetLocalIP()
	if err != nil {
		fmt.Println(err)
	}

	defaultLogger = &Logger{
		ServiceName: getServiceName(),
		IP:          localIP,
		logger: &logrus.Logger{
			Out:       os.Stderr,
			Formatter: newFormatter(),
			Hooks:     make(logrus.LevelHooks),
			Level:     getLevel(),
		},
	}
}

func getServiceName() string {
	serviceName := os.Getenv("LOGGER_SERVICENAME")
	if serviceName == "" {
		serviceName = defaultServiceName
	}
	return serviceName
}

func getLevel() logrus.Level {
	levelstr := os.Getenv("LOGGER_LEVEL")
	if levelstr == "" {
		return logrus.InfoLevel
	}
	level, err := strconv.Atoi(levelstr)
	if err != nil {
		fmt.Printf("get log level from env error:%v,level string:%s", err, levelstr)
		return logrus.InfoLevel
	}
	return logrus.Level(level)
}

func newFormatter() logrus.Formatter {
	formatter := os.Getenv("LOGGER_FORMATTER")
	if formatter == "" {
		formatter = formatterJSON
	}
	switch formatter {
	case formatterText:
		return new(logrus.TextFormatter)
	default:
		return new(logrus.JSONFormatter)
	}
}

//Logger 环境变量名LOGGER_FORMATTER，只能设置两个值：json、text
//TODO back log :watch level change
//TODO back log :static
//TODO review code: log
//TODO load config from envirement
type Logger struct {
	ServiceName string
	IP          string
	logger      *logrus.Logger
}

func (l *Logger) newField(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{
		"serviceName": l.ServiceName,
		"ip":          l.IP,
	}

	if ctx != nil {
		//TODO get traceid
	}

	pc, filename, line, ok := runtime.Caller(2)
	if ok {
		funcname := runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")
		fields["func"] = funcname
		fields["file"] = filepath.Base(filename)
		fields["line"] = line
	}
	return fields
}

//SetServiceName 优先于环境变量
// 环境变量名LOGGER_SERVICENAME
func SetServiceName(serviceName string) {
	defaultLogger.ServiceName = serviceName
}

//SetLevel 优先于环境变量
// 环境变量名LOGGER_LEVEL,5=Debug 4=Info,3=Warn,2=Error,1=Fatal,0=Panic
func SetLevel(level logrus.Level) {
	defaultLogger.logger.SetLevel(level)
}

//UseJSONFormatter 优先于环境变量
func UseJSONFormatter() {
	defaultLogger.logger.Formatter = new(logrus.JSONFormatter)
}

//UseTextFormatter 优先于环境变量
func UseTextFormatter() {
	defaultLogger.logger.Formatter = new(logrus.TextFormatter)
}

//AddESHook AddESHook
func AddESHook(level logrus.Level, urls ...string) {
	client, err := elastic.NewClient(elastic.SetURL(urls...))
	if err != nil {
		defaultLogger.logger.Panic(nil, err)
	}

	hook, err := elogrus.NewAsyncElasticHook(client, defaultLogger.ServiceName, level, "huspluslog")
	if err != nil {
		defaultLogger.logger.Panic(nil, err)
	}
	defaultLogger.logger.Hooks.Add(hook)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Debugf(format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Infof(format, args...)
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Printf(format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Warnf(format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Warningf(format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Errorf(format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Fatalf(format, args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Panicf(format, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Debug(args...)
}

func Info(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Info(args...)
}

func Print(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Print(args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Warn(args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Warning(args...)
}

func Error(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Error(args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Fatal(args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Panic(args...)
}

func Debugln(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Debugln(args...)
}

func Infoln(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Infoln(args...)
}

func Println(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Println(args...)
}

func Warnln(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Warnln(args...)
}

func Warningln(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Warningln(args...)
}

func Errorln(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Errorln(args...)
}

func Fatalln(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Fatalln(args...)
}

func Panicln(ctx context.Context, args ...interface{}) {
	defaultLogger.logger.WithFields(defaultLogger.newField(ctx)).Panicln(args...)
}

//GetLocalIP get local private ip
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			//TODO IP6?
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", nil
}
