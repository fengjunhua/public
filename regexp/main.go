package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"regexp"
	"time"
)

func main() {

	//file, err := ioutil.ReadFile("/home/fengjh/.config/powermanagementprofilesrc")
	//if err != nil {
	//	log.Println(err)
	//}

	//fmt.Println(string(file))
	//compile := regexp.MustCompile("\\[AC\\]\\[BrightnessControl\\]\n.*value=(\\d+)\n")

	//findString := compile.FindString(string(file))
	//fmt.Println(findString)

	//ret := compile.FindStringSubmatch(string(file))

	//fmt.Println(len(ret))
	//fmt.Println(ret[1])

	//compile2 := regexp.MustCompile("\\d+")
	//s := compile2.FindString(findString)

	//fmt.Println(s)

	//value := GetBrightnessSystemSetting()

	//fmt.Println("value is :", value)

	go MonitorBrightnessSystemSetting()

	func() {
		for {
			time.Sleep(1 * time.Minute)
		}
	}()
}

func MonitorBrightnessSystemSetting() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err1 := watcher.Add("/home/fengjh/.config/")
	if err1 != nil {
		log.Fatal(err1)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//log.Println("event:", event)
				//log.Println("event name:", event.Name)
				//log.Println("event Op:", event.Op)
				//log.Println("event type:", reflect.TypeOf(event.Op))

				if event.Name == "/home/fengjh/.config/powermanagementprofilesrc" && event.Op == fsnotify.Create {
					//log.Printf("modified file:", event.Name, event.Op)
					//log.Printf("modified file:", event.Name)
					//log.Println("true")
					log.Println(event.Name, event.Op)
					value := GetBrightnessSystemSetting()
					fmt.Println("value is :", value)
				}

			case err2, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err2)
			}
		}
	}()

	<-done
}

func GetBrightnessSystemSetting() string {

	var value string
	file, err := ioutil.ReadFile("/home/fengjh/.config/powermanagementprofilesrc")
	if err != nil {
		log.Println(err)
	}

	//fmt.Println(string(file))
	compile := regexp.MustCompile("\\[AC\\]\\[BrightnessControl\\]\n.*value=(\\d+)\n")

	ret := compile.FindStringSubmatch(string(file))

	if len(ret) == 0 {
		value = "null"
	} else {
		value = ret[1] + "%"
	}

	return value
}
