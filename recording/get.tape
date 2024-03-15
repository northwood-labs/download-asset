Output recording/get.gif
Output recording/get.ascii

Set Shell bash

Require download-asset

Set Width 1200
Set Height 600
Set FontSize 18
Set FontFamily "JetBrainsMono Nerd Font"
Set LineHeight 1.2
Set Theme "Galaxy"
Set WindowBar Colorful
Set Padding 30
Set Margin 0
Set MarginFill "#ffffff"
Set BorderRadius 10
Set Framerate 60
Set TypingSpeed 75ms
Set CursorBlink true

Sleep 1s
Type "download-asset get --help | less"
Sleep 500ms
Enter
Sleep 2s
Space
Sleep 2s
Space
Sleep 2s
Space
Sleep 2s
Space
Sleep 2s
Space
Type q

Sleep 1s
Hide
Type "clear"
Enter
Show
Type "download-asset get --verbose --owner-repo opentofu/opentofu"
Sleep 500ms
Enter

Sleep 8s

Sleep 1s
Hide
Type "clear"
Enter
Show
Type "download-asset get --verbose --owner-repo gruntwork-io/terragrunt"
Sleep 500ms
Enter

Sleep 15s

Sleep 1s
Hide
Type "clear"
Enter
Show
Type "download-asset get --verbose --owner-repo direnv/direnv"
Sleep 500ms
Enter

Sleep 5s

Sleep 1s
Hide
Type "alias ll='eza -aa -l -g -h --git --modified --time-style=long-iso -H --classify --icons --group-directories-first' && clear"
Enter
Show
Type "ll /usr/local/bin"
Sleep 500ms
Enter

Sleep 5s
Ctrl+D
