---
title: "Device Ownership: Longan Nano - Getting Started"
date: 2020-02-27
tags: ["risc-v", "firmware", "longan nano"]
draft: true
---
# Introduction
The [Sipeed Longan Nano](https://www.seeedstudio.com/Sipeed-Longan-Nano-RISC-V-GD32VF103CBT6-Development-Board-p-4205.html) is a small, affordable, 32-bit RISC-V chip.

Specs:

* GD32VF103[CBT6] 32-bit CPU
* 128 KB flash storage
* 32 KB sram memory
* 3 LEDs (red, green, and blue)
* 1 USB Type-C port
* 160x80 pixel LCD (0.96 inches)
* 2 Buttons (RESET and BOOT)

# The world's smallest RISC-V program
Seen below is a valid RISC-V program that does absolutely nothing!
```
addi zero, zero, 0
```
This particular combination of words and numbers tells the computer to add `0` to a special register named `zero`.
More specifically, `addi` (shorthand for "add immediate") is an instruction that tells the CPU to add the literal number `0` to the value currently stored in register `zero` and then store the resulting sum back into register `zero`.

What is a register, though?

A helpful analogy for explaining the relationship between memory, a CPU, and its registers is to compare them to a library, a librarian, and the books on their desk.
The librarian's job is to repair, sort, and organize books.
Even though the library contains thousands and thousands of books, they can only fit 32 on their desk at any given time.
If the librarian needs to work on a book that is already on their desk, then all is well!
They can instantly pick it up and get to work.
If the book is _not_ on their desk, however, then they need to get up, walk all the way over to correct shelf, grab the book, and bring it back.
If their desk was already full, then an old book must be stored back on a shelf or thrown away.

When it comes to computers, we deal with numbers instead of books.
All RISC-V CPUs come with 32 general-purpose registers for holding temporary values to be used by the program.
In the case of the Longan Nano, these registers are 32 bits long meaning that they can store any integer value from `-2,147,483,648` up to `2,147,483,647`.

# Setup
note that this only applies to ubuntu/apt systems for now  
split this big boi up by section  
```
# install required system libraries
sudo apt update
sudo apt install build-essential git libusb-1.0-0-dev
sudo apt install python3-pip python3-venv

# build dfu-util from source and install
git clone git://git.code.sf.net/p/dfu-util/dfu-util
cd dfu-util
./autogen.sh
./configure
make -j4
sudo make install
cd ..
rm -r dfu-util/

# create project directory
mkdir device_ownership
cd device_ownership/

# create a virtualenv and install simpleriscv
python -m venv venv
source venv/bin/activate
pip install simpleriscv
```

# Build
clean this up  
explain udev to avoid the sudo? :'(  
explain how upload / download is nutters in DFU  
link that DFU spec like a boss  

```
simpleriscv asm smallest.asm
udev magic
dfu-util -l
dfu-util -a 0 -D smallest.bin -s 0x08000000:0x00020000
dfu-util -a 0 -U firmware.bin -s 0x08000000:0x00020000
```

# Conclusion
did it even work? did it do anything? we can't really know  
in the next post, we will write a program that does something we can see!  
