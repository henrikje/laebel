#!/usr/bin/env bash
set -e

# Take screenshots of Laebel's home screen
echo "Taking screenshots of Laebel's home screen..."
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
  --headless --screenshot=laebel-example-screenshot.png --window-size=1024,1024 --hide-scrollbars \
  http://localhost:8000/

# For some odd reason, the graph is not properly rendered if I limit the window height to 868px,
# so I make it larger and crop the image afterwards.
echo "Cropping the screenshot..."
magick laebel-example-screenshot.png -crop 1024x863+0+0 laebel-example-screenshot.png

# Overlay screenshot on template
echo "Overlaying the screenshot on the template..."
magick screenshot-frame.png laebel-example-screenshot.png -geometry +26+66 -composite laebel-example-screenshot.png

# Opening the screenshot
echo "Opening the screenshot..."
open laebel-example-screenshot.png

# The web page is updated manually:
# 1. Open a new Private browsing window in Firefox (to avoid extensions polluting the HTML)
# 2. Navigate to http://localhost:8000/
# 3. Right-click the page and select "Save Page As..."
# 4. Click Cmd+Shift+G, and enter: /Users/henrikjernevad/Dropbox/Development/Go/laebel/examples/react-express-mysql
# 5. Enter "laebel-output.html" as name, and click Save.
# 6. Click "Replace".
echo "Don't forget to update the web page manually."