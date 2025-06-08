package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Instance struct {
	Id            int
	Uid           int
	VolumeName    string
	ContainerName string
	InstanceType  string
	IsRunning     bool
}

func (i *Instance) CreateInstanceType(uid int, userName string) *Instance {
	s := fmt.Sprintf("%d:%s", uid, userName)
	h := sha256.Sum256([]byte(s))
	hashString := hex.EncodeToString(h[:])

	containerName := hashString + "_instance"
	volumeName := hashString + "_volume"
	instanceType := "core"

	return &Instance{
		Uid:           uid,
		VolumeName:    volumeName,
		ContainerName: containerName,
		InstanceType:  instanceType,
	}
}
