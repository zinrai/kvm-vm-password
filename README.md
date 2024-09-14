# KVM VM Password Change Tool

This tool allows you to change the root password or a specific user's password for a KVM (Kernel-based Virtual Machine) virtual machine using the `virt-customize` command.

## Notes

This tool is not intended for use in production environments. It uses command-line arguments for password input, which is not secure for production use. This approach is used for simplicity and demonstration purposes only.

## Features

- Change the root password of a KVM virtual machine
- Change the password of a specific user in a KVM virtual machine
- Support for specifying a VM name or directly providing an image path
- Automatic detection and use of the first disk of the specified VM
- Verification that the specified image belongs to a known VM when using the image path option
- Ensures the VM is stopped before attempting to change passwords
- Supports both root and non-root user password changes with a single tool

## Prerequisites

- Sudo access
- virsh installed on the system
- virt-customize tool

## Installation

Build the tool:

```
$ go build
```

## Usage

The tool provides two main functionalities with options to specify either a VM name or an image path:

1. Changing the root password:
   ```
   $ ./kvm-vm-password -root -password <new_password> <vm_name>
   ```
   or
   ```
   $ ./kvm-vm-password -root -password <new_password> -image <image_path>
   ```

2. Changing a specific user's password:
   ```
   $ ./kvm-vm-password -user <username> -password <new_password> <vm_name>
   ```
   or
   ```
   $ ./kvm-vm-password -user <username> -password <new_password> -image <image_path>
   ```

Replace `<new_password>` with the desired password, `<vm_name>` with the name of your KVM virtual machine, and `<image_path>` with the path to the VM's disk image when using the -image option.

## Examples

Change root password for a VM:
```
$ ./kvm-vm-password -root -password mynewrootpass myvm
```

Change root password using image path:
```
$ ./kvm-vm-password -root -password mynewrootpass -image /path/to/vm/image.qcow2
```

Change password for user 'debian' in a VM:
```
$ ./kvm-vm-password -user debian -password mynewuserpass myvm
```

Change password for user 'debian' using image path:
```
$ ./kvm-vm-password -user debian -password mynewuserpass -image /path/to/vm/image.qcow2
```

## License

This project is open-source and available under the [MIT License](https://opensource.org/licenses/MIT).
