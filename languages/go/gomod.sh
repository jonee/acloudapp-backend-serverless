
if [ -f ./go.mod ]; then
    exit 0
fi

go mod init acloudapp.org
