set GOPATH=C:\agent_updater\agent\shared
:: copy gcpp into current folder

REM go get github.com/jstemmer/gotags
REM go get github.com/sadew7/gcpp

for /r . %%g in (*.go) do gcpp -D=EPD -i=%%~nxg -o=./ht/%%~nxg
for /r . %%g in (*.go) do gcpp -D=STELLUS -i=%%~nxg -o=./ti/%%~nxg
