/*
guicast: NATIVE WINDOWS VERSION
*/

package main

import (
	"flag"
	"fmt"
	"github.com/yasutakatou/ishell"
	"github.com/kbinani/screenshot"
	"github.com/yasutakatou/string2keyboard"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var debug bool = false
var cliHwnd uintptr
var targetHwnd uintptr

type (
	HANDLE uintptr
	HWND HANDLE
)

var (
	user32					= syscall.MustLoadDLL("user32.dll")
	procEnumWindows			= user32.MustFindProc("EnumWindows")
	procGetWindowTextW		= user32.MustFindProc("GetWindowTextW")
	procSetActiveWindow		= user32.MustFindProc("SetActiveWindow")
	procSetForegroundWindow	= user32.MustFindProc("SetForegroundWindow")
	procGetForegroundWindow = user32.MustFindProc("GetForegroundWindow")
	procGetWindowRect		= user32.MustFindProc("GetWindowRect")
)

type _RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

func GetWindowRect(hwnd HWND, rect *_RECT) (err error) {
	r1, _, e1 := syscall.Syscall(procGetWindowRect.Addr(), 7, uintptr(hwnd), uintptr(unsafe.Pointer(rect)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func main() {
  var shell=ishell.New()
 
  cliHwnd = getWindow("GetForegroundWindow")
 
  _debug := flag.Bool("debug",false,"[-debug=DEBUG MODE]")
  flag.Parse()
  debug = bool(*_debug)

  loopWait := 1
  loopCount := 1
  changeWindow := "ctrl,\\t"
  targetWindow := "Chrome"
  autoCapture := false
  captureWait := 1
  captureHang := 10
  capturePath := ""

  forcusWindow(targetWindow)

  shell.AddCmd(&ishell.Cmd{
    Name: "config",
    Help: "show config",
    Func: func(c *ishell.Context) {
      fmt.Println("wait:" + strconv.Itoa(loopWait))
      fmt.Println("count:" + strconv.Itoa(loopCount))
      fmt.Println("changeWindow:" + changeWindow)
      fmt.Println("targetWindow:" + targetWindow)
      fmt.Printf("autoCapture mode: %t\n", autoCapture)
      fmt.Println("captureWait:" + strconv.Itoa(captureWait))
      fmt.Println("captureHang:" + strconv.Itoa(captureHang))
      fmt.Println("capturePath:" + capturePath)
    },
  })

  shell.AddCmd(&ishell.Cmd{
    Name: "wait",
    Help: "cast loop wait(must >= 1). usage: wait (wait count)",
    Func: func(c *ishell.Context) {
      if len(c.Args) < 1 {
        fmt.Printf("usage: wait (wait count)\n")
        return
      }

      cnt, err := strconv.Atoi(c.Args[0])
      if cnt < 1 {
        fmt.Printf("cast loop wait(must >= 1)\n")
        return
      }
      if err != nil {
        log.Fatal(err)
      }
      loopWait = cnt
    },
  })

  shell.AddCmd(&ishell.Cmd{
    Name: "count",
    Help: "cast loop count(must >= 1). usage: count (loop count)",
    Func: func(c *ishell.Context) {
      if len(c.Args) < 1 {
        fmt.Printf("usage: count (loop count)\n")
        return
      }

      cnt, err := strconv.Atoi(c.Args[0])
      if cnt < 1 {
        fmt.Printf("cast loop count(must >= 1)\n")
        return
      }
      if err != nil {
        log.Fatal(err)
      }
      loopCount = cnt
    },
  })

  shell.AddCmd(&ishell.Cmd{
    Name: "target",
    Help: "change target window. usage: target (WINDOW NAME)",
    Func: func(c *ishell.Context) {
      if len(c.Args) < 1 {
        fmt.Printf("usage: target (WINDOW NAME)\n")
        return
      }

      tool := c.Args[0]
      if len(c.Args) > 1 {
        for i := 1; i < len(c.Args); i++ {
          tool += " "
          tool += c.Args[i]
        }
      } 

      targetWindow = tool
      forcusWindow(targetWindow)
    },
  })

  shell.AddCmd(&ishell.Cmd{
    Name: "capturePath",
    Help: "change capture save Path. usage: capturePath (PATH)",
    Func: func(c *ishell.Context) {
      if len(c.Args) < 1 {
        fmt.Printf("usage: capturePath (PATH)\n")
        return
      }

      tool := c.Args[0]
      if len(c.Args) > 1 {
        for i := 1; i < len(c.Args); i++ {
          tool += " "
          tool += c.Args[i]
        }
      } 

      capturePath = tool
    },
  })

  shell.AddCmd(&ishell.Cmd{
    Name: "change",
    Help: "change window for cast. usage: change (keyboard shortcut)",
    Func: func(c *ishell.Context) {
      if len(c.Args) < 1 {
        fmt.Printf("usage: change [keyboard shortcut]\n")
        return
      }
      changeWindow = c.Args[0]
    },
  })

  shell.AddCmd(&ishell.Cmd{
    Name: "autoCapture",
    Help: "autoCapture mode switch",
    Func: func(c *ishell.Context) {
      if autoCapture == true {
        autoCapture = false
      } else {
        autoCapture = true
      }
      fmt.Printf("autoCapture mode: %t\n", autoCapture)
    },
  })

  shell.AddCmd(&ishell.Cmd{
    Name: "list",
    Help: "all window listing",
    Func: func(c *ishell.Context) {
      listWindow()
    },
  })

  shell.AddCmd(&ishell.Cmd{
    Name: "captureWait",
    Help: "captureWait loop count(must >= 1). usage: captureWait (loop count)",
    Func: func(c *ishell.Context) {
      if len(c.Args) < 1 {
        fmt.Printf("usage: captureWait (loop count)\n")
        return
      }

      cnt, err := strconv.Atoi(c.Args[0])
      if cnt < 1 {
        fmt.Printf("captureWait loop count(must >= 1)\n")
        return
      }
      if err != nil {
        log.Fatal(err)
      }
      loopWait = cnt
    },
  })

  shell.AddCmd(&ishell.Cmd{
    Name: "captureHang",
    Help: "captureHang loop hangup count(must >= 1). usage: captureHang (hangup count)",
    Func: func(c *ishell.Context) {
      if len(c.Args) < 1 {
        fmt.Printf("usage: captureHang (hangup count)\n")
        return
      }

      cnt, err := strconv.Atoi(c.Args[0])
      if cnt < 1 {
        fmt.Printf("captureHang hangup count(must >= 1)\n")
        return
      }
      if err != nil {
        log.Fatal(err)
      }
      loopWait = cnt
    },
  })

  shell.AddCmd(&ishell.Cmd{
    Name: "onlyCapture",
    Help: "do Capture only",
    Func: func(c *ishell.Context) {
      for i := 0; i < loopCount; i++ {
        t := time.Now()
        const layout = "2006-01-02-15-04-05"
        captureFile := strconv.Itoa(i+1) + "_" + t.Format(layout) + ".png"

        getScreenCapture(capturePath + captureFile)

        hangCount := 0
        for {
          time.Sleep(time.Duration(captureWait) * time.Second)
          if Exists(capturePath + captureFile) == true {
            fmt.Println("Capture Success: " + captureFile + " (" + Filesize(capturePath + captureFile) + ")")
            break
          }

          if hangCount > captureHang {
            fmt.Println("Capture Failure!!")
            break
          }
          hangCount = hangCount + 1
        }

        sendKeyOrString(false, changeWindow, "")
      }
    },
  })

  shell.AddCmd(&ishell.Cmd{
    Name: "!",
    Help: "no capture",
    Func: func(c *ishell.Context) {

      doCmd := c.Args[0]
      if len(c.Args) > 1 {
        for i := 1; i < len(c.Args); i++ {
          doCmd += "|"
          doCmd += c.Args[i]
        }
      } 
      do(capturePath,changeWindow,targetWindow,doCmd,loopCount,loopWait,captureHang,captureWait,false)
    },
  })

  shell.AddCmd(&ishell.Cmd{
    Name: "default",
    Help: "default input is cast",
    Func: func(c *ishell.Context) {

      doCmd := c.Args[0]
      if len(c.Args) > 1 {
        for i := 1; i < len(c.Args); i++ {
          doCmd += "|"
          doCmd += c.Args[i]
        }
      }
      do(capturePath,changeWindow,targetWindow,doCmd,loopCount,loopWait,captureHang,captureWait,autoCapture)
    },
  })

  shell.Run()
}


func SetActiveWindow(hwnd HWND) {
	syscall.Syscall(procSetActiveWindow.Addr(), 4, uintptr(hwnd),0 ,0)
	syscall.Syscall(procSetForegroundWindow.Addr(), 5, uintptr(hwnd),0 ,0)
}

func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func getWindow(funcName string) uintptr {
	hwnd, _, _ := syscall.Syscall(procGetForegroundWindow.Addr(), 6, 0 ,0 ,0)
	if debug == true {
		fmt.Printf("currentWindow: handle=0x%x\n", hwnd)
	}
    return hwnd
}


func save(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

func listWindow() {
	var rect _RECT

	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		
		GetWindowRect(HWND(h),&rect)
		if debug == true {
			fmt.Printf("Verbose Window Title '%s' window: handle=0x%x ", syscall.UTF16ToString(b), h)
			fmt.Printf("window rect: ")
			fmt.Println(rect)
		}
		
		if (int(rect.Left) > 0 || int(rect.Top) > 0) {
			fmt.Printf("Window Title '%s' window: handle=0x%x\n", syscall.UTF16ToString(b), h)
		}
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)
}

func forcusWindow(title string) int {
	var hwnd syscall.Handle

	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}

		if debug == true {
			fmt.Printf("EnumWindows Search '%s' window: handle=0x%x\n", syscall.UTF16ToString(b), h)
		}
		
		if strings.Index(syscall.UTF16ToString(b), title) != -1 && fmt.Sprintf("%x",h) != fmt.Sprintf("%x",cliHwnd) {
			if debug == true {
				fmt.Printf("Found '%s' window: handle=0x%x\n", syscall.UTF16ToString(b), h)
			}

			// note the window
			hwnd = h
			targetHwnd = uintptr(h)
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)
	if hwnd == 0 {
		fmt.Printf("No window with title '%s' found, set exists windows title\n", title)
		return 1
	}
	return 0
}

func do(capturePath,changeWindow,targetWindow,doCmd string,loopCount,loopWait,captureHang,captureWait int, autoCapture bool) {
  for i := 0; i < loopCount; i++ {
    if strings.Index(doCmd,"|") != -1 {
      params := strings.Split(doCmd, "|")
      for r := 0; r < len(params); r++ {
        if strings.Index(params[r],"wait") != -1 {
          waits := strings.Replace(params[r], "wait", "", 1)
          cnt, _ := strconv.Atoi(waits)
          fmt.Println("wait: " + waits)
          time.Sleep(time.Duration(cnt) * time.Second)
        } else if strings.Index(params[r],"[") != -1 && strings.Index(params[r],"]") != -1 {
          keya := strings.Replace(params[r], "[", "", 1)
          keyb := strings.Replace(keya, "]", "", 1)

          sendKeyOrString(false, keyb, "")
        } else {
          sendKeyOrString(true, "", params[r])
        }
      }
    } else {
      if strings.Index(doCmd,"wait") != -1 {
        waits := strings.Replace(doCmd, "wait", "", 1)
        cnt, _ := strconv.Atoi(waits)
        fmt.Println("wait: " + waits)
        time.Sleep(time.Duration(cnt) * time.Second)
      } else if strings.Index(doCmd,"[") != -1 && strings.Index(doCmd,"]") != -1 {
        keya := strings.Replace(doCmd, "[", "", 1)
        keyb := strings.Replace(keya, "]", "", 1)

        sendKeyOrString(false, keyb, "")
      } else {
        sendKeyOrString(true ,"", doCmd)
      }
    }

    time.Sleep(time.Duration(loopWait) * time.Second)

    if autoCapture == true {
      t := time.Now()
      const layout = "2006-01-02-15-04-05"
      captureFile := strconv.Itoa(i+1) + "_" + t.Format(layout) + ".png"

      getScreenCapture(capturePath + captureFile)

      hangCount := 0
      for {
        time.Sleep(time.Duration(captureWait) * time.Second)
        if Exists(capturePath + captureFile) == true {
          fmt.Println("Capture Success: " + captureFile + " (" + Filesize(capturePath + captureFile) + ")")
          break
        }

        if hangCount > captureHang {
          fmt.Println("Capture Failure!!")
          break
        }
        hangCount = hangCount + 1
      }
    }

    sendKeyOrString(false, changeWindow, "")
  }
}

func Filesize(filename string) string {
  fileinfo, staterr := os.Stat(filename)

  if staterr != nil {
    fmt.Println(staterr)
    return "Error!"
  }

  return strconv.FormatInt(fileinfo.Size(),10)
}

func Exists(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil
}

func getScreenCapture(fileName string) {
	if targetHwnd != getWindow("GetForegroundWindow") {
		SetActiveWindow(HWND(targetHwnd))
		time.Sleep(500)
	}

	var rect _RECT
	GetWindowRect(HWND(targetHwnd),&rect)
	if debug == true {
		fmt.Printf("window rect: ")
		fmt.Println(rect)
	}

	img, err := screenshot.Capture(int(rect.Left) + 10, int(rect.Top) + 10, int(rect.Right-rect.Left) - 20, int(rect.Bottom-rect.Top) - 20)
	if err != nil {
	  panic(err)
	}
	save(img, fileName)
}

func sendKeyOrString(keyStr bool, sendKeys string, sendStrs string) {
	if targetHwnd != getWindow("GetForegroundWindow") {
		SetActiveWindow(HWND(targetHwnd))
		time.Sleep(500)
	}

	cCtrl := false
	cAlt := false
	var splitResult []string

	if keyStr == true {
		string2keyboard.KeyboardWrite(sendStrs, false, false)
		if debug == true {
			fmt.Printf("StringInput: ")
			fmt.Println(sendStrs)
		}
	} else {
		splitResult = strings.Split(sendKeys, ",")
		if debug == true { fmt.Printf("KeyInput: ") }
		for i := 0; i < len(splitResult); i++ {
			if splitResult[i] == "ctrl" {
				cCtrl = true
				if debug == true { fmt.Printf("ctrl + ") }
			} else if splitResult[i] == "alt" {
				cAlt = true
				if debug == true { fmt.Printf("alt + ") }
			} else {
				if debug == true { fmt.Println(splitResult[i]) }
				string2keyboard.KeyboardWrite(splitResult[i], cCtrl, cAlt)
			}
		}
	}
}

