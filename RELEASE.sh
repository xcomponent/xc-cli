if [ $# -eq 0 ]
  then
    echo "Tag has not been specified"
    exit 1
fi

rm -rf ./dist && \
  gox -osarch="darwin/amd64 linux/386 linux/amd64 windows/386 windows/amd64" \
      -output=dist/{{.Dir}}_{{.OS}}_{{.Arch}} && \
  zip -mqj dist/xc-windows-64.zip dist/xc-cli_windows_amd64.exe && \
  zip -mqj dist/xc-windows-32.zip dist/xc-cli_windows_386.exe && \
  zip -mqj dist/xc-linux-32.zip dist/xc-cli_linux_386 && \
  zip -mqj dist/xc-linux-64.zip dist/xc-cli_linux_amd64 && \
  zip -mqj dist/xc-macos-64.zip dist/xc-cli_darwin_amd64

ghr -u xcomponent $1 dist/
