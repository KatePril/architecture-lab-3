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
  -d 'bgrect 100 100 300 300' \
  http://localhost:17000/

# DEBUG POINT START (REMOVE WHEN SENDING)

curl -X POST \
  -H "Content-Type: text/plain" \
  -d 'update' \
  http://localhost:17000/

sleep 2

# DEBUG POINT END

curl -X POST \
  -H "Content-Type: text/plain" \
  -d 'figure 200 200' \
  http://localhost:17000/

curl -X POST \
  -H "Content-Type: text/plain" \
  -d 'update' \
  http://localhost:17000/

count=0
step=2

while [ $count -lt 50 ]
do
    count=$((count+1))
    curl -X POST \
      -H "Content-Type: text/plain" \
      -d "move $((200+count*step)) $((200+count*step))" \
      http://localhost:17000/

    curl -X POST \
      -H "Content-Type: text/plain" \
      -d 'update' \
      http://localhost:17000/

done