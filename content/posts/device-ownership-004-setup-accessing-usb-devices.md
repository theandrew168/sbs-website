---
date: 2020-02-25
title: "Device Ownership 004: Setup - Accessing USB Devices"
slug: "device-ownership-004"
tags: ["risc-v", "firmware", "device ownership", "setup"]
---

# Identifying the device
The ability to interact with most USB devices on Linux systems is restricted to the root user.
However, the [udev](https://en.wikipedia.org/wiki/Udev) device manager exists to enable more granular access to specific devices for non-root users.
In our case, we want to be able to perform read and write operations on the Longan Nano.

Udev looks for special "rules" files in multiple system directories.
The one we care about is `/etc/udev/rules.d/`.
These rules files contain lines of key-value pairs that can be used to filter specific devices and change the permissions associated with them.
To identity a specific USB device, we need to know two pieces of information: its vendor ID and its product ID.
Every unique USB device can be identified by this pair of IDs.
Once we know them, we can setup a rule that effectively says: "if you see the Longan Nano, allow non-root users to interact with it".

What would be a good way to find these IDs?
Since some basic USB device information _is_ available to non-root users by default, we can use a USB listing utility such as `lsusb`.
Let's plug in the Longan Nano and see what we can see!
After connecting a USB-C cable from your system to the device, try running `lsusb` and looking for something related to GigaDevice or DFU.

```
~/device_ownership$ lsusb
Bus 004 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
Bus 003 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
Bus 002 Device 002: ID 0bda:0316 Realtek Semiconductor Corp. USB3.0-CRW
Bus 002 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
Bus 001 Device 003: ID 13d3:56a6 IMC Networks Integrated Camera
Bus 001 Device 002: ID 8087:0a2b Intel Corp. 
Bus 001 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
```

Unfortunately, nothing obvious shows up!
This is actually expected and is related to how DFU works.
See, when in DFU mode, a USB device acts like (and looks like) something completely different.
By default, the Longan Nano does NOT enter DFU mode.
Instead, it simply runs whatever program is currently loaded on its flash storage which may or may not involve interacting with the host sytem over USB.
Unless explicitly programmed to do so, the Longan Nano won't even appear in the output of programs like `lsusb`.
In order to "see" the device so that we can reprogram it, we need to somehow force it into DFU mode when it powers on.

Different DFU-capcable chips have different ways of entering DFU mode.
In the case of the Longan Nano, the magic lies within its two buttons: RESET and BOOT.
The first button, RESET, does what you'd expect it to: it resets the device.
This is essentially the same as telling a modern computer to restart.
It resets the chip's power connection which clears all of the CPU's internal state.

The BOOT button, on the other hand, is a special button that tells the CPU to boot into DFU mode.
This is exactly what we're after!
In order to be sure that the CPU sees the signal from the BOOT button, it must be held down the while the chip is powered on or reset via the RESET button.
Since the buttons are so small, I find it easiest to hold the BOOT button down, then press and release the RESET button.
In short: press BOOT, press RESET, release RESET, release BOOT.

**Take note of this BOOT / RESET process! It will be used everytime we to download a new program to the device!**

Now that the Longan Nano is in DFU mode, it should be visible to us via the `lsusb` command.
Let's try it again:

```
~/device_ownership$ lsusb
Bus 004 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
Bus 003 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
Bus 002 Device 002: ID 0bda:0316 Realtek Semiconductor Corp. USB3.0-CRW
Bus 002 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
Bus 001 Device 003: ID 13d3:56a6 IMC Networks Integrated Camera
Bus 001 Device 002: ID 8087:0a2b Intel Corp. 
Bus 001 Device 004: ID 28e9:0189 GDMicroelectronics GD32 0x418 DFU Bootloade
Bus 001 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
```

There it is!
The device with the description `GDMicroelectronics GD32 0x418 DFU Bootloade` is definitely the Longan Nano.
The multiple instances of "GD" are short for "GigaDevice", the CPU's manufacturer.
To further increase our confidence, we can see most of the phrase `DFU Bootloader` in the device's description.
Now that we know which line corresponds to the Longan Nano, we can find the two IDs we are after in the string `ID 28e9:0189`.
The two IDs are separated by a colon.
The first one is the vendor ID and the second is the product ID.
For the Longan Nano device, its vendor ID is `23e9` and its product ID is `0189`.

# Accessing the device
With all of that digging out of the way, we can finally write the udev rule that we need to allow us to interact with the device.
The syntax used by udev rules is fairly straighforward.
Each comma-separated, key-value pair is either a filter (denoted by a comparison operator such as `==`) or an assignment (denoted by the assignment operator `=`).
In pseudo-code, this is what we want our udev rule to say.
```
if device.vendor_id == "28e9" and device.product_id == "0189":
    device.mode = "0666"  # this means "everyone can read and write"
```

In [udev syntax](https://linux.die.net/man/7/udev), our rule looks like this:
```
ATTRS{idVendor}=="28e9", ATTRS{idProduct}=="0189", MODE="0666"
```

As mentioned earlier, udev rules are added by means of special rules files.
We will name our file at `99-longan-nano.rules` and place it under `/etc/udev/rules.d/`.
Below is a command to create the rules file with the single rule that we need and reload udev.
```
cat <<EOF | sudo tee /etc/udev/rules.d/99-longan-nano.rules
ATTRS{idVendor}=="28e9", ATTRS{idProduct}=="0189", MODE="0666"
EOF
```

We then need to reload udev in order to pickup the new rule.
```
sudo udevadm control --reload
```

That should do it!
The last step now is to unplug and reinsert the USB cable from your host machine.
Udev applies its rules as devices are detected so a re-plug is necessary to pickup the new permissions.

With all the correct permissions in place and the device in DFU mode, we should be able to see something using dfu-util.
```
~/device_ownership$ dfu-util --list
Found DFU: [28e9:0189] ver=1000, devnum=7, cfg=1, intf=0, path="1-2", alt=1, name="@Option Bytes  /0x1FFFF800/01*016 g", serial="3CBJ"
Found DFU: [28e9:0189] ver=1000, devnum=7, cfg=1, intf=0, path="1-2", alt=0, name="@Internal Flash  /0x08000000/512*002Kg", serial="3CBJ"
```

Very cool!
This output actually shows two DFU alternatives: one for "Option Bytes" and one for "Internal Flash".
Since the GD32VF103 CPU executes code from it's internal flash storage, that second option is what we're after.
To be more specific, once the CPU is initialized, it starts executing code from address 0x08000000 in its flash.
