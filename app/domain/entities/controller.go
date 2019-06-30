package entities

type ResouceValues struct {
	// In Megabytes
	Memory int64
	// In mili CPU
	CPU int
}
type Resource struct {
	Requests ResouceValues
	Limits   ResouceValues
}
type Controller struct {
	Resource Resource
}
