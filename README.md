# Reflow Oven

This is a basic reflow oven for SMT reflow soldering. A real-time temperature graph can be viewed in any web browser.

![reflow oven](https://minor-industries.sfo2.digitaloceanspaces.com/hw/reflow-oven.jpg)

![reflow profile](https://minor-industries.sfo2.digitaloceanspaces.com/sw/reflow-oven-with-profile.png)

I've had good success building dozens of moderately complex boards using this setup
with [Chip Quik SMD291AX](https://www.digikey.com/en/products/detail/chip-quik-inc/SMD291AX50T3/5130159). I haven't
tried non-leaded solder but I suspect this toaster may struggle to reach the required temperature or heat up fast
enough.

The control algorithm is currently on-off (bang-bang). I'm certain better control is possible here (e.g. PID) but it
works just fine as-is for most home projects (see graph).

# Hardware

The hardware for this can be assembled quite easily.

- [Hamilton Beach 31401 toaster oven](https://hamiltonbeach.com/4-slice-toaster-oven-31401)
- [Digital Loggers IoT Power Relay](https://www.digital-loggers.com/iot2.html)
- Raspberry pi
- [mcp9600 themocouple amplifier](https://www.adafruit.com/product/4101)
- [themocouple](https://www.adafruit.com/product/270)

It uses an unmodified Hamilton Beach 31401 which was the wirecutter recommended budget pick at the time of writing. The
key is that it has fully manual controls. Every other DIY reflow oven project I've found used a
toaster oven that was discontinued (and perhaps that is also true by the time you're reading this).

In order to use an unmodified toaster oven, this uses a Digital Loggers IoT Power Relay. This turns the toaster oven
on and off according to the control algorithm. To measure temperature, it uses a thermocouple and an mcp9600
thermocouple amplifier.

Like other DIY reflow ovens out there, you have to open the toaster oven door when the temperature when it is time to
cool off, otherwise the temperature of your PCB will not fall fast enough.




