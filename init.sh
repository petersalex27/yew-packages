if cd ./$1; then
    rm -f go.mod
    if [ -z $2 ]; then
        go mod init github.com/petersalex27/yew-packages/$1
    else
        go mod init github.com/petersalex27/yew-packages/$2
    fi
    #go mod tidy
    cd ../
else
    echo no directory named $1
fi