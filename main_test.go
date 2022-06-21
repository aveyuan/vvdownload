package main

import "testing"

const url = "https://oss.vvcms.cn/uploads/files/2022/06/09/ad4314c72b104e9faac4df45f083775d.zip"
func TestDownload(t *testing.T) {
	path, err := Download(url)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
}