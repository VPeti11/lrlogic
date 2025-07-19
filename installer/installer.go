package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to clear screen:", err)
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func detectPackageManager() string {
	fmt.Println("Detecting package manager...")
	switch {
		case commandExists("apt"):
			fmt.Println("Detected apt package manager")
			return "apt"
		case commandExists("dnf"):
			fmt.Println("Detected dnf package manager")
			return "dnf"
		case commandExists("pacman"):
			fmt.Println("Detected pacman package manager")
			return "pacman"
		default:
			fmt.Println("No supported package manager found.")
			os.Exit(1)
	}
	return ""
}

func installDeps(manager string) {
	fmt.Println("Installing dependencies...")
	switch manager {
		case "apt":
			exec.Command("sudo", "apt", "update").Run()
			exec.Command("sudo", "apt", "install", "-y", "librsvg2-bin", "python3", "python3-pip", "go", "git").Run()
			fmt.Println("Installed dependencies via apt")
		case "dnf":
			exec.Command("sudo", "dnf", "install", "-y", "librsvg2", "python3", "python3-pip", "go", "git").Run()
			fmt.Println("Installed dependencies via dnf")
		case "pacman":
			exec.Command("sudo", "pacman", "-Syu", "--noconfirm", "librsvg", "python", "python-pip", "go", "git").Run()
			fmt.Println("Installed dependencies via pacman")
	}
}

func installPythonSVGPath(manager string) {
	fmt.Println("Installing python svg.path library...")
	if manager == "pacman" && commandExists("yay") {
		exec.Command("yay", "-S", "--noconfirm", "python-svg.path").Run()
		fmt.Println("Installed python-svg.path via yay")
	} else {
		exec.Command("python3", "-m", "pip", "install", "--upgrade", "pip").Run()
		args := []string{"-m", "pip", "install", "svg.path"}
		if manager == "pacman" {
			args = append(args, "--break-system-packages")
		}
		exec.Command("python3", args...).Run()
		fmt.Println("Installed svg.path via pip")
	}
}

func installPythonScript(inputPath, outputName string) {
	outputPath := filepath.Join("/usr/bin", outputName)
	os.WriteFile(outputPath, []byte("#!/usr/bin/env python3\n"), 0755)

	content, err := os.ReadFile(inputPath)
	if err == nil {
		f, _ := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY, 0755)
		defer f.Close()
		f.Write(content)
		exec.Command("chmod", "+x", outputPath).Run()
		fmt.Printf("Installed python script %s\n", outputName)
	} else {
		fmt.Printf("Failed to read %s: %v\n", inputPath, err)
	}
}

func installGoBinary(source string, outName string) {
	outPath := filepath.Join("/usr/bin", outName)
	fmt.Printf("Building Go binary %s...\n", outName)
	cmd := exec.Command("go", "build", "-o", outPath, source)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to build %s: %v\n", outName, err)
		return
	}
	exec.Command("chmod", "+x", outPath).Run()
	fmt.Printf("Installed Go binary %s\n", outName)
}

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("This installer only supports Linux.")
		return
	}

	clearScreen()
	fmt.Printf("LRLogic Installer\n")
	fmt.Println("By VPeti")
	time.Sleep(2 * time.Second)
	fmt.Println("Installing dependencies and tools...")
	time.Sleep(1 * time.Second)
	clearScreen()
	fmt.Println("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	fmt.Println("Cloning the latest source...")
	err := exec.Command("git", "clone", "https://gitlab.com/VPeti11/lrlogic.git").Run()
	if err != nil {
		fmt.Println("Failed to clone repo:", err)
		return
	}
	err = os.Chdir("lrlogic")
	if err != nil {
		fmt.Println("Failed to change directory:", err)
		return
	}

	manager := detectPackageManager()
	installDeps(manager)
	installPythonSVGPath(manager)

	installPythonScript("randomgen.py", "lrrandomgen")
	installPythonScript("svg2lrlogic/PY/transformsvg2lr.py", "pysvg2lr")

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
		return
	}

	absSource := filepath.Join(cwd, "svg2lrlogic", "Go", "main.go")
	installGoBinary(absSource, "svg2lrlogic")
	installGoBinary("main.go", "lrlogic")

	fmt.Println("All installations completed successfully.")
}
