#!/bin/sh

$ gpg --batch --gen-key <<EOF
    Key-Type: rsa
    Key-Length: 4096
    Name-Real: My Name
    Name-Comment: For Testing
    Name-Email: someuser@mail.com
    Expire-Date: 0
    %pubring the_key.pub
    %secring the_key.sec
    %commit
    %echo done
EOF
