#!/bin/bash

for f in *.tetrisviz; do
    tetrisviz --format pikchr "$f"
    tetrisviz --format svg "$f"
done

echo "# Examples" > examples.md

echo "<table>" >> examples.md
echo "<tr><th>tetrisviz</th><th>pikchr</th><th>svg</th></tr>" >> examples.md

for f in *.tetrisviz; do
    f="${f%.*}"
    echo "<tr>" >> examples.md
    { echo "<td><pre>"; cat "$f.tetrisviz"; echo; echo "</pre></td>"; } >> examples.md
    { echo "<td><pre>"; cat "$f.pikchr"; echo; echo "</pre></td>"; } >> examples.md
    echo "<td><img src=\"$f.svg\" style=\"width: 200px\"/></td>" >> examples.md
    echo "</tr>" >> examples.md
done

echo "</table>" >> examples.md
