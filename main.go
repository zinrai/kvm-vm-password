package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type DomainDisk struct {
	XMLName xml.Name `xml:"disk"`
	Source  struct {
		File string `xml:"file,attr"`
	} `xml:"source"`
}

type Domain struct {
	XMLName xml.Name     `xml:"domain"`
	Devices struct {
		Disks []DomainDisk `xml:"disk"`
	} `xml:"devices"`
}

func main() {
	rootFlag := flag.Bool("root", false, "Change root password")
	userFlag := flag.String("user", "", "Username to change password for")
	passwordFlag := flag.String("password", "", "New password")
	imageFlag := flag.String("image", "", "Path to the KVM image")
	flag.Parse()

	if *passwordFlag == "" || (flag.NArg() == 0 && *imageFlag == "") {
		fmt.Println("Missing required arguments")
		fmt.Println("Usage: kvm-vm-password (-root | -user <username>) -password <new_password> (-image <image_path> | <vm_name>)")
		os.Exit(1)
	}

	// Check for mutually exclusive options
	if *rootFlag && *userFlag != "" {
		fmt.Println("Error: -root and -user options are mutually exclusive")
		os.Exit(1)
	}

	if !*rootFlag && *userFlag == "" {
		fmt.Println("Error: Either -root or -user must be specified")
		os.Exit(1)
	}

	var targetImage string
	var err error

	if *imageFlag != "" {
		targetImage = *imageFlag
		fmt.Printf("Using specified image: %s\n", targetImage)
		if err := verifyImageBelongsToVM(targetImage); err != nil {
			fmt.Printf("Verification failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		vmName := flag.Arg(0)
		if !isVMStopped(vmName) {
			fmt.Printf("The VM '%s' is currently running. Please stop it before changing the password.\n", vmName)
			os.Exit(1)
		}
		targetImage, err = getVMDiskPath(vmName)
		if err != nil {
			fmt.Printf("Failed to get disk path for VM '%s': %v\n", vmName, err)
			os.Exit(1)
		}
		fmt.Printf("Using disk of VM '%s': %s\n", vmName, targetImage)
	}

	if *rootFlag {
		err = changeRootPassword(targetImage, *passwordFlag)
	} else {
		err = changeUserPassword(targetImage, *userFlag, *passwordFlag)
	}

	if err != nil {
		fmt.Printf("Password change failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Password changed successfully.")
}

func changeRootPassword(imagePath, password string) error {
	args := []string{
		"virt-customize",
		"-a", imagePath,
		"--root-password", fmt.Sprintf("password:%s", password),
	}

	return runVirtCustomize(args)
}

func changeUserPassword(imagePath, username, password string) error {
	args := []string{
		"virt-customize",
		"-a", imagePath,
		"--password", fmt.Sprintf("%s:password:%s", username, password),
	}

	return runVirtCustomize(args)
}

func runVirtCustomize(args []string) error {
	fmt.Printf("Executing command: sudo %s\n", strings.Join(args, " "))

	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("sudo virt-customize command failed: %w", err)
	}

	return nil
}

func isVMStopped(vmName string) bool {
	cmd := exec.Command("sudo", "virsh", "list", "--name", "--state-running")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to get list of running VMs: %v\n", err)
		return false
	}

	runningVMs := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, vm := range runningVMs {
		if vm == vmName {
			return false
		}
	}
	return true
}

func getVMDiskPath(vmName string) (string, error) {
	cmd := exec.Command("sudo", "virsh", "dumpxml", vmName)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get VM XML: %w", err)
	}

	var domain Domain
	err = xml.Unmarshal(output, &domain)
	if err != nil {
		return "", fmt.Errorf("failed to parse VM XML: %w", err)
	}

	if len(domain.Devices.Disks) == 0 {
		return "", fmt.Errorf("no disks found for VM")
	}

	return domain.Devices.Disks[0].Source.File, nil
}

func verifyImageBelongsToVM(imagePath string) error {
	cmd := exec.Command("sudo", "virsh", "list", "--name", "--all")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get list of all VMs: %w", err)
	}

	allVMs := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, vm := range allVMs {
		disks, err := getVMDisks(vm)
		if err != nil {
			fmt.Printf("Warning: Failed to get disks for VM '%s': %v\n", vm, err)
			continue
		}
		for _, disk := range disks {
			if disk == imagePath {
				fmt.Printf("Image '%s' belongs to VM '%s'\n", imagePath, vm)
				return nil
			}
		}
	}

	return fmt.Errorf("the specified image '%s' is not connected to any known VM", imagePath)
}

func getVMDisks(vmName string) ([]string, error) {
	cmd := exec.Command("sudo", "virsh", "dumpxml", vmName)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get VM XML: %w", err)
	}

	var domain Domain
	err = xml.Unmarshal(output, &domain)
	if err != nil {
		return nil, fmt.Errorf("failed to parse VM XML: %w", err)
	}

	var disks []string
	for _, disk := range domain.Devices.Disks {
		disks = append(disks, disk.Source.File)
	}

	return disks, nil
}
