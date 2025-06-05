package app

import (
	"crypto/md5"
	"encoding/hex"
	"selfhosted/static"
)

var CssVersion string

func init() {
	CssVersion = StaticAssetVersion("css/dist.css")
}

func StaticAssetVersion(path string) string {
	file, err := static.FS.ReadFile(path)
	if err != nil {
		panic(err)
	}

	version := md5.Sum(file)
	return hex.EncodeToString(version[:])
}
