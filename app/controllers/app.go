package controllers

import (
	"bytes"
	"fmt"
	"github.com/revel/revel"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func stdinfill(stdin io.WriteCloser) {
	var (
		err   error
		video []byte
		c     int64
	)
	video, err = ioutil.ReadFile("4.avi")
	if err != nil {
		fmt.Println("read file error : ", err)
		return
	}

	c, err = io.Copy(stdin, bytes.NewReader(video))
	if err != nil {
		fmt.Println("Copy error : ", err, " count = ", c)
		return
	}
	fmt.Println("count = ", c)
}


func (c App) GetVideo() revel.Result {
var stdout bytes.Buffer
	cmd := exec.Command("public/converter/ffmpeg1",   "-i", "-",  "-f", "webm", "-")
	cmd.Stderr = os.Stderr

	fmt.Println("!! StdinPipe")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		err = fmt.Errorf("StdinPipe ERROR : %+v\n", err)
		log.Fatal(err)
	}
	//defer func() {
	//	err = stdin.Close()
	//	if err != nil {
	//		err = fmt.Errorf("Close ERROR : %+v\n", err)
	//		log.Fatal(err)
	//	}
	//}()
	fmt.Println("!! StdoutPipe")
	cmd.Stdout = &stdout

	fmt.Println("!! Start")
	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("Start ERROR : %+v\n", err)
		log.Fatal(err)
	}

	fmt.Println("!! stdinfill")
	stdinfill(stdin)

	return c.RenderBinary(bytes.NewReader(stdout.Bytes()), "out.webm", revel.Inline, time.Now())
}
