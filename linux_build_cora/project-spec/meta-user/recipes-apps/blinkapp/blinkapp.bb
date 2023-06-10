#
# This file is the blinkapp recipe.
#

SUMMARY = "Simple blinkapp application"
SECTION = "PETALINUX/apps"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

SRC_URI = "file://blink.run"

S = "${WORKDIR}"

FILES_${PN} += "/home/root/blinkapp/*"

do_install() {
	install -d ${D}/home/root/blinkapp
	cp ${S}/blink.run ${D}/home/root/blinkapp/blink.run
}
