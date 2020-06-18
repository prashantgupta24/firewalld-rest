package db

//Instance DB interface for application
type Instance interface {
	Register(v interface{})
	Save(v interface{}) error
	Load(v interface{}) error
}
