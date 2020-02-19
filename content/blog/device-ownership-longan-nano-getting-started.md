---
title: "Device Ownership: Longan Nano - Getting Started"
date: 2020-02-19
tags: []
draft: true
---
### Intro
Longan Nano is a RISCV chip that uses the GD32VF103[CBT6]  
128KB flash  
32KB sram  
3 leds (RGB)  
USB-C w/ DFU  

### Goals
Make the device do something from scratch  
create something from nothing  
why not use the provided C-based SDK? because you don't learn anything!  
goal of this post is get custom code running on the device  

### Setup
note that this only applies to ubuntu/apt systems for now  
install python3  
install simpleriscv  
build and install dfu-util (explain the need)  

### Program
one-line file with a nop?  
addi zero, zero, 0  
compile  
upload  

### Conclusion
did it even work? did it do anything? we can't really know  
in the next post, we will write a program that does something we can see!  
