# Cora_Z7_10 project

Project is going to use Digilent's Cora Z7-10. (which is no longer produced). 
The process is the same for more common Z7-7S version but it has only 1 ARM core. CPU isolation approach would not be acceptable for that. 

This project is aimed to create a streamlined example of numerous techniques
that would let the users create complex designs based on Xilinx SoC chips. It is intended to use a developement board, but the process for the custom board is not going to be much different. 
All the necessary changes for custom boards are made at HDL stage and we'll try to cover those where it might be needed. 

The project is under development. Updates are happening when I have time. 

# Tools

This project uses Vivado/Vitis 2021.2 and corresponding petalinux version (2021.2). 
For some endgoal applications that would be written in Golang this project is going to use Go (1.20) and VSCode **go** extention. 

# Instructions

The detailed instructions are going to be in the `docs` folder. 

# Links

## Tools

I would suggest sticking to a not-so-recent versions as there are many examples and forum discussions for earlier versions.
Sometimes there are problems between versions bcause something has changed somewhere and you may end up figuring out what the differences are between versions of Xilinx tools and how to fix those. (Like, there are A LOT of small but confusing changes between 2019.1 - 2019.2 (Vitis), 2019.2 - 2020.1 (petalinux) and so on)
[Xilinx Vitis/Vivado 2021.2](https://www.xilinx.com/support/download/index.html/content/xilinx/en/downloadNav/vivado-design-tools/archive.html)
[Petalinux tools 2021.2](https://www.xilinx.com/support/download/index.html/content/xilinx/en/downloadNav/embedded-design-tools/archive.html)
[Golang in VSCode quickstart](https://learn.microsoft.com/en-us/azure/developer/go/configure-visual-studio-code)
[Reference designs for Cora Z7](https://digilent.com/reference/programmable-logic/cora-z7/start) <- not needed but might be useful

> You will need a Xilinx (AMD) with a form filled to be able to download Xilinx tools. 


