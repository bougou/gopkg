package main

import (
	"fmt"
	"log"

	// "github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	gonvml "github.com/NVIDIA/go-nvml/pkg/nvml"
)

func main() {
	ret := gonvml.Init()
	if ret != gonvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", gonvml.ErrorString(ret))
	}
	defer func() {
		ret := gonvml.Shutdown()
		if ret != gonvml.SUCCESS {
			log.Fatalf("Unable to shutdown NVML: %v", gonvml.ErrorString(ret))
		}
	}()

	count, ret := gonvml.DeviceGetCount()
	if ret != gonvml.SUCCESS {
		log.Fatalf("Unable to get device count: %v", gonvml.ErrorString(ret))
	}

	fmt.Printf("device count: %d\n", count)
	if count == 0 {
		return
	}

	for i := 0; i < int(count); i++ {
		fmt.Printf("\n\nenumerate %d device\n", i)
		device, ret := gonvml.DeviceGetHandleByIndex(i)
		if ret != gonvml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", i, gonvml.ErrorString(ret))
		}

		uuid, ret := device.GetUUID()
		if ret != gonvml.SUCCESS {
			log.Fatalf("Unable to get uuid of device at index %d: %v", i, gonvml.ErrorString(ret))
		}
		fmt.Printf("    device: %v\n", uuid)

		id, ret := device.GetBoardId()
		if ret != gonvml.SUCCESS {
			log.Fatalf("GetBoardId %d: %v", i, gonvml.ErrorString(ret))
		}
		fmt.Printf("    board id: %v\n", id)

		isMigDevice, ret := device.IsMigDeviceHandle()
		if ret != gonvml.SUCCESS {
			log.Fatalf("IsMigDeviceHandle failed, err: %v", gonvml.ErrorString(ret))
		}
		fmt.Printf("    isMigDevice: %v\n", isMigDevice)

		// MigMode: 0 (disable) 1 (enable)
		currentMode, pendingMode, ret := device.GetMigMode()
		if ret != gonvml.SUCCESS {
			log.Fatalf("GetMigMode failed, err: %v", gonvml.ErrorString(ret))
		}
		fmt.Printf("    current mig mode: %d, pending mig mode: %d\n", currentMode, pendingMode)

		if currentMode == 0 {
			fmt.Println("    mig disabled")
			continue
		}

		device.GetMaxMigDeviceCount()
		migCount, ret := device.GetMaxMigDeviceCount()
		if ret != gonvml.SUCCESS {
			log.Fatalf("Unable to get mig count")
		}
		fmt.Printf("    max mig count: %d\n", migCount)

		// 获取的是 GPU Instance Profile 的信息
		// 只能使用在 GPU Device 上，不能用在 MIG Device 上
		p1, ret := device.GetGpuInstanceProfileInfo(gonvml.GPU_INSTANCE_PROFILE_1_SLICE)
		if ret != gonvml.SUCCESS {
			log.Printf("GetGpuInstanceProfileInfo: %v", gonvml.ErrorString(ret))
		}
		fmt.Printf("    profile-1 instance count: %d\n", p1.InstanceCount)
		fmt.Printf("    profile-1 id: %d\n", p1.Id)
		fmt.Printf("    profile-1 memory slice MB: %d\n", p1.MemorySizeMB)
		fmt.Printf("    profile-1 compute slice: %d\n", p1.SliceCount)

		for j := 0; j < migCount; j++ {

			migDevice, ret := device.GetMigDeviceHandleByIndex(j)
			if ret != gonvml.SUCCESS {
				continue
			}

			uuid, ret := migDevice.GetUUID()
			if ret != gonvml.SUCCESS {
				log.Fatalf("Unable to get mig device uuid: %v", gonvml.ErrorString(ret))
			}
			fmt.Printf("\n\n    mig device: %s\n", uuid)

			isMigDevice, ret := migDevice.IsMigDeviceHandle()
			if ret != gonvml.SUCCESS {
				log.Fatalf("IsMigDeviceHandle failed, err: %v", gonvml.ErrorString(ret))
			}
			fmt.Printf("    isMigDevice: %v\n", isMigDevice)

			gi, ret := migDevice.GetGpuInstanceId()
			if ret != gonvml.SUCCESS {
				log.Fatalf("GetGpuInstanceId: %v", gonvml.ErrorString(ret))
			}
			fmt.Printf("    mig device gpu instance: %d\n", gi)

			ci, ret := migDevice.GetComputeInstanceId()
			if ret != gonvml.SUCCESS {
				log.Fatalf("GetComputeInstanceId: %v", gonvml.ErrorString(ret))
			}
			fmt.Printf("    mig device compute instance: %d\n", ci)

			mem, ret := migDevice.GetMemoryInfo()
			if ret != gonvml.SUCCESS {
				log.Fatalf("GetMemoryInfo: %v", gonvml.ErrorString(ret))
			}
			fmt.Printf("mem total: %d, mem free: %d, mem used: %d\n", mem.Total, mem.Free, mem.Used)

		}

	}

}
