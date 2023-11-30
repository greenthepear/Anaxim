**Anaxi** *(named after the ancient Greek philosopher and cartographer Anaximander)* is a civilisation simulator made in ~~Ebitengine~~ Fyne, with the main goal of making a visually appealing and interactive timeline of the world using cellular automata.

# State
Currently the "game" only features a basic migration simulation on a map (defmap.png) where you can ~~see the population of any cell with your mouse cursor~~* and control the speed. Some statistics are displayed on the right. ~~Clicking gives some debugging info in the terminal.~~*

\* *Migration to Fyne got rid of these functionalities for now, working on it.*

# Building
To build you need Go (1.21.4), a GCC compiler (for cgo) and depending on platform, some dependencies. Info below is mostly from the [Fyne.io guide](https://developer.fyne.io/started/), which you'll need if you want to build on different platforms than Windows and Linux and **please refer to it if you have issues**.

**Warning! `go build` will seem to hang for a while the first time you try to compile as it needs to build the graphics drivers using GCC.** *"So much for fast Go compile times!" well, technically it's not Go. Feel free to fork Fyne to work without cgo.*

## Linux

1. Install [Go](https://go.dev/doc/install)
2. Install gcc
3. Install graphics library header files 
    - Debian / Ubuntu: `sudo apt-get install libgl1-mesa-dev xorg-dev`
    - Fedora: `sudo dnf install libXcursor-devel libXrandr-devel mesa-libGL-devel libXi-devel libXinerama-devel libXxf86vm-devel`
    - Arch Linux: `sudo pacman -S xorg-server-devel libxcursor libxrandr libxinerama libxi`
    - Solus: `sudo eopkg it -c system.devel mesalib-devel libxrandr-devel libxcursor-devel libxi-devel libxinerama-devel`
    - openSUSE: `sudo zypper install libXcursor-devel libXrandr-devel Mesa-libGL-devel libXi-devel libXinerama-devel libXxf86vm-devel`
    - Void Linux: `sudo xbps-install -S base-devel xorg-server-devel libXrandr-devel libXcursor-devel libXinerama-devel`
    - Alpine Linux: `sudo apk add libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev linux-headers mesa-dev`
    - NixOS: `nix-shell -p libGL pkg-config xorg.libX11.dev xorg.libXcursor xorg.libXi xorg.libXinerama xorg.libXrandr xorg.libXxf86vm`
4. Compile! (See warning above)

        git clone https://github.com/greenthepear/Anaxi.git
        cd Anaxi
        go build

## Windows

1. Install [Go](https://go.dev/doc/install)
2. Install gcc
    - Save yourself the trouble and get [TDM-GCC](https://jmeubank.github.io/tdm-gcc/download/). You want the 64+32-bit MinGW-w64 edition.
4. Compile! (See warning above)

        git clone https://github.com/greenthepear/Anaxi.git
        cd Anaxi
        go build
          

# Running

You need to have the defmap.png file in the same directory as the executable. To change this run the program with a different `mappath` flag. This is explained below.

## Flags

Optionally, you can set flags from the command-line when launching Anaxi. They can be checked `./Anaxi -h`

    -mappath string
        Path to the map PNG file. (default "./defmap.png")
    -prerun int
        Generations to simulate before launching, min 50

So for example `./Anaxi -mappath=./defmapDEBUG.png -prerun=20000` will prerun the simulation for 20,000 generations on a map generated from defmapDEBUG.png.

## Map file

The PNG map file is converted to a map grid like this:
- All pixels with RGB = [0,0,255] become water cells.
- All other pixel become land cells.
- Their habitability level is determined by (255-R)/255, so more red is less habitable.

# Planned features
- ~~Underlying world map~~
- Tribes and countries
- Wars and diplomacy
- ~~Speed/pause controls~~
- GPU acceleration with OpenCL (scary!)
- Interactivity
- Stats and graphs!
- ~~Custom maps~~ (and random maps?)
- Map modes
