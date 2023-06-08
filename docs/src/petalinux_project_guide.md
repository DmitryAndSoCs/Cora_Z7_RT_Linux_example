---
title: "Cora Z7-10 petalinux project"
author: [Dmitrii Matafonov]
date: 
subtitle: "Detailed guide on petalinux project creation"
keywords: [FPGA, Xilinx, vhdl, documentation]
titlepage: true
titlepage-text-color: "000000"
titlepage-rule-color: "CCCCCC"
titlepage-rule-height: 4
logo: ../assets/os_icon.png
logo-width: 100 
page-background:
page-background-opacity:
titlepage-background: ../assets/title_page.png
links-as-notes: true
lot: false
lof: false
listings-disable-line-numbers: false
listings-no-page-break: false
disable-header-and-footer: false
header-left:
header-center:
header-right:
footer-left: "Dmitrii Matafonov (@DmitryAndSoCs)"
footer-center:
footer-right:
subparagraph: true
lang: en-US
...

# Introduction

This guide is intended to give a step-by-step guidelines on how to create a linux project based on the *.xsa file exported from Vivado. It would automatically configue base drivers and kernel setup, except for the custom IP module. The guidelines are going to give examples on how to write a kernel module to handle custom IP's interrupt and route it to a software. 

## Tools

Petalinux tools can be used only on linux machines. It works well with WSL2 on Windows, but the project must be in the EXT4 local virtual hard drive, or it would as slow as if you were running it on a machine from the 80s. 

> _**Note**_: Native installations support Ubuntu 20.04 maximum. You WILL experience a lot of issues if you run it on 22.04. 

