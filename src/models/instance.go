package models

type Instance struct {
	Id            int
	Uid           int
	VolumeName    string
	ContainerName string
	InstanceType  string
	IsRunning     bool
}

func (i *Instance) CreateInstanceType(id int, userName string) *Instance {
	return &Instance{}
}
