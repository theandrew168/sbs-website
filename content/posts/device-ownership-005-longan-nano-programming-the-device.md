---
date: 2020-02-26
title: "Device Ownership 005: Longan Nano - Programming the Device"
slug: "device-ownership-005"
tags: ["risc-v", "firmware", "device ownership", "longan nano"]
---

# Programming the device
We now have all the information we need to download our `smallest.bin` program on the device.
Remember from before that dfu-util shows us two alternative settings for upgrading firmware.
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
Stay tuned from some real action!
