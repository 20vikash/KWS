package models

import (
	"fmt"
	"regexp"
)

type Instance struct {
	Id            int
	Uid           int
	VolumeName    string
	ContainerName string
	InstanceType  string
	IsRunning     bool
}

func sanitize(s string) string {
	// Replace non-alphanumeric characters with underscores
	reg := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return reg.ReplaceAllString(s, "-")
}

func CreateInstanceType(uid int, userName string) *Instance {
	safeUserName := sanitize(userName)
	base := fmt.Sprintf("%d-%s", uid, safeUserName)

	const suffixInstance = "-instance"
	const suffixVolume = "_volume"
	const maxLen = 63

	maxBaseLen := maxLen - len(suffixInstance)
	if len(base) > maxBaseLen {
		base = base[:maxBaseLen]
	}
	containerName := base + suffixInstance

	maxBaseLen = maxLen - len(suffixVolume)
	volumeBase := base
	if len(base) > maxBaseLen {
		volumeBase = base[:maxBaseLen]
	}
	volumeName := volumeBase + suffixVolume

	return &Instance{
		Uid:           uid,
		VolumeName:    volumeName,
		ContainerName: containerName,
		InstanceType:  "core",
	}
}
