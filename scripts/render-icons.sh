#!/bin/bash
# Renders the SVG icons to PNGs using headless Chrome
# (ImageMagick's built-in SVG renderer drops strokes):
#   forms/icons/<name>.png    - 32x32 for the toolbar
#   forms/icons/<name>-16.png - 16x16 for menus, with thicker lines so they stay crisp
#   icon.png                  - 256x256 application icon
set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

CHROME=google-chrome
if [ -x "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" ]; then
    CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
fi

# render <svg> <size> <png>
render() {
    printf '<html><body style="margin:0;background:transparent"><img src="file://%s" width="%d" height="%d" style="display:block"></body></html>' "$1" "$2" "$2" > "$TMP/page.html"
    "$CHROME" --headless=new --disable-gpu --hide-scrollbars --default-background-color=00000000 \
        --force-device-scale-factor=1 --window-size="$2,$2" --screenshot="$3" "file://$TMP/page.html" >/dev/null 2>&1
}

cd "$ROOT/forms/icons"
for f in *.svg; do
    n=${f%.svg}
    render "$PWD/$f" 32 "$PWD/$n.png"
    sed 's/stroke-width="[0-9.]*"/stroke-width="1.8"/' "$f" > "$TMP/$n-16.svg"
    render "$TMP/$n-16.svg" 16 "$PWD/$n-16.png"
done

render "$ROOT/icon.svg" 256 "$ROOT/icon.png"
