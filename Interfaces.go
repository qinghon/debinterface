package debinterface

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Interfaces struct {
	backupPath string
	FilePath   string
	//Adapters        []Interface
	SourceDirectory string
	Source          string
}

func (f *Interfaces) Add(Adapter Interface) error {

	//f.Adapters = append(f.Adapters, Adapter)
	reader := NewReader("")
	err := reader.Read()
	if err == nil {
		for _, i := range reader.Adapters {
			if i.GetName() == Adapter.GetName() {
				return nil
			}
		}
	}
	if f.Source != "" || f.SourceDirectory != "" {
		if f.Source != "" {
			dir := filepath.Dir(f.Source)
			err := ioutil.WriteFile(path.Join(dir, Adapter.GetName()), []byte(Adapter.Export()), 0644)
			return err
		}
		if f.SourceDirectory != "" {
			dir := filepath.Dir(f.SourceDirectory)
			err := ioutil.WriteFile(path.Join(dir, Adapter.GetName()), []byte(Adapter.Export()), 0644)
			return err
		}
	}
	ff, err := os.OpenFile(f.FilePath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer ff.Close()
	_, err = ff.WriteString(Adapter.Export())
	return err
}
func (f *Interfaces) GetBlock(fi *os.File, Adapter Interface) (int, int, []int, error) {
	var lineNumber = 0

	br := bufio.NewReader(fi)
	var lineStr = ""
	var blockBegin, blockEnd, options int
	var delLines []int
	blockEnd = -1
	blockBegin = -1
	for {
		lineByte, _, e := br.ReadLine()
		lineStr = string([]rune(string(lineByte)))
		if e == io.EOF {
			break
		}
		lineNumber += 1
		if strings.TrimSpace(lineStr) == "" {
			continue
		}
		if strings.TrimSpace(lineStr)[0:1] == "#" {
			continue
		}
		if strings.Index(lineStr, "#") != -1 {
			lineStr = lineStr[:strings.Index(lineStr, "#")]
		}
		if strings.TrimSpace(lineStr) == "" {
			continue
		}
		sline := strings.Fields(lineStr)
		if sline[0] == "iface" {
			if sline[1] == Adapter.GetName() {
				blockBegin = lineNumber
			}
			if options > blockBegin {
				blockEnd = options
			}
			continue
		}
		if sline[0] == "auto" || sline[0] == "hotplug" {
			if sline[1] == Adapter.GetName() {
				delLines = append(delLines, lineNumber)
			}

			continue
		}
		options = lineNumber
	}
	//lineNumber=0
	return blockBegin, blockEnd, delLines, nil
}
func (f *Interfaces) Del(Adapter Interface) error {
	target, _ := Mktmp(false)
	defer os.Remove(target)
	fo, err := os.OpenFile(target, os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fo.Close()

	fi, err := os.Open(f.FilePath)
	if err != nil {
		return err
	}
	defer fi.Close()

	log.Println(target)
	br := bufio.NewReader(fi)
	//f.GetBlock(fi,fo,Adapter)
	blockBegin, blockEnd, delLines, err := f.GetBlock(fi, Adapter)
	if blockBegin < 0 {
		return os.ErrNotExist
	}
	log.Println(blockBegin, blockEnd, delLines)
	fi.Seek(0, 0)
	var skip bool
	var lineNumber = 0
	for {
		lineByte, _, e := br.ReadLine()
		//lineStr=string([]rune(string(lineByte)))
		if e == io.EOF {
			break
		}
		lineNumber += 1
		skip = false
		for i := 0; i < len(delLines); i++ {
			if lineNumber == delLines[i] {
				//log.Println("skip ",lineNumber,lineNumber==delLines[i])
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		if lineNumber < blockBegin || lineNumber > blockEnd {
			_, err := fo.Write(lineByte)
			//log.Println(string(lineByte))
			if err != nil {
				return err
			}
			_, err = fo.Write([]byte("\n"))
			continue
		}
		//_,err:=fo.Write([]byte(Adapter.Export()))
		//if err != nil {
		//	return err
		//}
	}
	//fo.Sync()
	fo.Seek(0, 0)
	//fi.Seek(0,0)
	fi.Close()
	fi, err = os.OpenFile(f.FilePath, os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = io.Copy(fi, fo)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (f *Interfaces) Update(Adapter Interface) error {
	target, _ := Mktmp(false)
	defer os.Remove(target)
	fo, err := os.OpenFile(target, os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fo.Close()

	fi, err := os.Open(f.FilePath)
	if err != nil {
		return err
	}
	defer fi.Close()

	log.Println(target)
	br := bufio.NewReader(fi)
	//f.GetBlock(fi,fo,Adapter)
	blockBegin, blockEnd, delLines, err := f.GetBlock(fi, Adapter)
	if blockBegin < 0 {
		return os.ErrNotExist
	}
	log.Println(blockBegin, blockEnd, delLines)
	fi.Seek(0, 0)
	var skip, updated bool
	var lineNumber = 0
	for {
		lineByte, _, e := br.ReadLine()
		//lineStr=string([]rune(string(lineByte)))
		if e == io.EOF {
			break
		}
		lineNumber += 1
		skip = false
		for i := 0; i < len(delLines); i++ {
			if lineNumber == delLines[i] {
				//log.Println("skip ",lineNumber,lineNumber==delLines[i])
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		if lineNumber < blockBegin || lineNumber > blockEnd {
			_, err := fo.Write(lineByte)
			//log.Println(string(lineByte))
			if err != nil {
				return err
			}
			_, err = fo.Write([]byte("\n"))
			continue
		}
		if updated {
			continue
		}
		_, err := fo.Write([]byte(Adapter.Export()))
		updated = true
		if err != nil {
			return err
		}
	}
	//fo.Sync()
	fo.Seek(0, 0)
	//fi.Seek(0,0)
	fi.Close()
	fi, err = os.OpenFile(f.FilePath, os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = io.Copy(fi, fo)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (f *Interfaces) Write() {

}
func (f *Interfaces) RollBack() {

}
func Mktmp(directory bool) (string, error) {
	tmpDIR := "/tmp"
	if os.Getenv("TMPDIR") != "" {
		tmpDIR = os.Getenv("TMPDIR")
	}
	randStr := RandStringRunes(10)
	target := path.Join(tmpDIR, "tmp."+randStr)
	if directory {
		return target, os.MkdirAll(target, 0700)
	}
	f, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	defer f.Close()
	return target, err

}
func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
