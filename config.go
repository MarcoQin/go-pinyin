package pinyin

import (
	"bufio"
	"fmt"
	"github.com/yanyiwu/gojieba"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var (
	// 结巴分词
	jieba      *gojieba.Jieba
	PinyinDict map[int]string
	PhraseDict map[string]string
)

func init() {
	dictDir := path.Join(path.Dir(getCurrentFilePath()), "dict")
	pinyinDictPath := path.Join(dictDir, "pinyin_dict")
	phraseDictPath := path.Join(dictDir, "phrase_dict")
	PinyinDict = make(map[int]string)
	PhraseDict = make(map[string]string)
	jieba = gojieba.NewJieba()
	parsePinyinDict(pinyinDictPath)
	parsePhraseDict(phraseDictPath)
}

func getCurrentFilePath() string {
	_, filePath, _, _ := runtime.Caller(1)
	return filePath
}

func parsePinyinDict(dictFile string) {
	fd, err := os.Open(dictFile)
	if err != nil {
		fmt.Printf("open file %s error", dictFile)
		panic(err)
	}
	defer fd.Close()
	rd := bufio.NewReader(fd)
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		// line: `U+4E2D: zhōng,zhòng  # 中`
		dataSlice := strings.Split(line, "  #")
		dataSlice = strings.Split(dataSlice[0], ": ")
		// 0x4E2D
		//hexCode := strings.Replace(dataSlice[0], "U+", "0x", 1)
		hexCode := strings.Replace(dataSlice[0], "U+", "", 1)
		i, err := strconv.ParseUint(hexCode, 16, 64)
		if err != nil {
			panic(err)
		}
		// zhōng,zhòng
		pinyin := dataSlice[1]
		PinyinDict[int(i)] = pinyin
	}
}

func parsePhraseDict(dictFile string) {
	fd, err := os.Open(dictFile)
	if err != nil {
		fmt.Printf("open file %s error", dictFile)
		panic(err)
	}
	defer fd.Close()
	rd := bufio.NewReader(fd)
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		// line: `街号巷哭: jiē hào xiàng kū`
		dataSlice := strings.Split(line, ": ")
		// 街号巷哭
		key := dataSlice[0]
		// zhōng,zhòng
		pinyin := strings.Trim(dataSlice[1], "\n")
		PhraseDict[key] = pinyin
	}
}
