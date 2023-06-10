# Cora_Z7_10 project

Project is going to use Digilent's Cora Z7-10. (which is no longer produced). 
The process is the same for more common Z7-7S version but it has only 1 ARM core. CPU isolation approach would not be acceptable for that. 

This project is aimed at creating a streamlined example of numerous techniques that would let the users create complex designs based on Xilinx SoC chips. It is intended to use a developement board, but the process for the custom board is not going to be much different. 

Key features: 
- Linux Real-time patch, CPU1 is isolated for the user application to run there
- Kernel module that fetches the interrupts from FPGA IP and passes it to a userspace app as a signal SIGUSR1
- Golang application takes advantage of the kernel features.
    > It blinks LEDs, but the approch would be the same for more complex applications.


All the necessary changes for custom boards are necessary at HDL stage and we'll try to cover those where it might be needed. 

- Check [Releases](https://github.com/DmitryAndSoCs/Cora_Z7_RT_Linux_example/releases) for the Cora Z7-10 loadable BOOT.BIN. Blink application is `~/blinkapp/blink.run` (on the dev board after boot).

    - To launch is use 

    ```bash 
    ~/blinkapp/blink.run # OR
    ~/blinkapp/blink.run --timer 39999999 # adjust the cycles number addording to the formula: period = 10 ns * (cycles +1)
    ```

The build instructions that contain the description of the techniques and step-by-step guides:

- Vivado project guide [[PDF]](https://github.com/DmitryAndSoCs/Cora_Z7_RT_Linux_example/blob/v0.0.2/docs/Cora_Z7_10_Vivado_project_guide.pdf), [[Markdown]](https://github.com/DmitryAndSoCs/Cora_Z7_RT_Linux_example/blob/v0.0.2/docs/src/vivado_project_guide.md)
- Petalinux project guide [[PDF]](https://github.com/DmitryAndSoCs/Cora_Z7_RT_Linux_example/blob/v0.0.2/docs/Cora_Z7_10_Petalinux_project_guide.pdf), [[Markdown]](https://github.com/DmitryAndSoCs/Cora_Z7_RT_Linux_example/blob/v0.0.2/docs/src/petalinux_project_guide.md)
- Golang application description is not developed yet


The project is under development. Updates are happening when I have time. 

# Tools

This project uses Vivado/Vitis 2021.2 and corresponding petalinux version (2021.2). 
For some endgoal applications that would be written in Golang this project is going to use Go (1.20) and VSCode **go** extention. 

# Instructions

The detailed instructions are going to be in the `docs` folder. They are written in markdown and rendered as PDFs with pandoc.

# Links

- Vivado downloads: [Version archive](https://www.xilinx.com/support/download/index.html/content/xilinx/en/downloadNav/vivado-design-tools/archive.html)
- Petalinux downloads: [Version archive](https://www.xilinx.com/support/download/index.html/content/xilinx/en/downloadNav/embedded-design-tools/archive.html)
- Petalinux required packages: [AR73296](https://support.xilinx.com/s/article/73296?language=en_US)
- May be handy. [Petalinux 2021.2 docker](https://github.com/carlesfernandez/docker-petalinux2)
- Golang in VSCode: [VSCode guide](https://learn.microsoft.com/en-us/azure/developer/go/configure-visual-studio-code)

## Tools

I would suggest sticking to a not-so-recent versions as there are many examples and forum discussions for earlier versions.
Sometimes there are problems between versions bcause something has changed somewhere and you may end up figuring out what the differences are between versions of Xilinx tools and how to fix those. (Like, there are A LOT of small but confusing changes between 2019.1 - 2019.2 (Vitis), 2019.2 - 2020.1 (petalinux) and so on)
[Xilinx Vitis/Vivado 2021.2](https://www.xilinx.com/support/download/index.html/content/xilinx/en/downloadNav/vivado-design-tools/archive.html)
[Petalinux tools 2021.2](https://www.xilinx.com/support/download/index.html/content/xilinx/en/downloadNav/embedded-design-tools/archive.html)
[Golang in VSCode quickstart](https://learn.microsoft.com/en-us/azure/developer/go/configure-visual-studio-code)
[Reference designs for Cora Z7](https://digilent.com/reference/programmable-logic/cora-z7/start) <- not needed but might be useful

> You will need a Xilinx (AMD) with a form filled to be able to download Xilinx tools. 


