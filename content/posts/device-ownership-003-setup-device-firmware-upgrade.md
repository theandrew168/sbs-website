---
date: 2020-02-24
title: "Device Ownership 003: Setup - Device Firmware Upgrade"
slug: "device-ownership-003"
tags: ["risc-v", "firmware", "device ownership", "setup"]
---

# Preparing for DFU
Without our assembled `smallest.bin` file in hand, the next step is to hook up our Longan Nano and upload the program to the chip's flash storage.
To accomplish this task, we must use the [Device Firmware Upgrade](https://en.wikipedia.org/wiki/USB#Device_Firmware_Upgrade) protocol.
This protocol enables a very simple method for upgrading the firmware of devices connected to your system over USB.
If you are curious about the details, the official specification for DFU can be found [here](https://www.usb.org/sites/default/files/DFU_1.1.pdf).

One important note about DFU is that the terms "upload" and "download" are expressed from the device's perspective.
So, "downloading" firmware means writing an assembled, binary file from your system to the device.
This is what we'll need to do in order to get our `smallest.bin` program onto the Longan Nano.
Conversely, "uploading" firmware with DFU means reading all of the data on the device's flash storage and saving it a file on your local system.

To program our device, we will be using a program called [dfu-util](git clone git://git.code.sf.net/p/dfu-util/dfu-util
).
Even though this program is available through the `apt` package system, we need to manually build a newer version that includes fixes for [multiple](https://sourceforge.net/p/dfu-util/dfu-util/ci/529fa5147613218c75dfa441c64df9b28910fe1c/) [bugs](https://sourceforge.net/p/dfu-util/dfu-util/ci/f2b7d4b1113ef6c3ada31a0654c9aefebcdb1de5/) in the GD32VF103 CPU's implementation of DFU.

Before building dfu-util, we need to install its dependencies.
```
sudo apt install autoconf build-essential git libusb-1.0-0-dev pkg-config
```

Now we can clone the source code, build the program, and install it.
The following steps include deleting the source code directory because we won't need it once dfu-util is installed.
```
git clone git://git.code.sf.net/p/dfu-util/dfu-util
cd dfu-util
./autogen.sh
./configure
make
sudo make install
cd ..
rm -r dfu-util/
```

To verify that the install was successful, try out the command:
```
dfu-util --version
```
