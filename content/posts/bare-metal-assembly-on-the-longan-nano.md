---
date: 2020-04-24
title: "Bare-Metal Assembly on the Longan Nano"
slug: "bare-metal-assembly-on-the-longan-nano"
tags: ["risc-v", "assembly"]
---
The [Sipeed Longan Nano](https://www.seeedstudio.com/Sipeed-Longan-Nano-RISC-V-GD32VF103CBT6-Development-Board-p-4205.html) is a small, affordable, 32-bit [RISC-V](https://en.wikipedia.org/wiki/RISC-V) chip.
Despite its minimalism, the Longan Nano provides enough power and peripherals to learn RISC-V assembly and build exciting programs along the way.

The Longan Nano comes with the following components:

* GigaDevice GD32VF103[CBT6] 32-bit CPU
* 8 MHz default clock speed (IRC8M)
* 108 MHz maximum clock speed
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

Since the [assembler](https://en.wikipedia.org/wiki/Assembly_language#Assembler) we are going to use is written in [Python](https://www.python.org), we will install all of the system packages necessary to install additional Python modules.
```
sudo apt install python3-pip python3-venv
```

To keep things organized, we will use a [virtual environment](https://docs.python.org/3/library/venv.html) to isolate our project-specific Python modules from the rest of our system.
```
python3 -m venv venv/
. ./venv/bin/activate
```
**Note that this last command above will be required any time you come back to work on the project!**

With a virtual environment in place, we can install the simple RISC-V assembler built for this project: [bronzebeard](https://pypi.org/project/bronzebeard/).
```
pip install bronzebeard
```

Last but not least, we can write our simple program to a file and assemble it. 
```
echo "addi zero, zero, 0" > smallest.asm
python3 -m bronzebeard.asm smallest.asm smallest.bin
```

Congratulations!
You just wrote and assembled the world's smallest RISC-V program!

# Preparing for DFU
Without our assembled `smallest.bin` file in hand, the next step is to hook up our Longan Nano and upload the program to the chip's flash storage.
To accomplish this task, we must use the [Device Firmware Upgrade](https://en.wikipedia.org/wiki/USB#Device_Firmware_Upgrade) protocol.
This protocol enables a very simple method for upgrading the firmware of devices connected to your system over USB.
If you are curious about the details, the official specification for DFU can be found [here](https://www.usb.org/sites/default/files/DFU_1.1.pdf).
Thankfully, the Bronzebeard project includes a minimal DFU uploader!
Consult the [setup instructions](https://github.com/theandrew168/bronzebeard#setup) in Bronzebeard's README for details on how to setup the proper USB tools.

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

# Programming the device
We now have all the information we need to upload our `smallest.bin` program to the device.
One thing that the DFU uploader needs is the USB identifier for the device.
If you are on a Linux or macOS system, the command `lsusb` can be helpful for finding this info.

In our case, we are programming the Longan Nano so the device ID is already known:
```
python3 -m bronzebeard.dfu 28e9:0189 smallest.bin
```

To take the Longan Nano out of DFU mode and back to normal, simply press the RESET button.
Get ready for the excitement!
Just kidding.
This program won't do anything.
Nothing will light up, nothing will flash, nothing will beep.
Do not discredit this achievement, though!

# All that work for nothing?
Who knew that doing absolutely nothing could be so much work?
Most of what we covered wasn't even related to RISC-V or the code: it was just setup work and preparation.
Don't let that worry you, though.
Executing bare-metal code like this is the first step in achieving great things.
For some next steps, try turning on some LEDs, reading an SD card, or even controlling the builtin LCD!
