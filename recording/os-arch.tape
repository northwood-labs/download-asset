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
Enter
Type "# Get the *current* OS and CPU architecture..."
Enter
Type "download-asset os-arch"
Sleep 500ms
Enter

Sleep 1s
Enter
Type "# When the current OS matches 'linux', return a different value..."
Enter
Type "download-asset os-arch --linux Linux"
Enter

Sleep 1s
Enter
Type "# When the current CPU architecture matches 'arm64', return a different value..."
Enter
Type "download-asset os-arch --linux Linux --arm64 AArch64"
Enter

Sleep 1s
Enter
Type "# Instead of {os}/{arch}, use a different pattern for matches..."
Enter
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
