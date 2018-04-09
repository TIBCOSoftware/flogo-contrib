//+build ignore

package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	fmt.Println("Running build script for the Lambda trigger")

	var cmd = exec.Command("")

	// appdir is the directory where the app is stored. For example if you app is called
	// lambda this would be <path>/lambda/src/lambda
	appDir := os.Args[1]

	// Clean up
	fmt.Println("Cleaning up previous executables")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "del", "/q", "handler", "handler.zip")
	} else {
		cmd = exec.Command("rm", "-f", "handler", "handler.zip")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = appDir

	err := cmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}

	// Build an executable for Linux
	fmt.Println("Building a new handler file")
	cmd = exec.Command("go", "build", "-o", "handler")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = appDir
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOPATH=%s", filepath.Join(appDir, "..", "..")), "GOOS=linux")

	err = cmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}

	// Zip the executable using the same code as build-lambda-zip
	fmt.Println("Zipping the new handler file")
	inputExe := filepath.Join(appDir, "handler")
	outputZip := filepath.Join(appDir, "handler.zip")
	if err := compressExe(outputZip, inputExe); err != nil {
		fmt.Printf("Failed to compress file: %v", err)
	}
}

func writeExe(writer *zip.Writer, pathInZip string, data []byte) error {
	exe, err := writer.CreateHeader(&zip.FileHeader{
		CreatorVersion: 3 << 8,     // indicates Unix
		ExternalAttrs:  0777 << 16, // -rwxrwxrwx file permissions
		Name:           pathInZip,
		Method:         zip.Deflate,
	})
	if err != nil {
		return err
	}

	_, err = exe.Write(data)
	return err
}

func compressExe(outZipPath, exePath string) error {
	zipFile, err := os.Create(outZipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	data, err := ioutil.ReadFile(exePath)
	if err != nil {
		return err
	}

	return writeExe(zipWriter, filepath.Base(exePath), data)
}
