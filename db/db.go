package db

//Instance DB interface for application
type Instance interface {
	Register(v interface{}) //needed for fileType. Can be left blank for other types of db
	Save(v interface{}) error
	Load(v interface{}) error
}
