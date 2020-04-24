---
date: 2020-03-07
title: "Device Ownership 001: Longan Nano - Getting Started"
slug: "device-ownership-001"
tags: ["risc-v", "firmware", "longan nano"]
---
# Introduction
The [Sipeed Longan Nano](https://www.seeedstudio.com/Sipeed-Longan-Nano-RISC-V-GD32VF103CBT6-Development-Board-p-4205.html) is a small, affordable, 32-bit [RISC-V](https://en.wikipedia.org/wiki/RISC-V) chip.
Despite its minimalism, the Longan Nano provides enough power and peripherals to learn RISC-V assembly and build exciting programs along the way.

The Longan Nano comes with the following components:

* 108 MHz GigaDevice GD32VF103[CBT6] 32-bit CPU
* 128 KB flash storage
* 32 KB sram memory
* 3 LEDs (red, green, and blue)
* 1 USB Type-C port
* 1 microSD card slot
* 160x80 pixel LCD (0.96 inches)
* 2 Buttons (RESET and BOOT)

With the specs out of the way, we can begin the journey of figuring out how to go from a blank text file to controlling all aspects of our hardware.
Get ready to discuss registers, design clean programs, and dive into some datasheets.

**The rest of the Longan Nano series requires that the device be connected to your host system via a USB-C cable! Go ahead get things wired up!**

# The world's smallest RISC-V program
Seen below is a valid RISC-V program that does absolutely nothing!
```
addi zero, zero, 0
```
This particular combination of words and numbers tells the computer to add 0 to a special register named `zero`.
More specifically, `addi` (shorthand for "add immediate") is an instruction that tells the CPU to add the literal number 0 to the value currently stored in register `zero` and then store the resulting sum back into register `zero`.

What is a register, though?

One analogy for explaining the relationship between a CPU and its registers is to consider a librarian.
The librarian's job is to organize and repair books.
Within arm's reach of their desk, the librarian has a small book shelf capable of holding only 32 books.
If the librarian needs to work on a book that is already in their shelf, then they can instantly pick it up and get to work.
If the book is _not_ on their shelf, however, then they need to stop working, walk deeper into the library, grab the book, and bring it back.
In order to make room for this new book, an old book must be taken back and stored somewhere in the library.

When it comes to computers, we deal with numbers instead of books.
All RISC-V CPUs come with 32 general-purpose registers for holding temporary values to be used by the program.
In the case of the Longan Nano, these registers are 32 bits long meaning that they can store any integer value from -2,147,483,648 up to 2,147,483,647.

All 32 of RISC-V's registers behave the same except for one: the `zero` register.
This register is special because it always holds the integer value 0.
If any other value attempts to get placed there, it effectively disappears.
This register is useful for discarding unwanted values and for when 0 is used in computation.

Knowing what we now know, we can look at our simple program in a bit more detail.
```
addi zero, zero, 0
```
The `addi` instruction takes three parameters: the destination register, the source register, and an immediate integer value.
These parameters are arranged as follows:
```
addi rd, rs1, imm
```

To recap, our simple program is adding the value currently stored in the `zero` register (which is always 0) to the value immediate integer value 0.
This sum (0 + 0 = 0) is then stored back into register `zero`.
Even though the program appears to write 0 to the `zero` register, it doesn't affect it at all.
As we learned with the special `zero` register, all values placed in it disappear.
Even if the result of the sum had been non-zero, the instruction would not have changed the `zero` register.

Now that our simple program is written and understood, how do we convert it to something that CPUs understand?

# Ready, set, assemble!
This section details the steps necessary to "assemble" our simple program.
Assembling is the process of converting human-readable assembly language text to a specific binary representation that a CPU can understand.

To start, we create a directory to hold the our project's development files.
In an open terminal window:
```
mkdir device_ownership
cd device_ownership/
```

Then, since the [assembler](https://en.wikipedia.org/wiki/Assembly_language#Assembler) we are going to use is written in [Python](https://www.python.org), we will install all of the system packages necessary to install additional Python modules.
```
sudo apt install python3-pip python3-venv
```

To keep things organized, we will use a [virtual environment](https://docs.python.org/3/library/venv.html) to isolate our project-specific Python modules from the rest of our system.
```
python3 -m venv venv/
source venv/bin/activate
```
**Note that that `source` command above will be required any time you come back to work on the project!**

With a virtual environment in place, we can install the RISC-V assembler used in this series: [simpleriscv](https://pypi.org/project/simpleriscv/).
```
pip install simpleriscv
```

Last but not least, we can write our simple program to a file and assemble it. 
```
echo "addi zero, zero, 0" > smallest.asm
simpleriscv asm smallest.asm
```

By default, simpleriscv places its output in a file named with a `.bin` extension in place of the original `.asm`.

Congratulations!
You just wrote and assembled the world's smallest RISC-V program!

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
Since the buttons are so small, I find it easiest to hold the BOOT button down, then press the RESET button at the same time.
After both buttons are pressed, both buttons can be released at the same time.

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

# Programming the device
We now have all the information we need to download our `smallest.bin` program on the device.
Remember from above that dfu-util shows us two alternative settings for upgrading firmware.
Since we want to program the chip's internal flash, we will need to specify the correct "alt" identifier (which is 0 in this case).
Lastly, we need to specify the "range" of the Longan Nano's flash storage: where it begins in memory and how large it is.
We know that it starts at address 0x08000000 and is 128 kilobytes in size.
Converting this size to bytes gives us 131072 bytes.
In hexadecimal, this number is represented as 0x20000.

Putting everything together in a way that dfu-util understands yields the following command.
```
dfu-util --download smallest.bin --alt 0 --dfuse-address 0x08000000:0x20000
```

To take the Longan Nano out of DFU mode and back to normal, simply press the RESET button.
Get ready for the excitement!
Just kidding.
This program won't do anything.
Nothing will light up, nothing will flash, nothing will beep.
Do not discredit this achievement, though!

# All that work for nothing?!
Who knew that doing absolutely nothing could be so much work?
Most of what we covered wasn't even related to RISC-V or the code: it was just setup work and preparation.
Don't let that worry you, though.
Almost all of the explanations in this post (DFU mode, udev rules, using dfu-util) won't have to be covered again.

In the next post we will do something slight more flashy: turning on an LED!
Moving forward from there, the future is even brighter.
We will look into using the device's timers to flash LEDs on and off.
We will even get into using that big LCD screen on the front.
Stick around for more of this series and some real action.

Thanks for sticking around to the end!
