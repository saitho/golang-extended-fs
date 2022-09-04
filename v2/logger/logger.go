package logger

type Wrapper interface {
	Debug(obj interface{})
	Error(obj interface{})
}
