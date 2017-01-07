#!/bin/bash

#!/bin/sh

echo "Run ./terminfo tests"
go test ./terminfo/
echo "Run ./termios tests"
go test ./termios/
echo "Done."
