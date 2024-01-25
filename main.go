package main

import (
   "fmt"
   "os"
   "log"
   "time"
   "os/exec"
 
   "golang.org/x/sys/windows"
)

func getPid(name string) uint32 {
   h, e := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
   if e != nil { return 0 }
   p := windows.ProcessEntry32{ Size: 568 }
   for {
      e := windows.Process32Next(h, &p)
      if e != nil { return 0 }
      if windows.UTF16ToString(p.ExeFile[:]) == name {
         return p.ProcessID
      }
   }
   return 0
}

func suspendProcess(pid uint32) error {
	handle, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, pid)
	if err != nil {
	   return err
	}
	defer windows.CloseHandle(handle)

	if r1, _, _ := windows.NewLazySystemDLL("ntdll.dll").NewProc("NtSuspendProcess").Call(uintptr(handle)); r1 != 0 {
		return fmt.Errorf("NtStatus='0x%.8X'", r1)
	}
	return nil
}

func main() {
   userDir, e := os.UserHomeDir()
   if e != nil {
      log.Fatal(e)
   }
 
   launcherPath := userDir + "\\.fnlauncher\\7.40\\FortniteGame\\Binaries\\Win64\\FortniteLauncher.exe"
   eacPath := userDir + "\\.fnlauncher\\7.40\\FortniteGame\\Binaries\\Win64\\FortniteClient-Win64-Shipping_EAC.exe"
   cmd := exec.Command(eacPath, "-noeac", "-fromfl=be", "-fltoken=f7b9gah4h5380d10f721dd6a")
   launcher := exec.Command(launcherPath)
 
   launcher.Start()
   suspendProcess(uint32(launcher.Process.Pid))
   cmd.Start()
   suspendProcess(uint32(cmd.Process.Pid))

   launcher.Process.Release()
   cmd.Process.Release()

   pName := "FortniteClient-Win64-Shipping.exe"
      for {
      time.Sleep(1 * time.Second)

      pid := getPid(pName)
      _, err := os.FindProcess(int(pid))
  
      if err != nil {
         defer launcher.Process.Kill()
         defer cmd.Process.Kill()
   
         break;
      }
   }
}
