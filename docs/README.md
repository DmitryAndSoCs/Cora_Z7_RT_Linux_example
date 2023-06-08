This documentation was built with [pandoc](https://pandoc.org/installing.html)

Original template is [Eisvogel](https://github.com/Wandmalfarbe/pandoc-latex-template)

The mod file was a combination of what I did myself to add VHDL highlighting in listings and some other similar template repositories on GitHub, I didn't write the names of those, unfortunately. If you find that I miss some mentions, please, let me know. 

The mod has to be installed in pandoc templates (depends on the system)

Command to build is 

```bash
cd ./src

pandoc vivado_project_guide.md -o ../Cora_Z7_10_Vivado_project_guide.pdf --from markdown --template eisvogel_mod --toc --listings

pandoc petalinux_project_guide.md -o ../Cora_Z7_10_Petalinux_project_guide.pdf --from markdown --template eisvogel_mod --toc --listings
```