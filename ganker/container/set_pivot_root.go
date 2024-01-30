package container

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/sirupsen/logrus"
)

func MountSet() error {

	//set "/" as a private mount namespace's mount point
	//as pivot is prohibited if parent mount is shared

	rootfsMountFlags := syscall.MS_PRIVATE | syscall.MS_REC
	if err := syscall.Mount("", "/", "", uintptr(rootfsMountFlags), ""); err != nil {
		logrus.WithField("method", "syscall.Mount").Error(err)
		return err
	}

	// get the current dir
	pwd, err := os.Getwd()
	if err != nil {
		logrus.WithField("method", "os.Getwd").Error(err)
		return err
	}

	// mount rootfs to the current dir
	if err := pivotRoot(pwd); err != nil {
		return err
	}

	procMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

	// mount /proc to /proc
	if err := syscall.Mount("", "/proc", "proc", uintptr(procMountFlags), ""); err != nil {
		logrus.Errorf("mount proc error %v", err)
		return err
	}

	tmpfsMountFlags := syscall.MS_NOSUID | syscall.MS_STRICTATIME

	// mount /tmpfs to /sys
	if err := syscall.Mount("tmpfs", "dev", "tmpfs", uintptr(tmpfsMountFlags), "mode=755"); err != nil {
		logrus.Errorf("mount tmpfs error %v", err)
		return err
	}

	return nil
}

// it is used to pivot rootfs to a new rootfs
func pivotRoot(root string) error {

	// create a new mount point,it is different from the original mount point("/")
	// now new "root" path changes to the new mount point
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}

	// create an old_put dir to store the old rootfs
	pivotDir := filepath.Join(root, ".pivot_root")

	// check if the dir exists
	if path, err := os.Stat(pivotDir); path != nil {
		fmt.Println(err)
		if err := os.Remove(pivotDir); err != nil {
			return fmt.Errorf("remove pivot_root dir %v", err)
		}
	}

	// create temporary dir pivotDir
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}

	// it will pivot rootfs(original path) to a new rootfs(path/.pivot_root)
	// now the rootfs is the new rootfs(root)
	if err := syscall.PivotRoot(root, pivotDir); err != nil {

		return fmt.Errorf("pivot_root %v", err)
	}

	// change current dir to rootfs
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	// umount temporary dir pivotDir
	// now the rootfs is the new rootfs(root), so the path of pivotDir is also changed
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	// remove temporary dir pivotDir
	return os.Remove(pivotDir)
}
