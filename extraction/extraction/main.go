package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type systemInfo struct {
	Name        string
	RAM         int64 //Store Bytes of Ram
	RAMSpeed    float64
	HDDSpeed    float64
	CacheMemory []uint64
	CPUCores    uint32
	CPUThreads  uint32
	CPUSpeed    float64
	MaxCPUSpeed float64
	Latency     float64
	DownSpeed   float64
	UpSpeed     float64
}

type dataset struct {
	ComputerPower float64 `json:"computerPower"`
	DownloadSpeed float64 `json:"downSpeed"`
	Ylabels       string  `json:"yLabels"`
}

func main() {
	//generateKeys()
	testing()
}

func testing() {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Password: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, text)

	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	// Display mnemonic and keys
	fmt.Println("Mnemonic: ", mnemonic)
	fmt.Println("Master private key: ", masterKey)
	fmt.Println("Master public key: ", publicKey)
}

func generateKeys() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	publicKey := &privateKey.PublicKey

	fmt.Println("Private Key ", privateKey)
	fmt.Println("Public Key ", publicKey)

}

func calculateSystemInfo() {
	host, err := ghw.Host()
	if err != nil {
		fmt.Println("Error in Getting Host Informations", err.Error())
	}

	info := systemInfo{}
	info.RAM = host.Memory.TotalUsableBytes
	info.Name = host.CPU.Processors[0].Model
	info.CPUCores = host.CPU.TotalCores
	info.CPUThreads = host.CPU.TotalThreads
	// fmt.Println("Printing BIOS information \t", host.BIOS.JSONString(true))
	// fmt.Println("Printing Baseboard information \t", host.Baseboard.JSONString(true))
	// fmt.Println("Printing Chassis information \t", host.Chassis.JSONString(true))
	// fmt.Println("Printing Memory information \t", host.Memory.JSONString(true))
	// fmt.Println("Printing Product information \t", host.Product.JSONString(true))

	//fmt.Println(host.Topology.JSONString(true))

	for _, node := range host.Topology.Nodes {
		for _, cache := range node.Caches {
			//fmt.Println("Cache Level Befor if ", cache.Level)
			if len(info.CacheMemory) < int(cache.Level) {
				for i := len(info.CacheMemory); i <= int(cache.Level)-1; i++ {
					//fmt.Println("Cache Level", cache.Level)
					info.CacheMemory = append(info.CacheMemory, 0)
				}
			}
		}
	}

	for _, node := range host.Topology.Nodes {
		for _, cache := range node.Caches {
			info.CacheMemory[cache.Level-1] += cache.SizeBytes
		}
	}
	// fmt.Println("Size of Cache Memory ", len(info.CacheMemory))
	// fmt.Println("Values of Cahce Memory ", info.CacheMemory)
	//fmt.Println(host.JSONString(true))
	fmt.Println(info)
	var speed string
	array := strings.Split(info.Name, " ")
	for x := 0; x < len(array); x++ {
		loCase := strings.ToLower(array[x])
		if strings.Contains(loCase, "ghz") {
			speed = loCase
		}

	}
	// fmt.Println(array)
	// fmt.Println(speed)
	speed = strings.TrimSuffix(speed, "ghz")
	// fmt.Println(speed)

	info.CPUSpeed, err = strconv.ParseFloat(speed, 64)
	if err != nil {
		fmt.Println("Error in Parsing Float Values ", err.Error())
	}
	// fmt.Println(info)
	info.MaxCPUSpeed = info.CPUSpeed * float64(info.CPUCores)
	// fmt.Println(info)

	// floatVar := 22.23423425523452
	// fmt.Println(floatVar)
	// fmt.Println(math.Floor(floatVar))
	// fmt.Println(math.Trunc(floatVar))
	// fmt.Println(math.Round(floatVar))
	// fmt.Println(math.Floor(floatVar))
	//println(host.JSONString(true))
	//ghw.Host()

	cmd := exec.Command("speedtest-cli")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	result := out.String()

	newResult := strings.Split(result, "\n")
	var download []string
	var upload []string
	var latency []string

	for _, line := range newResult {
		// fmt.Println("Line Number ", i, "\t Result is \t", strings.ToLower(line))
		if strings.Contains(strings.ToLower(line), "download") {
			download = strings.Split(strings.ToLower(line), " ")
		}
		if strings.Contains(strings.ToLower(line), "upload") {
			upload = strings.Split(strings.ToLower(line), " ")
		}
		if strings.Contains(strings.ToLower(line), "hosted by") {
			latency = strings.Split(strings.ToLower(line), " ")
		}

	}

	// fmt.Println("Download \t", download[1])
	// fmt.Println("Upload \t\t", upload[1])
	// fmt.Println("Latecny ", latency[len(latency)-2])
	info.DownSpeed, _ = strconv.ParseFloat(download[1], 64)
	info.UpSpeed, _ = strconv.ParseFloat(upload[1], 64)
	info.Latency, _ = strconv.ParseFloat(latency[len(latency)-2], 64)

	// fmt.Println("DMICODE RESULT", result)
	// fmt.Println("Missing lInes")

	cmd = exec.Command("lshw", "-json")
	cmd.Stdin = strings.NewReader("some input")
	// var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	data := dataset{}
	power := computerPower(float64(info.RAM)*0.000001, float64(info.CPUCores), info.UpSpeed)
	data.ComputerPower = power
	data.DownloadSpeed = info.DownSpeed

	if power > 100000 && power < 110000 && info.DownSpeed > 20.0 {
		data.Ylabels = "Lower Node"
	} else if power >= 110000 && power < 180000 && info.DownSpeed > 50.0 {
		data.Ylabels = "Medium Node"
	} else if power >= 180000 && info.DownSpeed > 100.0 {
		data.Ylabels = "Lightening Node"
	} else {
		data.Ylabels = "Useless"
	}
	dataArray := []dataset{}
	fmt.Println(dataArray)
	dataArray = append(dataArray, data)
	fmt.Println("Printing Data 1", dataArray)

	infoArray := []systemInfo{}
	infoArray = append(infoArray, info)

	file, _ := json.Marshal(infoArray)
	ioutil.WriteFile("extractedData.json", file, 0777)

	fmt.Println("Info Array Data", infoArray)
	newFile, _ := json.Marshal(dataArray)

	_ = ioutil.WriteFile("dataset.json", newFile, 0777)

	reqBody, err := json.Marshal(dataArray)
	if err != nil {
		fmt.Println("Error in Marshelling JSON Body.", err.Error())
		return
	}
	resp, err := http.Post("http://localhost:4400/api/dev/generateCSVfromJSON", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error in Sending Post Request", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in Reading Body ", err.Error())
		return
	}
	fmt.Println(string(body))
	fmt.Println("Computer Power", computerPower(float64(info.RAM)*0.000001, float64(info.CPUCores), info.UpSpeed))

}

func computerPower(
	ram float64,
	cores float64,
	speed float64) float64 {

	fmt.Println("Printing Ram ", ram)
	part1 := ram * 1.0
	part2 := 1.0
	part3 := cores * speed * 1.0

	x := (part1 + part2) * part3
	return math.Floor(x*100) / 100

}
