Output recording/os-arch.gif
Output recording/os-arch.ascii

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
Type "download-asset os-arch --help | less"
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
Type q

Sleep 1s
Type "download-asset os-arch"
Sleep 500ms
Enter

Sleep 1s
Type "download-asset os-arch --linux Linux"
Enter

Sleep 1s
Type "download-asset os-arch --linux Linux --arm64 AArch64"
Enter

Sleep 1s
Type "download-asset os-arch --linux Linux --arm64 AArch64 --pattern ''"
Left
Type "{{}}"
Left 2
Type ".OS"
Right 2
Type " running on "
Type "{{}}"
Left 2
Type ".Arch"
Enter

Sleep 5s
Ctrl+D
