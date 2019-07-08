package main
//文件夹：实验报告+源文件
//源地址和目的地址修改
import (
	"math/rand"
	"time"
	"fmt"
	"os"
	"bufio"
	"io"
	"strings"
)

//保存帧数据的全局变量
var AllFrames [200]string
var Frames1 [100]string
var Frames2 [100]string

//保存转发表数据的全局变量
var SourceAddress [200]string//保存的源地址
var DestinationAddress [200]string//保存的目的地址

type NetBridge struct{
	Address [200]string//目的地址
	Port [200]string//对应的接口号
	Time [200]int//保存时间
}
var NetBridgeTable NetBridge

func main(){
	//WriteAllMAC()
	ReadAllMAC()
	FramesSlice()
	Forwarding()
}

//缓存区满删除最旧的数据
func DelOld(Count int){
	var i=0
	var TimOldest=NetBridgeTable.Time[i]
	if(Count>=5){
		for{
			if(i>=4){
				break
			}
			if(NetBridgeTable.Time[i+1]<TimOldest){
				TimOldest=NetBridgeTable.Time[i+1]
			}
			i++
		}
	}
	fmt.Println("缓存区满，删除数据"+NetBridgeTable.Address[TimOldest]+NetBridgeTable.Port[TimOldest])
	NetBridgeTable.Address[TimOldest] = ""
	NetBridgeTable.Port[TimOldest] = ""
	NetBridgeTable.Time[TimOldest] = 0
}

//找到转发表的空位
func FindSpace()int {
	var i= 0
	var Space = -1
	for {
		if (i >= 5) {
			break
		}
		if (NetBridgeTable.Address[i] == "") {
			Space = i
		}
		i++
	}
	if (Space==-1){
		DelOld(i)
		Space=FindSpace()
	}
	return Space
}


//时间减少1s
func TimeReduce(){
	var i=0
	for {
		if (i >= 5) {
			break
		}
		if (NetBridgeTable.Time[i] != 0) {
			NetBridgeTable.Time[i] -= 1
			if (NetBridgeTable.Time[i] == 1) {
				fmt.Println("定期清空数据："+NetBridgeTable.Address[i]+NetBridgeTable.Port[i])
				NetBridgeTable.Address[i] = ""
				NetBridgeTable.Port[i] = ""
				NetBridgeTable.Time[i] = 0
			}
		}
		i++
	}
}

//查找是否在同一网段
func CheckInOne(Source string,Destination string)int{
	var ok=-1;
	var i,j=0,0;
	for{
		if(SourceAddress[i]==""||DestinationAddress[j]==""||(Source==SourceAddress[i]&&Destination==SourceAddress[j])){
			break;
		}
		if(Source!=SourceAddress[i]){
			i++
		}
		if(Destination!=SourceAddress[j]){
			j++
		}
	}
	if((i<=6&&j<=6)||(i>6&&j>6)){
		ok=1
	}
	return ok;
}

//搜索转发表，找到相同地址返回i,找不到返回-1
func CheckTable(Address string)int{
	var i=0
	for {
		if (NetBridgeTable.Address[i] == "") {
			i = -1
			break
		}
		if (NetBridgeTable.Address[i] == Address) {
			break
		}
		i++
	}
	return i
}

//显示转发表内容
func ShowTable(){
	fmt.Println()
	fmt.Println("转发表内容")
	var i=0
	for{
		if(i>=5){
			break
		}
		if(NetBridgeTable.Address[i]!=""){
			fmt.Println(NetBridgeTable.Address[i]+"  "+NetBridgeTable.Port[i])
			//fmt.Println(NetBridgeTable.Time[i])
		}
		i++
	}
	fmt.Println()
}

