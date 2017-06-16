# CADUBI - Creative ASCII Drawing Utility By Ian

This Perl program was written sometime around 1997 -- many years before I knew
anything about Perl or Unix -- and the code hasn't changed much since then. It
has, however, become very widely adopted and is available as a package on most
Linux distributions.

CADUBI is an application written in Perl that allows you to draw
text-based images that are viewable on typical unix-based consoles.
Usually the applications that emulate these consoles support various
text modes, such as background and foreground colors, bold, and
inverse. This text art, commonly called "ASCII art," is used in
various places such as online BBSes, email and login prompts.

## Screen shots

<img src="https://raw.githubusercontent.com/statico/cadubi/master/screenshots/catfile.png" width="150"/> <img src="https://raw.githubusercontent.com/statico/cadubi/master/screenshots/edit.png" width="150"/> <img src="https://raw.githubusercontent.com/statico/cadubi/master/screenshots/help.png" width="150"/>

## Installation

* [Homebrew](http://brew.sh/): `brew install cadubi`
* [Debian](https://packages.debian.org/wheezy/cadubi) and [Ubuntu](http://packages.ubuntu.com/precise/cadubi): `apt-get install cadubi`
* From source:
  * Perl 5.002 or later
  * `Term::ReadKey` Perl module

## Usage

CADUBI has a 'pen' which describes the current mode. Properties of
the pen are the painting character, foreground color, background
color, bold, inverse, and blink. Whenever you paint or use the text
mode, the characters drawn on the screen will have the properties of
the pen. The current mode of the pen is shown at the bottom of the
console and is what will be drawn on screen when you paint.

Move around the cursor with the <kbd>i</kbd>, <kbd>j</kbd>, <kbd>k</kbd> and <kbd>l</kbd> keys. Holding
down shift and typing these keys will move the cursor five spaces
instead of one. Pressing return/enter will move the cursor down one
line and all the way to the left of the console.

To paint the current pen on the screen, press the space bar. To
delete a character, press the delete/backspace key. You'll notice
that editing is much like common text editors, such as pico or joe.
You can also delete with the <kbd>tilde</kbd> key, which makes moving & painting
(right hand) and erasing (left hand) much easier.

The pen character is the character that is drawn when you paint using
the space bar. To change the character, press <kbd>p</kbd> and then the
character you would like it to be.

To set the foreground or background colors for the cursor, press <kbd>f</kbd>
for foreground or <kbd>b</kbd> for background, and then a corresponding color
    code. The color codes are case-insensitive and are listed below:

     0 or N   Normal (standard text)
     1 or W   White
     2 or R   Red
     3 or G   Green
     4 or Y   Yellow
     5 or B   Blue
     6 or M   Magenta
     7 or C   Cyan
     8 or K   Black
     
If you can't remember the codes above, you can always hit <kbd>Ctrl-h</kbd>
to view the Quick Help which will display a summary of all the keys,
color codes and examples of how they look.

Typically, foreground text colors are the same as background colors,
unless the text is bold. If the text is bold, foreground colors are
usually lighter than the background color, making text easier to read
when the text has the same foreground and background color. Refer to
the Quick Help (<kbd>Ctrl-h</kbd>) to see what the colors look like on your
console.

Bold and inverse are two widely-supported modes. Bold is toggled with
the <kbd>g</kbd> key, and inverse is toggled with the <kbd>v</kbd> key. Blink, though
regarded as highly annoying, can be toggled with by pressing
<kbd>Shift-w</kbd>.

Text mode is an extremely useful feature. Once in the text mode you
can type as if you were using a normal text editor, and all the
characters drawn onscreen will use the same mode as the pen. To enter
text mode, press the <kbd>t</kbd> key. To exit, press escape.

To exit the CADUBI application, press <kbd>Ctrl-x</kbd>. Quick help can be
accessed by pressing <kbd>Ctrl-h</kbd>. In case it is needed, pressing
<kbd>Ctrl-w</kbd> will refresh the entire screen by redrawing each
character.

### Reading and writing files

To read a file and use it with CADUBI, type <kbd>Ctrl-r</kbd>. To write
a file, type <kbd>Ctrl-o</kbd>. You will be prompted for a filename.

When CADUBI reads a file, it will only read as much that will fit in
the workspace (the area of the console minus the bottom row [status
bar]). To gain more workspace, see the `-s` operator in 'COMMAND LINE
USAGE' below.

CADUBI optimizes its output files to display properly and take up as
little space as possible. All CADUBI output can be viewed with the
'cat' utility.

## Command line usage

    Usage: cadubi [OPTIONS] [FILE]

    Available options:
      -h, --help              what you're looking at now
      -m, --mute              turn off beeping
      -s, --size [W] [H]      sets the size of the console for use with
                              CADUBI, where W is number of columns and H
                              is number of rows.
      -v, --version           show CADUBI version

    Example:
      Will make the cadubi workspace 160 columns wide, 48 rows high,
      disable beeping, and open the file 'bacon.txt':
      
          cadubi --mute --size 160 48 bacon.txt

      Will display the version of CADUBI, copyright and author:

          cadubi -v

## Tips

Whenever you are prompted to type in information, such as the name of
a file to read/write to, you can hit escape to cancel. You can also
hit escape to get out of text mode.

When using the `-s` or `--size` command line option, make sure your
console actually _is_ that size or the text won't wrap properly and
CADUBI will look funny.
