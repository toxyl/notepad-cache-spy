package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/toxyl/glog"
)

var (
	MAGIC_BYTES_UNSAVED = []byte{0x4E, 0x50, 0x00, 0x00, 0x01}
	MAGIC_BYTES_SAVED   = []byte{0x4E, 0x50, 0x00, 0x01}
	log                 = glog.NewLoggerSimple("npcs")
)

func init() {
	glog.LoggerConfig.ShowDateTime = false
	glog.LoggerConfig.ShowSubsystem = false
	glog.LoggerConfig.ShowRuntimeMilliseconds = false
	glog.LoggerConfig.SplitOnNewLine = true
	glog.LoggerConfig.CheckIfURLIsAlive = false
}

type Data struct {
	ID      []byte
	Unknown []byte
	Source  string
	File    string
	Content string
}

func checkError(err error) {
	if err != nil {
		log.Error("Failed: %s", glog.Error(err))
		os.Exit(1)
	}
}
func main() {
	var files []string
	if len(os.Args) < 2 {
		log.Warning("You didn't provide a filename, will process all files in the current users home dir.")
		log.Warning("Provide one or more filenames to process specific files instead.")
		files, _ = GetTabCacheFiles()
	} else {
		files = os.Args[1:]
	}

	for _, f := range files {
		data, err := DecodeBinaryFile(f)
		if err != nil {
			log.Warning("Could not decode %s: %s", glog.Auto(f), glog.Error(err))
			continue
		}
		log.Blank("")
		log.Blank(glog.Bold() + glog.Underline() + glog.WrapWhite(data.Source) + glog.Reset())
		log.Info("Type:    0x%X", data.ID)
		log.Info("File:    %s", data.File)
		log.Info("Unknown: %x", data.Unknown)
		log.Info("Content:")
		log.Blank("%s", data.Content)
	}
}

func DecodeBinaryFile(filename string) (data Data, err error) {
	reader := NewFileReader()
	reader.Open(filename)
	defer reader.Close()

	data.Source = filename
	data.ID = MAGIC_BYTES_SAVED

	isSavedFile := reader.LookAhead(MAGIC_BYTES_SAVED)
	if !isSavedFile {
		if !reader.LookAhead(MAGIC_BYTES_UNSAVED) {
			checkError(fmt.Errorf("this is not a Notepad file"))
		}
		data.ID = MAGIC_BYTES_UNSAVED
	}

	nSkip := 8 // usually we can skip 8 bytes ahead, unless...
	filepath := "<memory>"
	if isSavedFile {
		if strLen, err := reader.ReadInt(); err == nil {
			strLen *= 2
			filepath, _ = reader.ReadString(strLen)
		}
	} else {
		// there are several possible constructs in unsaved files
		// check the next two bytes to determine which we have to use
		t, err := reader.ReadBytes(2)
		checkError(err)

		if t[0] == t[1] { // both bytes we just read are identical
			if t[0] == 0 {
				nSkip += 5 // zero seems to be a special case
			} else {
				nSkip -= 3 // then the content starts 3 bytes earlier
			}
		}
	}

	data.File = filepath

	if isSavedFile {
		nSkip = 57
	}

	data.Unknown, err = reader.ReadBytes(nSkip)
	checkError(err)

	contentBuf := bytes.NewBuffer(nil)
	for {
		b, err := reader.ReadBytes(1)
		if err != nil {
			break
		}
		contentBuf.Write(b)
	}
	b := contentBuf.Bytes()
	if isSavedFile {
		b = b[0 : len(b)-14]
	} else {
		b = b[0 : len(b)-4]
	}
	data.Content = strings.TrimSpace(strings.ReplaceAll(string(b), "\r", "\n"))

	return data, nil
}
