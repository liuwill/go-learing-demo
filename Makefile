m3u8:
	cd demo/m3u8 && go run main.go --mode=remote --name="${TARGET}"

m3u8-simple:
	cd demo/m3u8 && echo "${PATH_INFO}" && go run main.go --mode=simple --info="${PATH_INFO}" --name="${TARGET}"
