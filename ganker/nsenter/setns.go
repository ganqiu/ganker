//go:build linux

package nsenter

/*
#define _GNU_SOURCE
#include <fcntl.h>
#include <sched.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <errno.h>
#include <string.h>

// __attribute__((constructor)) will make this function run before main() if this package is imported
__attribute__((constructor)) void enter_namespace(void) {
	char *container_pid;
	container_pid = getenv("container_pid");
	if (container_pid) {
		fprintf(stdout, "C :container_pid = %s\n", container_pid);
	}else {
		//if not set, return
		return;
	}
	char *container_cmd;
	// get the command to execute in the container
	container_cmd = getenv("container_cmd");
	if (container_cmd) {
		fprintf(stdout, "C :container_cmd = %s\n", container_cmd);
	}else {
		return;
	}
	int i;
	char nspath[1024];
	char *namespaces[] = { "ipc", "uts", "net", "pid", "mnt"};
	for (i=0; i<5; i++) {
		// join the path of namespace
		sprintf(nspath, "/proc/%s/ns/%s", container_pid, namespaces[i]);
		int fd = open(nspath, O_RDONLY);
		// invoke setns to join the namespace , if success, return 0
		if (setns(fd, 0) == -1) {
			return;
		}
		close(fd);
	}
	// 进入所有Namespace后执行指定的命令
	int res = system(container_cmd);
	exit(0);
	return;
}
*/
import "C"

func EnterNamespace() {
}
