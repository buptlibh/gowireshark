package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/randolphcyg/gowireshark"
)

const inputFilepath = "../pcaps/s7comm_clean.pcap"
const inputFilepath2 = "../pcaps/wincc_s400_production.pcap"

func TestEpanVersion(t *testing.T) {
	fmt.Println(gowireshark.EpanVersion())
}

func TestEpanPluginsSupported(t *testing.T) {
	fmt.Println(gowireshark.EpanPluginsSupported())
}

func TestDissectPrintFirstFrame(t *testing.T) {
	err := gowireshark.DissectPrintFirstFrame(inputFilepath)
	if err != nil {
		fmt.Println(err)
	}

}

func TestDissectPrintAllFrame(t *testing.T) {
	err := gowireshark.DissectPrintAllFrame(inputFilepath)
	if err != nil {
		fmt.Println(err)
	}

}

func TestDissectPrintFirstSeveralFrame(t *testing.T) {
	err := gowireshark.DissectPrintFirstSeveralFrame(inputFilepath, 200)
	if err != nil {
		fmt.Println(err)
	}

}

func TestDissectPrintSpecificFrame(t *testing.T) {
	err := gowireshark.DissectPrintSpecificFrame(inputFilepath, 5000)
	if err != nil {
		fmt.Println(err)
	}
}

func TestCapFileMulSeq(t *testing.T) {
	var err error

	fmt.Println("@@@@@@@@@@@@@")
	err = gowireshark.DissectPrintSpecificFrame(inputFilepath, 5000)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("$$$$$$$$$$$$$")
	err = gowireshark.DissectPrintSpecificFrame(inputFilepath2, 50)
	if err != nil {
		fmt.Println(err)
	}
}

//func TestCapFileMulGoroutine(t *testing.T) {
//	var err error
//
//	go func() {
//		fmt.Println("@@@@@@@@@@@@@")
//		err = gowireshark.DissectPrintSpecificFrame(inputFilepath, 5000)
//		if err != nil {
//			fmt.Println(err)
//		}
//	}()
//
//	go func() {
//		fmt.Println("$$$$$$$$$$$$$")
//		err = gowireshark.DissectPrintSpecificFrame(inputFilepath2, 50)
//		if err != nil {
//			fmt.Println(err)
//		}
//	}()
//
//	time.Sleep(time.Second * 3)
//}

