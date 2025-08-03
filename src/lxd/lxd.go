package lxd_kws

import (
	"log"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
)

type LXDKWS struct {
	Conn lxd.InstanceServer
}

func (lxdkws *LXDKWS) PullUbuntuImage() error {
	remote, err := lxd.ConnectSimpleStreams("https://cloud-images.ubuntu.com/releases/", nil)
	if err != nil {
		log.Println("Failed to connect to lxc remote")
		return err
	}

	alias, _, err := remote.GetImageAlias("22.04")
	if err != nil {
		log.Println("Cannot get ubuntu image alias")
		return err
	}

	image, _, err := remote.GetImage(alias.Target)
	if err != nil {
		log.Println("Failed to get image information")
		return err
	}

	op, err := lxdkws.Conn.CopyImage(remote, *image, &lxd.ImageCopyArgs{
		Aliases: []api.ImageAlias{
			{Name: "ubuntu-22:04", Description: "Stable version of ubuntu cloud"},
		},
	})
	if err != nil {
		log.Println("Cannot copy remote image to local lxc")
		return err
	}

	err = op.Wait()
	if err != nil {
		log.Println("Something failed when downloading ubuntu image")
		return err
	}

	return nil
}
