# zap demo - github star magnet
# rounded corners, perfect pacing, visually stunning

Output assets/demo.gif

Set FontSize 18
Set Width 1200
Set Height 700
Set Padding 30
Set Margin 50
Set MarginFill "#000000"
Set BorderRadius 12
Set WindowBar "Colorful"
Set Framerate 60
Set PlaybackSpeed 1.0
Set TypingSpeed 50ms

Set Theme { "name": "Pitch Black", "background": "#000000", "foreground": "#cdd6f4", "cursor": "#f5e0dc", "black": "#000000", "brightBlack": "#585b70", "red": "#f38ba8", "brightRed": "#f38ba8", "green": "#a6e3a1", "brightGreen": "#a6e3a1", "yellow": "#f9e2af", "brightYellow": "#f9e2af", "blue": "#89b4fa", "brightBlue": "#89b4fa", "magenta": "#cba6f7", "brightMagenta": "#cba6f7", "cyan": "#89dceb", "brightCyan": "#89dceb", "white": "#cdd6f4", "brightWhite": "#ffffff" }

# setup - ensure binary is built
Hide
Type "cd /home/programmersd/Projects/github/zap"
Enter
Type "go build -o zap ./cmd/zap"
Enter
Sleep 2s
Type "mkdir -p /tmp/zap-demo/photos /tmp/zap-demo/backup"
Enter
Type "for i in {1..8}; do dd if=/dev/zero of=/tmp/zap-demo/photos/img_482$i.jpg bs=1M count=16 2>/dev/null; done"
Enter
Sleep 2s
Type "dd if=/dev/zero of=/tmp/zap-demo/video.mp4 bs=1M count=120 2>/dev/null"
Enter
Sleep 2s
Type "clear"
Enter
Sleep 400ms
Show

# stunning intro - show the gradient banner with help
Type "./zap --help"
Sleep 700ms
Enter
Sleep 3500ms

Type "clear"
Sleep 300ms
Enter
Sleep 700ms

# demo 1: multiple files - show the beautiful progress
Type "# copy 8 photos (128 mb)"
Sleep 900ms
Enter
Sleep 400ms
Type "./zap /tmp/zap-demo/photos/*.jpg /tmp/zap-demo/backup/"
Sleep 900ms
Enter
Sleep 5200ms

Sleep 1400ms

Type "clear"
Sleep 300ms
Enter
Sleep 700ms

# demo 2: single large file - showcase speed calculation
Type "# single large file with live speed tracking"
Sleep 1100ms
Enter
Sleep 400ms
Type "./zap /tmp/zap-demo/video.mp4 /tmp/zap-demo/backup/"
Sleep 900ms
Enter
Sleep 4800ms

Sleep 1400ms

Type "clear"
Sleep 300ms
Enter
Sleep 700ms

# demo 3: verbose mode - show individual file tracking
Type "# verbose mode - track every file"
Sleep 1100ms
Enter
Sleep 400ms
Type "./zap -v -m /tmp/zap-demo/backup/*.jpg /tmp/zap-demo/photos/"
Sleep 900ms
Enter
Sleep 5500ms

Sleep 1400ms

Type "clear"
Sleep 300ms
Enter
Sleep 700ms

# demo 4: show completion message beauty
Type "# lightning fast operations"
Sleep 900ms
Enter
Sleep 400ms
Type "./zap /tmp/zap-demo/video.mp4 /tmp/zap-demo/final.mp4 -f"
Sleep 700ms
Enter
Sleep 4500ms

Sleep 1600ms

Type "clear"
Sleep 300ms
Enter
Sleep 500ms

# finale - back to help to show the full feature set
Type "./zap --help"
Sleep 500ms
Enter
Sleep 3500ms

# cleanup
Hide
Type "rm -rf /tmp/zap-demo"
Enter
