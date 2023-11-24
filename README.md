**Anaxi** *(named after the ancient Greek philosopher and cartographer Anaximander)* is a civilisation simulator made in ~~Ebitengine~~ Fyne, with the main goal of making a visually appealing and interactive timeline of the world using cellular automata.

# State
Currently the "game" only features a basic migration simulation on a map (defmap.png) where you can ~~see the population of any cell with your mouse cursor~~ and control the speed. ~~Clicking gives some debugging info in the terminal.~~*

\* Migration to Fyne got rid of these functionalities for now, working on it.

# Building
If you have git and Go (at least 1.21.4):

    git clone https://github.com/greenthepear/Anaxi.git
    cd Anaxi
    go build

# Running

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
- Custom maps (and random maps?)
- Map modes
