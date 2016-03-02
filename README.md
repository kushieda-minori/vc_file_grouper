# vc_file_grouper
application that tries to group together game data files from [Valkyrie Crusade](http://mynet.co.jp/service/valkyrie.html) mobile app game

#Running
Compile from source or use an executable from the [latest release](https://github.com/kushieda-minori/vc_file_grouper/releases/latest). Use the 64bit version only if you have a 64bit operating system.

Run from the command line:

##Windows
* Open a command prompt (```cmd.exe``` or powershell)
* Change to the directory where the executable is: ```cd "c:\Users\MyUserName\Downloads"```
* Run the executable for your system
 * For Windows XP and 32bit Windows 7 run: ```vc_file_grouper_Win32.exe```
 * For Windows 8+ run: ```vc_file_grouper_Win64.exe```
* If you want to specify the location of your data files you can run it like this:<br />
 ```vc_file_grouper_Win64.exe "c:\Users\MyUserName\Downloads\My VC Data"```

##OSX
* Open a terminal (Applications -> Utilities -> Terminal, or I prefer [iTerm2](https://www.iterm2.com/)
* Change to the directory where the executable is: ```cd "~/Downloads"```
* Run the executable for your system: ```vc_file_grouper_OSX```
* If you want to specify the location of your data files you can run it like this:<br />
 ```vc_file_grouper_OSX "~/Downloads/My VC Data"```

##Unix
I shouldn't have to tell you. It's basically the same as OSX anyway.

#Using the program
The program starts a web-service. You should see a URL print in your terminal/command promt that looks like http://localhost:8585/ . Open this URL up in your favorite broswer (IE, Firefox, Chrome, etc). Once the application opens in your browser, you are set to go.

If you didn't specify a datafile location on the command line, or wish to change it, you can do so from within the web-application.
