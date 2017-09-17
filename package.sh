rm -rf ./target &&
  echo "Purging target folder." &&
  mkdir ./target &&
  echo "Building XC CLI for multiple platform." &&
  gox -os="darwin linux windows" -arch="amd64 386" &&
  sleep 1
  echo "Packaging OSX 32 bit distribution." &&
  cp xc_darwin_386 target/xc | zip -mqj target/xc-osx-32.zip target/xc &&
  echo "Packaging OSX 64 bit distribution." &&
  cp xc_darwin_amd64 target/xc | zip -mqj target/xc-osx-64.zip target/xc &&
  echo "Packaging Linux 32 bit distribution." &&
  cp xc_linux_386 target/xc | zip -mqj target/xc-linux-32.zip target/xc &&
  echo "Packaging Linux 64 bit distribution." &&
  cp xc_linux_amd64 target/xc | zip -mqj target/xc-linux-64.zip target/xc &&
  echo "Packaging Windows 32 bit distribution." &&
  cp xc_windows_386.exe target/xc.exe | zip -mqj target/xc-windows-32.zip target/xc.exe &&
  echo "Packaging Windows 64 bit distribution." &&
  cp xc_windows_amd64.exe target/xc.exe | zip -mqj target/xc-windows-64.zip target/xc.exe
