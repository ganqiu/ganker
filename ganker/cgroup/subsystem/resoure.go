package subsystem

const processIdPath = "tasks"

type ResourceConfig struct {
	Memory   string // Memory limit
	CpuShare string // CPU time slice allocation
	CpuQuota string // scheduled CPU quota by cfs
}

type Subsystem interface {
	Name() string                               // return the name of subsystem
	Set(path string, res *ResourceConfig) error // set the resource limit of the process
	Apply(path string, pid int) error           // add the process to the subsystem
	Delete(path string) error                   // remove the process from the subsystem
}

// Init the subsystem
var SubsystemsInit = []Subsystem{

	&MemorySubSystem{},
	&CpuShareSubSystem{},
	&CpuQuotaSubSystem{},
}
