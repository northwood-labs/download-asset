Output recording/latest-tag.gif
Output recording/latest-tag.ascii

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
Type "# Get the latest 'release' of OpenTofu..."
Enter
Type "download-asset latest-tag --owner-repo 'opentofu/opentofu'"
Sleep 500ms
Enter

Sleep 1s
Enter
Type "# With verbose logging..."
Enter
Type "download-asset latest-tag --owner-repo 'opentofu/opentofu' --verbose"
Enter

Sleep 1s
Enter
Type "# Ignore releases and go straight to the list of tags (may include pre-releases)..."
Enter
Type "download-asset latest-tag --owner-repo 'opentofu/opentofu' --skip-to-tags"
Enter

Sleep 1s
Enter
Type "# Now, let's do Terraform..."
Enter
Type "download-asset latest-tag --owner-repo 'hashicorp/terraform'"
Enter

Sleep 1s
Type "download-asset latest-tag --owner-repo 'hashicorp/terraform' --skip-to-tags"
Enter

Sleep 1s
Enter
Type "# Let's add a version constraint..."
Enter
Type "download-asset latest-tag --owner-repo 'hashicorp/terraform' --constraint '< 1.6.0'"
Enter

Sleep 1s
Enter
Type "# Some repositories do not use 'releases' at all, and only use tags..."
Enter
Type "download-asset latest-tag --owner-repo 'golang/go'"
Enter

Sleep 5s
Ctrl+D
