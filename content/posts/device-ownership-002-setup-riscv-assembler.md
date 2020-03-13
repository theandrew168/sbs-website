---
date: 2020-02-23
title: "Device Ownership 002: Setup - RISC-V Assembler"
slug: "device-ownership-002"
tags: ["risc-v", "firmware", "device ownership", "setup"]
---

# Ready, set, assemble!
This post details the steps necessary to "assemble" our simple program.
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
