---
date: 2020-02-19
title: "Device Ownership 002: Longan Nano - Controlling LEDs"
slug: "device-ownership-002"
tags: ["risc-v", "firmware", "longan nano"]
draft: true
---
# Intro
where we left off: we can write riscv assembly, upload it to the nano, and run it  
but we couldn't see any output. how can we do that?  
the board has 3 LEDs... lets light em up!  

# Datasheets and Schematics
where do we get em?  
where are the leds on the schematic?  
what is GPIO?  

# GPIO Basics
talk about it  
write the code  
but its not working! what did we miss?  

# RCU Basics
ahh... the GPIO isn't "on" (it has no clock signal running to it)  
how do we enable that?  
hit the books  
more code!  
get it working  

# Conclusion
there you have it!  
consider what we've done: came from nothing and made this little computer do stuff  
this is only the beginning  
next we will control the LCD  
maybe even a bit of graphics!!!  
