# kvm-vm-password

This tool allows you to change the root password or a specific user's password for a KVM (Kernel-based Virtual Machine) virtual machine using the `virt-customize` command.

## Notes

The password is read from a file via the `-password-file` option and passed to `virt-customize` using its `file:` selector, so the plaintext never appears on any command line or in the environment. The file is read by the sudo'd `virt-customize` process, so it must be readable by root. Managing the file (permissions and removal) is the caller's responsibility.

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

## Usage

The tool provides two main functionalities with options to specify either a VM name or an image path:

1. Changing the root password:
   ```
   $ ./kvm-vm-password -root -password-file <path> <vm_name>
   ```
   or
   ```
   $ ./kvm-vm-password -root -password-file <path> -image <image_path>
   ```

2. Changing a specific user's password:
   ```
   $ ./kvm-vm-password -user <username> -password-file <path> <vm_name>
   ```
   or
   ```
   $ ./kvm-vm-password -user <username> -password-file <path> -image <image_path>
   ```

Replace `<path>` with a file whose first line is the desired password, `<vm_name>` with the name of your KVM virtual machine, and `<image_path>` with the path to the VM's disk image when using the -image option.

## Examples

Write the password to a file first:
```
$ umask 077; printf '%s' mynewrootpass > /tmp/pw
```

Change root password for a VM:
```
$ ./kvm-vm-password -root -password-file /tmp/pw myvm
```

Change root password using image path:
```
$ ./kvm-vm-password -root -password-file /tmp/pw -image /path/to/vm/image.qcow2
```

Change password for user 'debian' in a VM:
```
$ ./kvm-vm-password -user debian -password-file /tmp/pw myvm
```

Change password for user 'debian' using image path:
```
$ ./kvm-vm-password -user debian -password-file /tmp/pw -image /path/to/vm/image.qcow2
```

## License

This project is licensed under the [MIT License](./LICENSE).
