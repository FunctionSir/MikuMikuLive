# Makefile for MikuMikuLive
default:
	@echo '__  __ __  __ _'
	@echo '|  \/  |  \/  | |'
	@echo '| |\/| | |\/| | |'
	@echo '| |  | | |  | | |___'
	@echo '|_|  |_|_|  |_|_____|'
	@echo 'MikuMikuLive Version 0.1-alpha (Miku)'
	@echo Building using GOOS and GOARCH you set before...
	@mkdir MikuMikuLive
	@go build -o MikuMikuLive/mml
	@cp root.html MikuMikuLive
	@cp media.d MikuMikuLive -r
	@cp conf.d MikuMikuLive -r
	@cp ldf.d MikuMikuLive -r
	@cp README.md MikuMikuLive
	@cp LICENSE MikuMikuLive
	@echo Done! You can copy/move the dir MikuMikuLive to anywhere you want and use it.

linux-amd64:
	@echo '__  __ __  __ _'
	@echo '|  \/  |  \/  | |'
	@echo '| |\/| | |\/| | |'
	@echo '| |  | | |  | | |___'
	@echo '|_|  |_|_|  |_|_____|'
	@echo 'MikuMikuLive Version 0.1-alpha (Miku)'
	@echo Building for linux/amd64...
	@mkdir MikuMikuLive
	@GOOS=linux GOARCH=amd64 go build -o MikuMikuLive/mml
	@cp root.html MikuMikuLive
	@cp media.d MikuMikuLive -r
	@cp conf.d MikuMikuLive -r
	@cp ldf.d MikuMikuLive -r
	@cp README.md MikuMikuLive
	@cp LICENSE MikuMikuLive
	@echo Done! You can copy/move the dir MikuMikuLive to anywhere you want and use it.

linux-arm64:
	@echo '__  __ __  __ _'
	@echo '|  \/  |  \/  | |'
	@echo '| |\/| | |\/| | |'
	@echo '| |  | | |  | | |___'
	@echo '|_|  |_|_|  |_|_____|'
	@echo 'MikuMikuLive Version 0.1-alpha (Miku)'
	@echo Building for linux/arm64...
	@mkdir MikuMikuLive
	@GOOS=linux GOARCH=arm64 go build -o MikuMikuLive/mml
	@cp root.html MikuMikuLive
	@cp media.d MikuMikuLive -r
	@cp conf.d MikuMikuLive -r
	@cp ldf.d MikuMikuLive -r
	@cp README.md MikuMikuLive
	@cp LICENSE MikuMikuLive
	@echo Done! You can copy/move the dir MikuMikuLive to anywhere you want and use it.

linux-ppc64le:
	@echo '__  __ __  __ _'
	@echo '|  \/  |  \/  | |'
	@echo '| |\/| | |\/| | |'
	@echo '| |  | | |  | | |___'
	@echo '|_|  |_|_|  |_|_____|'
	@echo 'MikuMikuLive Version 0.1-alpha (Miku)'
	@echo Building for linux/ppc64le...
	@mkdir MikuMikuLive
	@GOOS=linux GOARCH=ppc64le go build -o MikuMikuLive/mml
	@cp root.html MikuMikuLive
	@cp media.d MikuMikuLive -r
	@cp conf.d MikuMikuLive -r
	@cp ldf.d MikuMikuLive -r
	@cp README.md MikuMikuLive
	@cp LICENSE MikuMikuLive
	@echo Done! You can copy/move the dir MikuMikuLive to anywhere you want and use it.

freebsd-amd64:
	@echo '__  __ __  __ _'
	@echo '|  \/  |  \/  | |'
	@echo '| |\/| | |\/| | |'
	@echo '| |  | | |  | | |___'
	@echo '|_|  |_|_|  |_|_____|'
	@echo 'MikuMikuLive Version 0.1-alpha (Miku)'
	@echo Building for freebsd/amd64...
	@mkdir MikuMikuLive
	@GOOS=freebsd GOARCH=amd64 go build -o MikuMikuLive/mml
	@cp root.html MikuMikuLive
	@cp media.d MikuMikuLive -r
	@cp conf.d MikuMikuLive -r
	@cp ldf.d MikuMikuLive -r
	@cp README.md MikuMikuLive
	@cp LICENSE MikuMikuLive
	@echo Done! You can copy/move the dir MikuMikuLive to anywhere you want and use it.

freebsd-arm64:
	@echo '__  __ __  __ _'
	@echo '|  \/  |  \/  | |'
	@echo '| |\/| | |\/| | |'
	@echo '| |  | | |  | | |___'
	@echo '|_|  |_|_|  |_|_____|'
	@echo 'MikuMikuLive Version 0.1-alpha (Miku)'
	@echo Building for freebsd/arm64...
	@mkdir MikuMikuLive
	@GOOS=freebsd GOARCH=arm64 go build -o MikuMikuLive/mml
	@cp root.html MikuMikuLive
	@cp media.d MikuMikuLive -r
	@cp conf.d MikuMikuLive -r
	@cp ldf.d MikuMikuLive -r
	@cp README.md MikuMikuLive
	@cp LICENSE MikuMikuLive
	@echo Done! You can copy/move the dir MikuMikuLive to anywhere you want and use it.

windows-amd64:
	@echo '__  __ __  __ _'
	@echo '|  \/  |  \/  | |'
	@echo '| |\/| | |\/| | |'
	@echo '| |  | | |  | | |___'
	@echo '|_|  |_|_|  |_|_____|'
	@echo 'MikuMikuLive Version 0.1-alpha (Miku)'
	@echo Building for windows/amd64...
	@mkdir MikuMikuLive
	@GOOS=windows GOARCH=amd64 go build -o MikuMikuLive/mml.exe
	@cp root.html MikuMikuLive
	@cp media.d MikuMikuLive -r
	@cp conf.d MikuMikuLive -r
	@cp ldf.d MikuMikuLive -r
	@cp README.md MikuMikuLive
	@cp LICENSE MikuMikuLive
	@echo Done! You can copy/move the dir MikuMikuLive to anywhere you want and use it.

windows-arm64:
	@echo '__  __ __  __ _'
	@echo '|  \/  |  \/  | |'
	@echo '| |\/| | |\/| | |'
	@echo '| |  | | |  | | |___'
	@echo '|_|  |_|_|  |_|_____|'
	@echo 'MikuMikuLive Version 0.1-alpha (Miku)'
	@echo Building for windows/arm64...
	@mkdir MikuMikuLive
	@GOOS=windows GOARCH=arm64 go build -o MikuMikuLive/mml.exe
	@cp root.html MikuMikuLive
	@cp media.d MikuMikuLive -r
	@cp conf.d MikuMikuLive -r
	@cp ldf.d MikuMikuLive -r
	@cp README.md MikuMikuLive
	@cp LICENSE MikuMikuLive
	@echo Done! You can copy/move the dir MikuMikuLive to anywhere you want and use it.

clean:
	@echo '__  __ __  __ _'
	@echo '|  \/  |  \/  | |'
	@echo '| |\/| | |\/| | |'
	@echo '| |  | | |  | | |___'
	@echo '|_|  |_|_|  |_|_____|'
	@echo 'MikuMikuLive Version 0.1-alpha (Miku)'
	@echo Cleaning...
	@rm MikuMikuLive -rf
	@echo Done!