//转发
func Forwarding(){
	var i=0
	var NetBridgeCount=0
	for {
		if (SourceAddress[i] == "" || DestinationAddress[i] == "") {
			break
		}
		if(i>=200){
			i=0
		}
		time.Sleep(1000000000)
		InOneNet := CheckInOne(SourceAddress[i], DestinationAddress[i])

		FindSourceMACCount := CheckTable(SourceAddress[i])
		FindDestinationMACCount := CheckTable(DestinationAddress[i])
		fmt.Println("目的地址:" + DestinationAddress[i]+"源地址:" + SourceAddress[i])

		NetBridgeCount=FindSpace()
		//广播地址
		if(DestinationAddress[i]=="FF-FF-FF-FF-FF-FF") {
			NetBridgeTable.Address[NetBridgeCount] = SourceAddress[i]
			if (FindSourceMACCount == -1) {
				NetBridgeTable.Address[NetBridgeCount] = SourceAddress[i]
				if (i%2 == 0) {
					NetBridgeTable.Port[NetBridgeCount] = "接口1"
					NetBridgeTable.Time[NetBridgeCount] = 9
				}
				if (i%2 == 1) {
					NetBridgeTable.Port[NetBridgeCount] = "接口2"
					NetBridgeTable.Time[NetBridgeCount] = 9
				}
			}
			if (i%2 == 0) {
				fmt.Println("转发接口2")
			}
			if (i%2 == 1) {
				fmt.Println("转发接口1")
			}
		} else if (InOneNet == 1) {//在同一网段
			//源地址不在转发表-》记录
			if (FindSourceMACCount == -1) {
				NetBridgeTable.Address[NetBridgeCount] = SourceAddress[i]
				if (i%2 == 0) {
					NetBridgeTable.Port[NetBridgeCount] = "接口1"
					NetBridgeTable.Time[NetBridgeCount]=9
				}
				if (i%2 == 1) {
					NetBridgeTable.Port[NetBridgeCount] = "接口2"
					NetBridgeTable.Time[NetBridgeCount]=9
				}
				//源地址不在转发表+目的地址不在转发表
				if (FindDestinationMACCount == -1) {
					if (i%2 == 0) {
						fmt.Println("转发接口2")
					}
					if (i%2 == 1) {
						fmt.Println("转发接口1")
					}
				}else{//源地址不在转发表+目的地址在转发表
				fmt.Println("不转发")
				}
			}else{//源地址在转发表
			//源地址在转发表+目的地址不在转发表
				if (FindDestinationMACCount == -1) {
					if (i%2 == 0) {
						fmt.Println("转发接口2")
					}
					if (i%2 == 1) {
						fmt.Println("转发接口1")
					}
				}else{//源地址在转发表+目的地址在转发表
					fmt.Println("不转发")
				}

			}
		}else{//不在同一网段
			//源地址不在转发表
			if (FindSourceMACCount == -1) {
				NetBridgeTable.Address[NetBridgeCount] = SourceAddress[i]
				if (i%2 == 0) {
					NetBridgeTable.Port[NetBridgeCount] = "接口1"
					NetBridgeTable.Time[NetBridgeCount]=9
				}
				if (i%2 == 1) {
					NetBridgeTable.Port[NetBridgeCount] = "接口2"
					NetBridgeTable.Time[NetBridgeCount]=9
				}
				//源地址不在转发表+目的地址不在转发表
				if (FindDestinationMACCount == -1) {
					if (i%2 == 0) {
						fmt.Println("转发接口2")
					}
					if (i%2 == 1) {
						fmt.Println("转发接口1")
					}
				} else {//源地址不在转发表+目的地址在转发表
					if (i%2 == 0) {
						fmt.Println("转发接口2")
					}
					if (i%2 == 1) {
						fmt.Println("转发接口1")
					}
				}
			} else {//源地址在转发表
				//源地址在转发表+目的地址不在转发表
				if (FindDestinationMACCount == -1) {
					if (i%2 == 0) {
						fmt.Println("转发接口2")
					}
					if (i%2 == 1) {
						fmt.Println("转发接口1")
					}
				} else {//源地址在转发表+目的地址在转发表
					if (i%2 == 0) {
						fmt.Println("转发接口2")
					}
					if (i%2 == 1) {
						fmt.Println("转发接口1")
					}
				}
			}
		}
		i++
		TimeReduce()
		ShowTable()
	}

}

//把帧数据切割放到两个全局变量中
func FramesSlice(){
	var i=0
	for{
		if (AllFrames[i]=="") {
			break
		}
		FramesSplit:= strings.Split(AllFrames[i],"\t")
		DestinationAddress[i]=FramesSplit[0]
		SourceAddress[i]=FramesSplit[1]
		i++
	}
}

//把两个帧的数据存到一个帧中
func FramsSum(){
	var i,j,k=0,0,0
	for{
		if(Frames1[j]==""&&Frames2[k]==""){
			break
		}
		if(Frames1[j]!=""){
			AllFrames[i]=Frames1[j]
			i++;j++
		}
		if(Frames2[k]!="") {
			AllFrames[i] = Frames2[k]
			i++;k++
		}
	}
}

//读取所有MAC文件的地址
func ReadAllMAC(){
	fmt.Println("接口1的帧数据")
	ReadOneMac("接口1.txt")
	fmt.Println("接口2的帧数据")
	ReadOneMac("接口2.txt")
	FramsSum()
}

//按行读取一个文件的MAC地址并存入全局变量
func ReadOneMac(filename string)error{
	ok,err:=os.OpenFile(filename,os.O_RDONLY,0)
	var Count=0
	if err != nil {
		fmt.Println("ReadOneMac os.Openfile failed. err: " + err.Error())
	} else {
		rd:=bufio.NewReader(ok)
		for{
			line,err:=rd.ReadString('\n')
			line = strings.TrimSpace(line)//去空白，即去'\n'
			if(filename=="接口1.txt"){
				Frames1[Count]=line
			}
			if(filename=="接口2.txt"){
				Frames2[Count]=line
			}
			fmt.Println(line)
			Count+=1
			if err != nil || io.EOF == err {
				break
			}
		}
	}
	defer ok.Close()
	return err
}

//把六个随机MAC分别写入两个文件
func WriteAllMAC(){
	var i=0
for {
	WriteOneMAC("接口1.txt")
	WriteOneMAC("接口2.txt")
	i++
	if(i>=12){
		break;
	}
}
}

//把一个随机MAC地址写入文件
func WriteOneMAC(filename string)error{
	MACAddressOne:=RandMAC()
	ok,err:=os.OpenFile(filename,os.O_WRONLY,0644)
	if err != nil {
		fmt.Println("WriteOneMac os.Openfile failed. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := ok.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		if(n==0){
			_, err = ok.WriteAt([]byte(MACAddressOne), n)
		}else if(n%35<6) {
		_, err = ok.WriteAt([]byte("\n"+MACAddressOne), n)
		} else{
			_, err = ok.WriteAt([]byte("\t"+MACAddressOne), n)}
	}
	defer ok.Close()
	return err
}

//随机产生MAC地址
func RandMAC()string{
	var MACAddress string
	//!!!byte rune string的关系与转换
	var RandOne rune//rune类型约等于char，但是是int32
	rand.Seed(time.Now().UnixNano())//UnixNano() better than Unix()
	for {
		RandInt := rand.Intn(22)
		if RandInt>9&&RandInt<17{
			continue;
		}else {RandInt+=48}
		RandOne = rune(RandInt)
		RandString:=string(RandOne)
		MACAddress+=RandString
		if(len(MACAddress)>=17) {
			fmt.Println(MACAddress)
			break;
		}
		if((len(MACAddress)%3)==2){
			MACAddress+="-"
		}

	}
return MACAddress
}
