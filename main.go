package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func isDir(path string) (bool, error) {
	isDirectory := false
	file, err := os.Open(path)

	if err == nil {
		if fi, err := file.Stat(); err == nil && fi.IsDir() {
			isDirectory = true
		}
	}

	defer file.Close()
	return isDirectory, err
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.EqualFold(a, b) {
			return true
		}
	}
	return false
}

func filesInFolder(root string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	return files, err
}

func createTextFile(root string, files []os.FileInfo) (string, error) {
	availableVideoExtensions := []string{".mov", ".mp4", ".avi", ".mkv", ".flv"}
	textFilePath := path.Join(root, "merge.txt")

	newFile, err := os.Create(textFilePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Mode().IsRegular() && stringInSlice(filepath.Ext(file.Name()), availableVideoExtensions) {
			line := fmt.Sprintf("file '%s'\n", path.Join(root, file.Name()))
			newFile.WriteString(line)
		}
	}

	defer newFile.Close()

	return textFilePath, nil
}

func executeCmd(textFile string, output string) {
	cmdTxt := fmt.Sprintf("ffmpeg -y -f concat -safe 0 -i %s -c copy %s", textFile, output)
	fmt.Printf("4- Running %s...\n", cmdTxt)

	os.Setenv("PATH", "/usr/bin:/sbin:/usr/local/bin")
	cmd := exec.Command("ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", textFile, "-c", "copy", output)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	} else {
		fmt.Printf("5- Done! %s\n", output)
	}
}

func main() {
	processedWithErrors := false

	if len(os.Args) == 1 {
		fmt.Println("Error no folders selected, select 1 folder as minimum.")
		fmt.Println("Run: ./ffmpeg_join [folder1] [folder2]...")
		os.Exit(1)
	}

	argsWithoutProg := os.Args[1:]
	for _, folder := range argsWithoutProg {
		isDirectory, err0 := isDir(folder)
		fmt.Printf("0- Check if is folder %s: %v\n", folder, isDirectory)
		processedWithErrors = processedWithErrors || err0 != nil

		if isDirectory && err0 == nil {
			fmt.Printf("1- Scanning %s...\n", folder)
			files, err1 := filesInFolder(folder)
			processedWithErrors = processedWithErrors || err1 != nil

			if err1 == nil {
				fmt.Printf("2- %v files found.\n", len(files))
				textFile, err2 := createTextFile(folder, files)
				processedWithErrors = processedWithErrors || err2 != nil

				if err2 == nil {
					fmt.Printf("3- Created text file %s\n", textFile)
					extension := filepath.Ext(files[0].Name())
					executeCmd(textFile, folder+extension)
				}
			}
		}
	}

	if processedWithErrors {
		os.Exit(2)
	}
}
