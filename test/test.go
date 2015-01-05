package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/tbuckley/whistler"
	"io/ioutil"
	"log"
	"math"
	"path"
	"regexp"
	"time"
)

var (
	filenamePattern = regexp.MustCompile("(.*?)_.*")

	negative = flag.Bool("negative", false, "Indicate that you are recording a negative example")
	testdir  = flag.String("testdir", "", "The directory where test data should be stored")
	name     = flag.String("name", "recording", "An optional name for the recording")
	file     = flag.String("file", "", "A specific filename to use")
)

func main() {
	flag.Parse()

	switch flag.Arg(0) {
	case "record":
		record()
	case "test":
		runRecordings()
	}
}

func getFileName() string {
	var sign string
	if *negative {
		sign = "neg"
	} else {
		sign = "pos"
	}
	date := time.Now().Format("15040501022006")
	filename := fmt.Sprintf("%s_%s_%s.json", sign, *name, date)
	return path.Join(*testdir, filename)
}

func filenameIsPositive(filename string) (bool, error) {
	match := filenamePattern.FindStringSubmatch(filename)
	if len(match) != 2 {
		return false, errors.New("incorrect filename pattern")
	}
	return match[1] == "pos", nil
}

type RecordState struct {
	points [][]whistler.SineWave
}

func (s *RecordState) Name() string {
	return "RECORD"
}
func (s *RecordState) Handle(point []whistler.SineWave) whistler.State {
	s.points = append(s.points, point)
	return nil
}

type RecordFactory struct {
	recorder *RecordState
}

func (f *RecordFactory) New() *whistler.Matcher {
	return &whistler.Matcher{
		StartState: f.recorder,
		MatchState: new(RecordState),
	}
}

func record() {
	whistler.Initialize()
	defer whistler.Terminate()

	whistle, err := whistler.New()
	if err != nil {
		panic(err)
	}
	defer whistle.Close()

	recorder := new(RecordState)
	matchChan := whistle.Add(&RecordFactory{recorder})
	go func() {
		for {
			<-matchChan
		}
	}()

	whistle.Listen()

	_, err = fmt.Scanln()
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(recorder.points)
	if err != nil {
		panic(err)
	}

	filename := getFileName()
	ioutil.WriteFile(filename, data, 0755)

	fmt.Printf("Wrote file %s\n", filename)
}

func playRecording() {

}

func runRecordings() {
	fails := 0
	passes := 0

	if *file == "" {
		buffer := new(bytes.Buffer)
		log.SetOutput(buffer)

		files, err := ioutil.ReadDir(*testdir)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			if testFile(file.Name()) {
				passes += 1
			} else {
				fails += 1
			}
		}
	} else {
		if testFile(*file) {
			passes += 1
		} else {
			fails += 1
		}
	}

	if fails > 0 {
		fmt.Printf("[FAIL] %d / %d passed\n", passes, fails+passes)
	} else {
		fmt.Printf("[PASS] %d / %d passed\n", passes, fails+passes)
	}
}

func testFile(filename string) bool {
	filepath := path.Join(*testdir, filename)
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	points := make([][]whistler.SineWave, 0)
	err = json.Unmarshal(data, &points)
	if err != nil {
		panic(err)
	}

	matcher := whistler.Kikee.New()
	matches := 0
	for _, point := range points {
		filteredPoints := make([]whistler.SineWave, 0)
		for _, wave := range point {
			if wave.Amplitude > 0.005 {
				filteredPoints = append(filteredPoints, wave)
			}
		}
		point = filteredPoints
		if len(point) > 0 {
			log.Printf("[WAVE] =======================")
			l := int(math.Min(2.0, float64(len(point))))
			for _, wave := range point[:l] {
				log.Printf("[WAVE] %s", wave.String())
			}
		}
		if len(point) > 2 {
			point = point[:2]
		}
		if matcher.Match(point) {
			matches += 1
		}
	}

	shouldMatch, err := filenameIsPositive(filename)
	if err != nil {
		panic(err)
	}

	if (matches > 0) == shouldMatch {
		fmt.Printf("[PASS] %s has %d matches\n", filename, matches)
		return true
	} else {
		fmt.Printf("[FAIL] %s has %d matches\n", filename, matches)
		return false
	}
}