/*
RESULT: none
*/
func TestDissectPrintSpecificFrameOutOfBounds(t *testing.T) {
	err := gowireshark.DissectPrintSpecificFrame(inputFilepath, 5448)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetSpecificFrameHexData(t *testing.T) {
	res, err := gowireshark.GetSpecificFrameHexData(inputFilepath, 3)
	if err != nil {
		fmt.Println(err)
	}

	for i, item := range res.Offset {
		fmt.Println(i, item)
	}
	for i, item := range res.Hex {
		fmt.Println(i, item)
	}
	for i, item := range res.Ascii {
		fmt.Println(i, item)
	}
}

/*
	{
		"1": {
			"_index": "packets-20xx-0x-0x",
			"offset": ["0000", "0010", ...],
			"hex": ["00 1c 06 1c 69 e4 20 47 47 87 d4 96 08 00 45 00", ...],
			"ascii": ["....i. GG.....E.", ...],
			"_source": {
				"layers": {
	                "frame": {
						...
					},
					"eth": {
						"eth.dst": "00:1c:06:1c:69:e4",
						"eth.dst_tree": {
	                        ...
						}
					},
	                ...
				}
			}
		},
	    ...
	}
*/
func TestGetSpecificFrameProtoTreeInJson(t *testing.T) {
	specificFrameDissectRes, err := gowireshark.GetSpecificFrameProtoTreeInJson(inputFilepath, 3)
	if err != nil {
		fmt.Println(err)
	}

	for k, frameData := range specificFrameDissectRes {
		fmt.Println("# Frame num:" + k)
		fmt.Println("## WsIndex:", frameData.WsIndex)
		fmt.Println("## Offset:", frameData.Offset)
		fmt.Println("## Hex:", frameData.Hex)
		fmt.Println("## Ascii:", frameData.Ascii)
		fmt.Println("## Layers:")
		for layer, layerData := range frameData.WsSource.Layers {
			fmt.Println("### Layers num:", layer, layerData)
		}
	}

}

func TestGetAllFrameProtoTreeInJson(t *testing.T) {
	allFrameDissectRes, err := gowireshark.GetAllFrameProtoTreeInJson(inputFilepath)
	if err != nil {
		fmt.Println(err)
	}

	// Do not print the content in case the amount of data is too large
	fmt.Printf("frame count:%d\n", len(allFrameDissectRes))

	allFrameDissectRes2, err := gowireshark.GetAllFrameProtoTreeInJson(inputFilepath)
	if err != nil {
		fmt.Println(err)
	}

	// Do not print the content in case the amount of data is too large
	fmt.Printf("frame count:%d\n", len(allFrameDissectRes2))
}

/*
Get interface device list, TODO but not return addresses of device (cause c logic error)
Result:
bluetooth-monitor Bluetooth Linux Monitor 56
nflog Linux netfilter log (NFLOG) interface 48
nfqueue Linux netfilter queue (NFQUEUE) interface 48
dbus-system D-Bus system bus 48
dbus-session D-Bus session bus 48
enp0s5  22
any Pseudo-device that captures on all interfaces 54
lo  55
*/
func TestGetIfaceList(t *testing.T) {
	iFaces, err := gowireshark.GetIfaceList()
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range iFaces {
		fmt.Println(k, v.Description, v.Flags)
	}
}

/*
Get interface device nonblock status, default is false
Result:
device: enp0s5  nonblock status: false
*/
func TestGetIfaceNonblockStatus(t *testing.T) {
	ifaceName := "enp0s5"
	status, err := gowireshark.GetIfaceNonblockStatus(ifaceName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("device:", ifaceName, " nonblock status:", status)
}

func TestSetIfaceNonblockStatus(t *testing.T) {
	ifaceName := "enp0s5"
	status, err := gowireshark.SetIfaceNonblockStatus(ifaceName, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("device:", ifaceName, "after set, now nonblock status is:", status)
}

/*
Test infinite loop fetching and parsing packets.
Set the num parameter of the DissectPktLive function to -1
to process packets in an infinite loop.
*/
func TestDissectPktLiveInfinite(t *testing.T) {
	ifName := "enp0s5"
	pktNum := -1
	promisc := 1
	timeout := 20
	livePkgCount := 0

	// socket server: start socket server and wait data to come
	err := gowireshark.RunSock()
	if err != nil {
		fmt.Println(err)
	}

	// user read unmarshal data from go channel
	go func() {
		for {
			select {
			case pkg := <-gowireshark.PkgDetailLiveChan:
				livePkgCount++
				pkgByte, _ := json.Marshal(pkg)
				fmt.Printf("Processed pkg:【%d】pkg len:【%d】\n", livePkgCount, len(pkgByte))
				//fmt.Println(pkg)
			default:
			}
		}
	}()

	// start socket client, capture and dissect packet.
	err = gowireshark.DissectPktLive(ifName, pktNum, promisc, timeout)
	if err != nil {
		fmt.Println(err)
	}

	select {}

}

/*
Set the num parameter of the DissectPktLive function to a specific num like 20
to process packets in a limited loop.
*/
func TestDissectPktLiveSpecificNum(t *testing.T) {
	ifName := "enp0s5"
	pktNum := 20
	promisc := 1
	timeout := 20
	livePkgCount := 0

	// socket server: start socket server and wait data to come
	err := gowireshark.RunSock()
	if err != nil {
		fmt.Println(err)
	}

	// user read unmarshal data from go channel
	go func() {
		for {
			select {
			case pkg := <-gowireshark.PkgDetailLiveChan:
				livePkgCount++
				pkgByte, _ := json.Marshal(pkg)
				fmt.Printf("Processed pkg:【%d】pkg len:【%d】\n", livePkgCount, len(pkgByte))
				//fmt.Println(pkg)
			default:
			}
		}
	}()

	// start socket client, capture and dissect packet.
	err = gowireshark.DissectPktLive(ifName, pktNum, promisc, timeout)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second)

}

func TestStopDissectPktLive(t *testing.T) {
	ifName := "enp0s5"
	pktNum := -1
	promisc := 1
	timeout := 20
	livePkgCount := 0

	// socket server: start socket server and wait data to come
	err := gowireshark.RunSock()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("socket server started!")

	// user read unmarshal data from go channel
	go func() {
		for {
			select {
			case pkg := <-gowireshark.PkgDetailLiveChan:
				//time.Sleep(time.Millisecond * 200)
				livePkgCount++
				pkgByte, _ := json.Marshal(pkg)
				fmt.Printf("Processed pkg:【%d】pkg len:【%d】\n", livePkgCount, len(pkgByte))
			default:
			}
		}

	}()

	go func() {
		fmt.Println("Simulate manual stop real-time packet capture!")
		time.Sleep(time.Second * 2)
		err := gowireshark.StopDissectPktLive()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("############ stop capture successfully! ##############")
	}()

	fmt.Println("start c client, start capture function")
	// start socket client, capture and dissect packet.
	err = gowireshark.DissectPktLive(ifName, pktNum, promisc, timeout)
	if err != nil {
		fmt.Println(err)
	}

	select {}
}
