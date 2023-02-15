package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func pobieransko(url, path string) {
	fmt.Println("UWAGA!!!! Rozpoczynam pobierańsko")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	f, _ := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"trwa pobierańsko",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
}

func rozpakowywansko(zipname string) {
	fmt.Println("rozpakowuje")
	archive, err := zip.OpenReader(zipname)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	dirname := zipname[:len(zipname)-4]
	os.Mkdir(dirname, 0644)

	bar := progressbar.Default(
		int64(len(archive.File)),
		"trwa rozpakowywańsko",
	)

	for _, f := range archive.File {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}
		plik, _ := f.Open()
		defer plik.Close()

		fp, err := os.Create(dirname + "/" + filepath.Base(f.Name))
		if err != nil {
			panic(err)
		}
		defer fp.Close()
		io.Copy(fp, plik)
		bar.Add(1)
	}
}

func odpalansko(argumenty ...string) {
	fmt.Println("\nStartujemy!!!!!!!!!!!!!!\n")

	c := exec.Command(argumenty[0], argumenty[1:]...)
	c.Stderr, c.Stdout, c.Stdin = os.Stderr, os.Stdout, os.Stdin
	c.Run()
}

func main() {
	hostname, _ := os.Hostname()

	var miner string

	if len(os.Args) > 1 {
		miner = os.Args[1]
	} else {
		fmt.Print(`Witamy w skrypciku górniczym! :>
Jeśli możesz to uruchom skrypcik jako administrator inaczej bedzie strasznie wolno kopało

W razie problemow napisz na aleksander.wajcht2021p@zsepoznan.pl

Wybierz kopareczke: (jak nie wiesz co to wybierz 1)
 1. xmrig
 2. cpuminer
>>> `)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadByte()
		miner = string(input)
	}

	if miner == "1" {
		pobieransko("https://github.com/xmrig/xmrig/releases/download/v6.17.0/xmrig-6.17.0-gcc-win64.zip", "xmrig.zip")

		rozpakowywansko("xmrig.zip")

		odpalansko(".\\xmrig\\xmrig.exe",
			"-p", "x",
			"-o", "xmr.2miners.com:2222",
			"-u", "43bEsE2i48KAez5aQ3FUq54dfGfQ3HSdZitYxRFfdS93dY586VueJfHdH59o5gUSzgfyWCnkonZhyWh8P4GuU6RX2dXYB5k."+hostname,
		)

	} else if miner == "2" {
		pobieransko("https://github.com/JayDDee/cpuminer-opt/releases/download/v3.20.3/cpuminer-opt-3.20.3-windows.zip", "cpuminer.zip")

		rozpakowywansko("cpuminer.zip")

		odpalansko(".\\cpuminer\\cpuminer-sse2.exe",
			"--algo=allium",
			"-o", "stratum+tcp://freshgarlicblocks.net:3032",
			"-u", "MUdUpEPM1aJJqQJNJUqaPjdujTo854noqn",
			"-p", "x",
		)
	} else {
		fmt.Println("nie ma takiej opcji!!!!!!!!!!!!!")
	}
}
