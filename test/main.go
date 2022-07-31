package main

import (
	"fmt"
	xfs "github.com/saitho/golang-extended-fs"
	"github.com/saitho/golang-extended-fs/sftp"
)

type DebugLogger struct {
}

func (DebugLogger) Debug(obj interface{}) {
	fmt.Println(fmt.Sprintf("DEBUG: %s", obj))
}
func (DebugLogger) Error(obj interface{}) {
	fmt.Println(fmt.Sprintf("ERROR: %s", obj))
}

func main() {
	// Config
	sftp.Config.SshHost = "168.119.173.121"
	sftp.Config.LoadLocalSigners = true
	sftp.Config.Logger = DebugLogger{}
	remoteTestFolder := "ssh:///stackhead/testfolder"

	// CREATE FOLDER
	if err := xfs.CreateFolder(remoteTestFolder); err != nil {
		panic(err.Error())
	}

	// CHANGE OWNER
	if err := xfs.Chown(remoteTestFolder, 1412, 1412); err != nil {
		panic(err.Error())
	}

	// WRITE FILE
	if err := xfs.WriteFile(remoteTestFolder+"/test.txt", "Text"); err != nil {
		panic(err.Error())
	}
	// HAS FILE
	hasFile, err := xfs.HasFile(remoteTestFolder + "/test.txt")
	if err != nil {
		panic(err.Error())
	}
	if !hasFile {
		panic("HasFile failed. expected true.")
	}
	// READ FILE
	content, err := xfs.ReadFile(remoteTestFolder + "/test.txt")
	if err != nil {
		panic(err.Error())
	}
	if content != "Text" {
		panic("READ file failed. expected text does not match.")
	}

	// DELETE FILE
	if err := xfs.DeleteFile(remoteTestFolder + "/test.txt"); err != nil {
		panic(err.Error())
	}

	// DELETE EMPTY FOLDER
	if err := xfs.DeleteFolder(remoteTestFolder, false); err != nil {
		panic(err.Error())
	}

	// CREATE FOLDER
	if err := xfs.CreateFolder(remoteTestFolder); err != nil {
		panic(err.Error())
	}
	// COPY FILE
	if err := xfs.CopyFile("./README.md", remoteTestFolder+"/README.md"); err != nil {
		panic(err.Error())
	}
	// READ FILE
	content, err = xfs.ReadFile(remoteTestFolder + "/README.md")
	if err != nil {
		panic(err.Error())
	}
	if len(content) == 0 {
		panic("READ file failed. expected text is empty.")
	}
	// DELETE FOLDER
	if err := xfs.DeleteFolder(remoteTestFolder, true); err != nil {
		panic(err.Error())
	}
}
