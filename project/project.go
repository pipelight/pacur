package project

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pacur/pacur/constants"
	"github.com/pacur/pacur/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Project struct {
	Root string
}

func (p *Project) Init() (err error) {
	err = utils.MkdirAll(filepath.Join(p.Root, "mirror"))
	if err != nil {
		return
	}

	for _, release := range constants.Releases {
		err = utils.MkdirAll(filepath.Join("pkgname", release))
		if err != nil {
			return
		}
	}

	return
}

func (p *Project) getTargets() (targets []os.FileInfo, err error) {
	targets, err = ioutil.ReadDir(p.Root)
	if err != nil {
		err = &FileError{
			errors.Wrapf(err, "repo: Failed to read dir '%s'", p.Root),
		}
		return
	}

	return
}

func (p *Project) createArch(distro, release, path string) (err error) {
	archDir := filepath.Join(path, "arch")

	err = utils.Exec("", "docker", "run", "--rm", "-t", "-v",
		path+":/pacur", constants.DockerOrg+distro, "create",
		distro)
	if err != nil {
		return
	}

	err = utils.Rsync(archDir, filepath.Join(p.Root, "mirror", "arch"))
	if err != nil {
		return
	}

	err = utils.RemoveAll(archDir)
	if err != nil {
		return
	}

	return
}

func (p *Project) createRedhat(distro, release, path string) (err error) {
	yumDir := filepath.Join(path, "yum")

	err = utils.Exec("", "docker", "run", "--rm", "-t", "-v",
		path+":/pacur", constants.DockerOrg+distro+"-"+release, "create",
		distro+"-"+release)
	if err != nil {
		return
	}

	err = utils.Rsync(yumDir, filepath.Join(p.Root, "mirror", "yum"))
	if err != nil {
		return
	}

	err = utils.RemoveAll(yumDir)
	if err != nil {
		return
	}

	return
}

func (p *Project) createDebian(distro, release, path string) (err error) {
	aptDir := filepath.Join(path, "apt")

	err = utils.Exec("", "docker", "run", "--rm", "-t", "-v",
		path+":/pacur", constants.DockerOrg+distro+"-"+release, "create",
		distro+"-"+release)
	if err != nil {
		return
	}

	err = utils.Rsync(aptDir, filepath.Join(p.Root, "mirror", "apt"))
	if err != nil {
		return
	}

	err = utils.RemoveAll(aptDir)
	if err != nil {
		return
	}

	err = utils.RemoveAll(filepath.Join(path, "conf"))
	if err != nil {
		return
	}

	err = utils.RemoveAll(filepath.Join(path, "db"))
	if err != nil {
		return
	}

	return
}

func (p *Project) createTarget(target, path string) (err error) {
	distro, release := getDistro(target)
	if err != nil {
		return
	}

	switch distro {
	case "archlinux":
		err = p.createArch(distro, release, path)
	case "centos":
		err = p.createRedhat(distro, release, path)
	case "debian", "ubuntu":
		err = p.createDebian(distro, release, path)
	default:
		err = &UnknownType{
			errors.Newf("repo: Unknown repo type '%s'", target),
		}
	}

	return
}

func (p *Project) Pull() (err error) {
	for _, release := range constants.Releases {
		err = utils.Exec("", "docker", "pull", constants.DockerOrg+release)
		if err != nil {
			return
		}
	}

	return
}

func (p *Project) iterPackages(handle func(string, string) error) (err error) {
	projects, err := utils.ReadDir(p.Root)
	if err != nil {
		return
	}

	for _, project := range projects {
		if project.Name() == "mirror" || !project.IsDir() {
			continue
		}

		projectPath := filepath.Join(p.Root, project.Name())

		packages, e := utils.ReadDir(projectPath)
		if e != nil {
			err = e
			return
		}

		for _, pkg := range packages {
			err = handle(pkg.Name(), filepath.Join(projectPath, pkg.Name()))
			if err != nil {
				return
			}
		}
	}

	return
}

func (p *Project) Build() (err error) {
	err = p.iterPackages(func(release, path string) (err error) {
		err = utils.Exec("", "docker", "run", "--rm", "-t", "-v",
			path+":/pacur", constants.DockerOrg+release)
		if err != nil {
			return
		}

		return
	})
	if err != nil {
		return
	}

	return
}

func (p *Project) Repo() (err error) {
	targets, err := p.getTargets()
	if err != nil {
		return
	}

	for _, target := range targets {
		image := target.Name()
		if image == "mirror" || !target.IsDir() {
			continue
		}
		path := filepath.Join(p.Root, image)

		err = p.createTarget(image, path)
		if err != nil {
			return
		}
	}

	return
}
