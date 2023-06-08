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