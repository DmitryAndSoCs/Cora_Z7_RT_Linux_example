# How to build the sources from Vivado

1. Open Vivado 2021.2
2. In the TCL console on the bottom:

```bash
cd ./your/folder/vivado_hdl/cora_hw_base
source ./cora_hw_base.tcl
```
> Note: because of Vivado's recreation flow, there would be a new folder named cora_hw_base where you source the script. It would contain cora_hw_base.srcs as well. If you want to keep your changes in git:
> 1. create a new TCL script (File -> Project -> Write Tcl...)
> 2. Copy the new contents of the .srcs folder one level above (where it was ogirinally tracked) and the new TCL script. 
> 3. Commit.


It would recreate the project from the sources. 
To build everything, click `Generate bitstream` in Vivado's Flow Navigator.
