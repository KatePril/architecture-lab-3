#!/bin/bash
curl -X POST \
  -H "Content-Type: text/plain" \
  -d 'reset' \
  http://localhost:17000/

curl -X POST \
  -H "Content-Type: text/plain" \
  -d 'green' \
  http://localhost:17000/

curl -X POST \
  -H "Content-Type: text/plain" \
  -d 'bgrect 0.25 0.25 0.75 0.75' \
  http://localhost:17000/


curl -X POST \
  -H "Content-Type: text/plain" \
  -d 'figure 0.5 0.5' \
  http://localhost:17000/

curl -X POST \
  -H "Content-Type: text/plain" \
  -d 'update' \
  http://localhost:17000/

count=0
steps=50
step=$(echo "scale=4; 1.0 / $steps" | bc)

while [ $count -lt $steps ]
do
    pos=$(echo "scale=4; $count * $step" | bc)
    curl -X POST \
      -H "Content-Type: text/plain" \
      -d "move $pos $pos" \
      http://localhost:17000/

    curl -X POST \
      -H "Content-Type: text/plain" \
      -d 'update' \
      http://localhost:17000/

    count=$((count + 1))
    sleep 0.05
done