If you run the newest one, a pre-made docker image might be a good solution for you: [petalinux-docker](https://github.com/carlesfernandez/docker-petalinux2). The installation is very clear and simple.

For native installations the links are down below. 

Petalinux 2021.2: [download link](https://www.xilinx.com/support/download/index.html/content/xilinx/en/downloadNav/embedded-design-tools/archive.html)

Required packages: [download link](https://support.xilinx.com/s/article/73296?language=en_US)y

_Minor suggestion_: I would recommend installing mc (Midnight commander) for the convenience of editing and navigating in the console. It works great on your machine and on the embedded linux we would be building. 

```bash {caption="Optional: mc installation"}
sudo apt-get install mc
```

How to install it on the board will be covered in the scope of this document.

## SD card

SD card has to be formatted to FAT32 filesystem. 512 Mb SD card is usually enough, but I use 4, 8 or 16 Gb ones. 
If one wants the linux image to have persistent storage (save files and changes across reboots), they need to make 2 partitions: FAT32 for bootloader and EXT4 for rootfs. Detailed instructions are [here](https://docs.xilinx.com/r/2021.2-English/ug1144-petalinux-tools-reference-guide/Partitioning-and-Formatting-an-SD-Card). The build in repository utilizes InitRAM. 

## External kernel sources

Since this project is aimed at creating a Real-time patched version of Xilinx kernel, it is necessary to download the external kernel sources package to be able to patch it and build easily without messing with the original petalinux installation.

[Linux Xilinx kernel 2021.2](https://github.com/Xilinx/linux-xlnx/releases/tag/xilinx-v2021.2).

For more specific applications, you may want to utilize different kernels from other providers, such as Analog devices, they provide additional kernel drivers for their own products, it may be useful for automotive applications. 

The patching approach remains the same for those versions of kernel. 

# Petalinux project development

Before using petalinux tools, one needs to source the settings file located in petalinux installation dir. 

```bash {caption="Sourcing petalinux settings"}
source ~/petalinux/settings.sh # adjust to your petalinux installation folder
```

## Petalinux tools quick reference

```bash {caption="Quick reference for petalinux commands"}
# To create a new folder with initial petalinux project template use
petalinux-create --type project --template zynq --name *your_project_name*

# To make initial linux project autoconfig from the HDL design use 
petalinux-config --get-hw-description ~/*your_project_folder*
# After that you can just use petalinux-config to edit the existing settings

# /*To change kernel config (add kernel drivers,FS system, kernel tweaks) use */
petalinux-config -c krenel

# To change the contents of the build's rootfs use 
petalinux-config -c rootfs

# To change the build's u-boot settings use
petalinux-config -c u-boot

# To build the project use
petalinux-build

# To fully clean the build use 
petalinux-build -x mrproper

# To build some component (kernel, for example) separately use
petalinux-build -c kernel
# or u-boot, rootfs, u-boot, device-tree

# To pack the build into writable files under ??/images/linux/
petalinux-package --boot --fsbl --fpga --u-boot --force

# In case of InitRAM pacakging type use
petalinux-package --boot --fsbl --fpga --u-boot --kernel --force
# to embed everything into a single BOOT.BIN

# To boot linux in qemu emulator
petalinux-package --prebuilt --force
petalinux-boot --qemu --prebuilt 3

# If necessary, you can try to boot the whole image over JTAG using
petalinux-package --prebuilt --force
petalinux-boot --jtag --prebuilt 3 
# It's a VERY long process
```

All the changes to the default tamplate are stored in `project-spec`.

|         Path          | Description |
|:---------------------:|:-----------:|
| *project_folder*/project-spec/meta-user/recipes-apps  | **sources for user applications**|
| *project_folder*/project-spec/meta-user/recipes-bsp | **sources for device-tree, u-boot and FSBL** |
| *project_folder*/project-spec/meta-user/recipes-kernel | **sources for linux kernel changes** |
| *project_folder*/project-spec/meta-user/recipes-modules | **sources for linux kernel modules** |


## New project

To create a new project, make a folder where you'd want to keep the projects. If one wants to do this in this repository, they may use a command similar to this:

```bash {caption="Creating a new project"}
cd ~/Cora_Z7/linux_build_cora # it assumes the repository is stored in your home dir
petalinux-create --type project --template zynq --name your_project_name # choose a new name for the project
```

> _**Note**_: this step is the same for custom boards as well.

To make the default configuration based on the design created at the previous step in Vivado, issue:

```bash {caption="First configuration based on Vivado design"}
petalinux-config --get-hw-description ~/Cora_Z7/linux_build_cora/hw # adjust to the used path for *.xsa file
```

This is the "main", overall configuration menu. It has the basic settings for with adjustments for the specific Vivado design already in place. At this staep no changes are needed, but it would be nice for the reader to make themselves familiar with the settings. 

Click `ESC` 2 times to exit the settings menu and let the tool finish the configuration.

It is not necessary for many simple projects to have an external kernel, but for this course because of the Real-Time patch it is going to be needed. 

Download Xilinx krenel for 2021.2: [link](https://github.com/Xilinx/linux-xlnx/releases/tag/xilinx-v2021.2).

Untar it somewhere ( `tar -xvzf ./linux-xlnx-xilinx-v2021.2` to untar locally), I suggest using the `linux_cora_build` folder, one level above from petalinux project.

Issue `petalinux-config` and navigate to `Linux components selection`. 

1. Choose `linux-kernel`.
2. Choose `ext-local-src`.
3. The new menu entry will appear called `External linux-kernel local source settings`. Choose it.
4. Edit `External linux-kernel local source path` to point to the untared linux sources.
6. _Optional_. If persistent storage is needed (i.e. files are preserved across reboots, changes are permanent), in the settings top menu choose `Image packaging configuration` -> `Rootfs filesystem type` -> choose `EXT4 (SD/eMMC/SATA/USB)`. The SD card must have 2 partitions for this ([link](https://docs.xilinx.com/r/2021.2-English/ug1144-petalinux-tools-reference-guide/Partitioning-and-Formatting-an-SD-Card)).
5. `ESC` until it exits. Save changes when prompted.

Now the first basic changes are made, the first test image may be built. 


## Building and packaging the project

Building is simple, just issue

```bash {caption="Building the linux image"}
petalinux-build
```

 and wait for the completion. There should be no errors since it is built with default settings from a template. The build products (built linux image) will be located in `project_folder/images/linux`. 

Packaging options are different and are based on the `Image packaging configuration`. 

1. If no changes were made at configuration stage, the default option is `InitRAM`. There are several options on how to package it. 

   1. `petalinux-package --boot --fsbl --fpga --u-boot --force`
   It packages bootloaders and kernel+rootfs separately. 
   If one used this option, they need to copy `BOOT.BIN`, `boot.scr` and `image.ub` to the FAT32 SD card partition. 
   2.  `petalinux-package --boot --fsbl --fpga --u-boot --kernel --force`
   It packages everything into a single BOOT.bin image, including kernel and rootfs. It would be the only file one needs to copy to the FAT32 partition.

> _**Note**_: in this packaging mode the changes are not saved across reboots. It is useful if one needs a static image. 

2. If the chosen packaging mode is SD card, the packaging command is `petalinux-package --boot --fsbl --fpga --u-boot --force`. 
Files to copy: `BOOT.BIN`, `boot.scr`, `image.ub` to FAT32 partition. With **sudo file manager** (`sudo mc`, for example) copy the **contents** of the rootfs.tar.gz into the EXT4 partition. 

Insert the SD card into the slot and power up the board. Connect to the serial console using putty or other terminal. Make sure it boots, after that you may attach an ethernet cable to have SSH functionality. By default, the login/password is root/root. 

## Modifying the project

Generally, all the modifications that are necessary are stores in `project-spec` folder. `meta-user` folder contains recipes (Yocto) that need to be introduced into the final build. Here are the main folders in `meta-user` that are important. 

|   Folder   |       Description      |
|:----------:|:----------------------:|
|recipes-apps|Contains the recipes for applications. Create new app templates with `petalinux-create --type apps --name app_name --enable`. It would appear here for the user to edit.|
|recipes-bsp | Contains the recipes for device-tree (structure to connect the kernel drivers to the hardware), u-boot and FSBL|
|recipes-core| Contains the recipes for core system functionality. For example, if you want to change the default behaviour of the SSH server, it would be there to place the files| 
|recipes-kernel| Usually it stores the changes that the user has made through nemuconfig (`petalinux-config -c kernel`)|
|recipes-modules| It stores the recipes for kernel modules|

### FSBL patch

First stage bootloader is a standalone baremetal application provided by Xilinx that set up the processor for running and that handles the contents of BOOT.BIN file. By default, there are no debug prints (or any prints) from it, but sometimes it may be necessary to have those. Unfortunately, there are bugs in the menu config in petalinux that accept changes from the user for the build flags for FSBL, but they don't make it into the final image. 

To bypass that issue, there is an option to customize the build process for the proces. 

Navigate to (create if it doesn't exist) to `./project-spec/meta-user/recipes-bsp/embeddedsw/`. Create an empty folder called `files` and a file called `fsbl-firmware_%.bbappend`, if they don't exist. 

Add the following contents to the `fsbl-firmware_%.bbappend` file.

``` {caption="Adding debug functionality to FSBL"}
#Add debug for FSBL(optional)
XSCTH_BUILD_DEBUG = "1"
  
#Enable appropriate FSBL debug or compiler flags
YAML_COMPILER_FLAGS_append = " -DFSBL_DEBUG_INFO -DRSA_SUPPORT"
```

In this example, not only debug prints are added to the fsbl, but also RSA support, it is good to have this option at the early stage just not to come back to this if one needs to enable secure boot in the future. 

### Kernel module

There is a simplier way to have access to FPGA resources with generic-uio drivers that support interrupts, but this example is intended to give more details and flexibility to the reader to be able to create their own custom modules and projects, handle the interrupts their own way. The users would be able to have custom interrupt handlers to make more of their FPGA designs than just RW access.

To create a template for kernel module that would handle the interrupt from the FPGA, use

```bash {caption="New kernel module in petalinux project"}
petalinux-create -type modules --name module-name --enable
```
> _**Note**_: Don't use '_' in naming.

In this repository it is named `fpgatimer`.

Petalinux would create a template for a kernel module. 
Replace the contents of `./project-spec/meta-user/recipes-modules/fpgatimer/files/fpgatimer.c` with the contents from [**Appendix A**](#appendix-a) where the source code is presented. 

This module registers itself based on the interrupt number provided through device tree and relays the hardware interrupt to the userspace using OS signal SIGUSR1. 

It gives better performance compared to regular polling and frees the application from the wasteful constant polling. It may be very important in real-time (or performance sensitive) applications. 

### Device-tree

To connect the interrupt in the FPGA fabric to the kernel module, a change to the device tree is needed. 

Changes to the default configuration are made in the `./project-spec/meta-user/recipes-bsp/device-tree/files/system-user.dtsi`

> This file overwrites the device-tree nodes if they are present in the default device tree. Otherwise, it adds the missing configurations.

```C
/*Add this to the end of the file*/
/*Interrupt handler custom kernel module*/
&simple_hardware_timer_0 {
	compatible = "homeuser,fpgatimer";
	interrupt-parent = <&intc>;
	interrupts = <0 29 1>;
	};
```
The entry name is taken from Vivado design. It overwrites the default blob entry. 

> If unsure how the entry was names, after the first build the default generated device tree can be found at `./components/plnx_workspace/device-tree/device-tree/pl.dtsi`. This is an gutogenerated folder and no changes would be saved here.

Interrupt parent comes from the fact that the interrupt is connected to the main PS interrupt. Interrupts property reflects that this is not an SPI interrupt (0), interrupt number comes from Vivado design and 1 means rising edge triggering.

> IRQ number given by linux will be different from what was entered at this stage. 


### Custom application

In order to have your application embedded, there is an approach to compile sources as the part of the build process, but my personal preferrance is to embedd precompied and tested applications. 

In order to do that, create a new application with 

```bash
petalinux-create --type apps --name your_app_name --enable
```

Navigate to `./project-spec/meta-user/recipes-apps/your_app_name/files` and replace everything in that folder with the files from your precompiled application. 

One level above, modify your `your_app_name.bb` in the following fashion. 

```bash
#
# This file is the your_app_name recipe.
#

SUMMARY = "Short your_app_name description"
SECTION = "PETALINUX/apps"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

SRC_URI = "file://file1 \
           file://file2 \
           file://file3 \
           file://file4 \
           file://file5 \
          "

S = "${WORKDIR}"

FILES_${PN} += "/home/root/your_app_name_folder/*"

do_install() {
    install -d ${D}/home/root/your_app_name_folder
    cp ${S}/file1 ${D}/home/root/your_app_name_folder/file1
    cp ${S}/file2 ${D}/home/root/your_app_name_folder/file2
    cp ${S}/file3 ${D}/home/root/your_app_name_folder/file3
    cp ${S}/file4 ${D}/home/root/your_app_name_folder/file4
    cp ${S}/file5 ${D}/home/root/your_app_name_folder/file5
}
```

### Custom application auto-start 

> Auto-start is not featured in the repository's recipes.

If boot-time auto-start is necessary, a startup script is a nice-to-have.

```bash
petalinux-create --type apps --name startup --enable
```

Navigate to the newly created application's folder and edit the *.bb file similar to the following:

```bash
#
# This file is the startup recipe.
#

SUMMARY = "Startup script which starts from init.d and can be edited at ~/startup.sh"
SECTION = "PETALINUX/apps"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

SRC_URI = "file://startup.sh \
            file://startup_init \
            "

S = "${WORKDIR}"            

FILESEXTRAPATHS_prepend := "${THISDIR}/files:"

inherit update-rc.d

INITSCRIPT_NAME = "startup_init"
INITSCRIPT_PARAMS = "start 00 5 ."


do_install() {
        install -d ${D}${sysconfdir}/init.d/
        install -m 0755 ${S}/startup_init ${D}${sysconfdir}/init.d/startup_init
        install -d ${D}/home/root
        install -m 0755 startup.sh ${D}/home/root/startup.sh
}

FILES_${PN} += "${sysconfdir}/*"
FILES_${PN} += "/home/root/*"
```

The `startup_init` could be the following:

```bash
#!/bin/sh

if [ "$1" = "start" ]; then
    echo " "
    echo "Launching custom startup.sh"
    sh /home/root/startup.sh
fi

if [ "$1" = "stop" ]; then
    echo "Startup executing stop"
fi

exit 0
```

The `startup.sh` is the call for your application to launch or whatever is necessary to be done at the startup. 

## Real-time patch

To make the kernel give more throughput for computing and application's reactions, there is a real-time patch specific to the kernel version.

Guidelines for this approach were taken from [Hackster.io](https://www.hackster.io/LogicTronix/real-time-optimization-in-petalinux-with-rt-patch-on-mpsoc-5f4832). 

Link to the patch specific to 2021.2 kernel version: [5.10-rt16](https://mirrors.edge.kernel.org/pub/linux/kernel/projects/rt/5.10/older/patch-5.10-rc7-rt16.patch.gz)

Link to the patches storage: [Linux real-time patches](https://mirrors.edge.kernel.org/pub/linux/kernel/projects/rt/)

The necessary patch is also stored in this repository in `project-spec` folder. 

To patch the kernel call patch from the kernel sources directory.

```bash {caption="Patching kernel for real-time"}
cd ~/Cora_Z7/linux-xlnx-xilinx-v2021.2 # adjust to your location
zcat ../linux_build_cora/project-spec/kernel-5.10-realtime_patch.gz | patch -p1 # adjust to your location
```

### Kernel settings
After patching issue 

```bash
petalinux-config -c kernel
```

to adjust the settings to take advantage of the patch.

Settings to change:

1. In the menuconfig choose `General setup`.
2. Navigate to `Preemption model` and choose `Fully preemptible kernel (Real-time)`
3. Make sure `High resolution timer` is ticked in `Timers subsystem`. 
4. Navigate to `Kernel features` in the root of menuconfig, find `Timer frequency` and choose `1000 Hz`. 


### CPU isolation

> _**Note**_: these instructions are for Dual-core ARM processors. 

To isolate the application from kernel interruptions, some changes are necessary in default kernel boot arguments. To change those, issue 

```bash
petalinux-config
```

1. Navigate to `DTG Settings` -> `Kernel bootargs` and copy autogenerated bootargs (Highlight + ctrl + shift + c).
2. Turn off `Generate bootargs automatically`
3. Paste the copied bootargs and add `isolcpus=1 NOHZ_FULL`

This isolated the CPU #1 from the kernel. NOHZ_FULL makes the application non-preemtible, but needs to utilize SCHED_FIFO for this option to work. 

The software that this guide would cover as the next step utilizes these features. 


# Checking board's build (Launching)

Place files on the SD card and power up the board. 
Connect to the serial to see the bootlog. 

After booting, issue 

```bash
uname -a
```

The output should be similar to the following: `Linux cora 5.10.0-rt16-xilinx-v2021.2 #1 SMP PREEMPT_RT`

If that happened, the realtime patch was successfull.

Check the presense and the contents of isolated CPU entries: 

```bash
cat /sys/devices/system/cpu/isolated
```

The output should be 1 (In case isolcpus was set to 1)

Check the dmesg output for kernel module loading messages: 

```bash
dmesg | grep fpgatimer
```

There should be prints from the kernel module. 

Check the presence of the interrupts

```bash
cat /proc/interrupts
```

There should be an entry called fpgatimer. Non-zero interrupts values mean that the interrupts are handled by the module. 

Check that the timer module in FPGA is available manually.

```bash
devmem 0x43c00000 # current period threshold
devmem 0x43c00004 # if the bit for software polling was set
```

The values should be non-zero and there should be no bus error. 

> Bus error would mean that there was an attept to access a non-existent address

Check that the LED's are avaiable for blinking.

```bash
devmem 0x41210000 32 0x000000FF # turn the LEDs on
devmem 0x41210000 32 0x00000000 # turn the LEDs off
```

If everything is present, the base build is correct. The system can be used for software creation and\or other ways of utilization.

# Appendix A

```C {caption="Kernel module for hardware FPGA interrupts relay"}
#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/interrupt.h>
#include <linux/irq.h>
#include <linux/platform_device.h>
#include <linux/of.h>
#include <linux/of_device.h>
#include <linux/of_irq.h>
#include <linux/slab.h>
#include <linux/sched.h>
#include <linux/signal.h>
#include <linux/fs.h>
#include <linux/sched/signal.h>
#include <asm/uaccess.h>

MODULE_LICENSE("GPL");

// The soft OS signal we're going to use to relay the FPGA interrupt
#define SIG_NUM SIGUSR1

#define DRIVER_NAME "fpgatimer"

static struct of_device_id fpgatimer_driver_of_match[] = {
    { .compatible = "homeuser,fpgatimer", },
    {}
};
MODULE_DEVICE_TABLE(of, fpgatimer_driver_of_match);

// Init PID file with 0 
static int pid = 0;

// Interrupt handler fetches the contents of 'pid' attribute file
// for the PID where it needs to send the singal. If the PID is 0, it doesn't send anything.
// Debug prints are commented out not to spam dmesg. 
// Application needs to write 0 into the PID file before exiting. 
// If it was killed (PID != 0), don't spam dmesg. 
static irqreturn_t fpgatimer_isr(int irq, void *dev_id)
{
    struct task_struct *task;
    int ret;
    
    // Check if the PID is 0
    if (pid == 0) {
        //printk(KERN_INFO "No app, PID == 0, skipping interrupt\n");
        return IRQ_HANDLED;
    }
   
    /* Find the task associated with the PID */
    task = pid_task(find_vpid(pid), PIDTYPE_PID);
    if (!task) {
        //printk(KERN_ERR "Could not find the task with PID %d\n", pid);
        return IRQ_NONE;
    }

    /* Send the signal */
    ret = send_sig(SIG_NUM, task, 0);
    if (ret < 0) {
        printk(KERN_ERR "Error sending signal to application with PID %d\n", pid);
        return IRQ_NONE;
    }

    return IRQ_HANDLED;
}


// Function to show the contents of the attribute file
static ssize_t pid_show(struct kobject *kobj, struct kobj_attribute *attr,
                      char *buf)
{
    return sprintf(buf, "%d\n", pid);
}

// Function to enable write (store) into the attribute file
static ssize_t pid_store(struct kobject *kobj, struct kobj_attribute *attr,
                       const char *buf, size_t count)
{
    int ret;

    ret = kstrtoint(buf, 10, &pid);
    if (ret < 0)
        return ret;

    return count;
}

static struct kobj_attribute pid_attribute =
    __ATTR(pid, 0664, pid_show, pid_store);

// Find the actual IRQ number and register the interrupt
static int fpgatimer_driver_probe(struct platform_device* dev)
{
    printk(KERN_INFO "fpgatimer: probing driver...\n");

    unsigned int irq;
    irq = irq_of_parse_and_map(dev->dev.of_node, 0);
    printk(KERN_INFO "fpgatimer: found matching irq = %d\n", irq);
    if (request_irq(irq, fpgatimer_isr, 0, DRIVER_NAME, &dev->dev))
        return -1;
    printk(KERN_INFO "fpgatimer: registered irq\n");
    
    return 0;
}

// Unregister the interrupt upon removal
static int fpgatimer_driver_remove(struct platform_device* dev)
{
    printk(KERN_INFO "fpgatimer: removing driver...\n");

    free_irq(of_irq_get(dev->dev.of_node, 0), &dev->dev);

    return 0;
}

// Driver's struct
static struct platform_driver fpgatimer_driver = {
    .probe = fpgatimer_driver_probe,
    .remove = fpgatimer_driver_remove,
    .driver = {
        .name = DRIVER_NAME,
        .owner = THIS_MODULE,
        .of_match_table = fpgatimer_driver_of_match,
    },
};


// kobject
static struct kobject *fpgatimer_kobj;

// Module init function
static int __init fpgatimer_init(void)
{
    printk(KERN_INFO "fpgatimer: init...\n");

    int retval;
    
    // Create a kobject and add it to the sysfs
    fpgatimer_kobj = kobject_create_and_add(DRIVER_NAME, kernel_kobj);
    if (!fpgatimer_kobj) {
        printk(KERN_WARNING "failed to create kobject\n");
        return -ENOMEM;
    }

    // Create the 'pid' attribute file
    retval = sysfs_create_file(fpgatimer_kobj, &pid_attribute.attr);
    if (retval) {
        printk(KERN_WARNING "failed to create sysfs file\n");
        kobject_put(fpgatimer_kobj);
        return retval;
    }
    
    // register platform driver
    if (platform_driver_register(&fpgatimer_driver)) {                                                     
        printk(KERN_WARNING "failed to register platform driver \"%s\"\n", DRIVER_NAME);
        return -1;                                                
    }
    printk(KERN_INFO "fpgatimer: registered platform driver\n");

    return 0;
}

// Stop routine
static void __exit fpgatimer_exit(void)
{
    printk(KERN_INFO "fpgatimer: stopped\n");
    platform_driver_unregister(&fpgatimer_driver);
    sysfs_remove_file(fpgatimer_kobj, &pid_attribute.attr);
    kobject_put(fpgatimer_kobj);
}

module_init(fpgatimer_init);
module_exit(fpgatimer_exit);

MODULE_AUTHOR ("Dmitrii Matafonov");
MODULE_DESCRIPTION("FPGA hard interrupt to userspace relay");
MODULE_LICENSE("GPL v2");
MODULE_ALIAS("custom:fpga-timer");
```