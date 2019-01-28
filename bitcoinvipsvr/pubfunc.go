package bitcoinvipsvr

import (
	"bytes"
	"encoding/binary"
)

/*
// CheckErr aaa
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// GetCurrentPath aaa
func GetCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	CheckErr(err)
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

// InitConfig InitConfig
func InitConfig(path string) map[string]string {
	fmt.Printf("%s", path)
	myMap := make(map[string]string)

	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	//创建一个输出流向该文件的缓冲流*Reader
	r := bufio.NewReader(f)
	for {
		//读取，返回[]byte 单行切片给b
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		//去除单行属性两端的空格
		s := strings.TrimSpace(string(b))
		//fmt.Println(s)

		//判断等号=在该行的位置
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		//取得等号左边的key值，判断是否为空
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}

		//取得等号右边的value值，判断是否为空
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		//这样就成功吧配置文件里的属性key=value对，成功载入到内存中c对象里
		myMap[key] = value
	}
	return myMap
}
*/
func readInt32(buf []byte, begpos int, total int) (bool, int32) {
	endpos := begpos + 4
	if endpos <= total {
		var lvalue32 int32
		readbuf := bytes.NewReader(buf[begpos:endpos])
		err := binary.Read(readbuf, binary.LittleEndian, &lvalue32)
		if nil == err {
			return true, lvalue32
		}
	}
	return false, 0
}

func readUInt32(buf []byte, begpos int, total int) (bool, uint32) {
	endpos := begpos + 4
	if endpos <= total {
		var lvalue32 uint32
		readbuf := bytes.NewReader(buf[begpos:endpos])
		err := binary.Read(readbuf, binary.LittleEndian, &lvalue32)
		if nil == err {
			return true, lvalue32
		}
	}
	return false, 0
}

func readInt8(buf []byte, begpos int, total int) (bool, int8) {
	endpos := begpos + 1
	if endpos <= total {
		var lvalue8 int8
		readbuf := bytes.NewReader(buf[begpos:endpos])
		err := binary.Read(readbuf, binary.LittleEndian, &lvalue8)
		if nil == err {
			return true, lvalue8
		}
	}
	return false, 0
}

func readInt64(buf []byte, begpos int, total int) (bool, int64) {
	endpos := begpos + 8
	if endpos <= total {
		var lvalue64 int64
		readbuf := bytes.NewReader(buf[begpos:endpos])
		err := binary.Read(readbuf, binary.BigEndian, &lvalue64)
		if nil == err {
			return true, lvalue64
		}
	}
	return false, 0
}

func readuint64(buf []byte, begpos int, total int) (bool, uint64) {
	endpos := begpos + 8
	if endpos <= total {
		var lvalue64 uint64
		readbuf := bytes.NewReader(buf[begpos:endpos])
		err := binary.Read(readbuf, binary.LittleEndian, &lvalue64)
		if nil == err {
			return true, lvalue64
		}
	}
	return false, 0
}

func readString(buf []byte, begpos int, maxpos int) (string, int, bool) {
	var endpos int
	endpos = begpos + 4
	res, strLen := readInt32(buf, begpos, maxpos)
	if res {
		endpos += int(strLen)
		if res && (endpos <= maxpos) {
			strText := string(buf[begpos+4 : endpos])
			return strText, endpos, true
		}
	}
	return "", 0, false
}
