package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Parser interface {
	Parse(file *os.File) error
	String() string
	AddRecord(record string) error
}

type DockerfileParser struct {
	version string
}

func (d *DockerfileParser) Parse(file *os.File) error {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "ENV TAG") {
			d.version = strings.Split(text, ":")[1]
			return nil
		}
	}
	return fmt.Errorf("version not found")
}

func (d *DockerfileParser) AddRecord(record string) error {
	d.version = record
	return nil
}

func (d *DockerfileParser) String() string {
	return d.version
}

type VersionParser struct {
	datas         []*VersionStruct
	latestVersion string
}

type VersionStruct struct {
	Header  string
	Author  []string
	Context []string
}

func NewVersionParser(latestVersion string) *VersionParser {
	return &VersionParser{
		latestVersion: latestVersion,
	}
}

func (d *VersionParser) validate(lastVersion string) error {
	if lastVersion == d.latestVersion {
		return errors.New("no change committed")
	}
	nowYear := time.Now().Format("06")
	recordYear := strings.Split(d.latestVersion, "_")[1][:2]
	if nowYear != recordYear {
		return fmt.Errorf("the year provided is out of date. year: %s", recordYear)
	}
	if lastYear := strings.Split(lastVersion, "_")[1][:2]; lastYear != recordYear {
		lastTag := strings.Split(lastVersion, ".")[1]
		latestTag := strings.Split(d.latestVersion, ".")[1]
		if lastTag == latestTag {
			return fmt.Errorf("tag is same in different years. tag: %s", lastTag)
		}
	}
	return nil
}

func (d *VersionParser) Parse(file *os.File) error {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "###") {
			data := &VersionStruct{
				Header: strings.Split(text, " ")[1],
			}
			scanner.Scan()
			data.Author = strings.Split(scanner.Text(), " ")[2:]
			for scanner.Scan() {
				if text := scanner.Text(); text == "" {
					break
				} else {
					data.Context = append(data.Context, strings.Split(text, " ")[1])
				}
			}
			d.datas = append(d.datas, data)
		}
	}
	return d.validate(d.datas[len(d.datas)-1].Header)
}

func (d *VersionParser) AddRecord(record string) error {
	records := strings.Split(record, " ")
	d.datas[len(d.datas)-1].Context = append(d.datas[len(d.datas)-1].Context, records...)
	return nil
}

func (d *VersionParser) String() string {
	buffer := bytes.NewBuffer([]byte{})
	for i, data := range d.datas {
		if i == len(d.datas)-1 {
			fmt.Fprintf(buffer, "### %s\n", d.latestVersion)
			fmt.Fprintf(buffer, "+ Author %s %s\n", data.Author[0], time.Now().Format("2006.01.02"))
		} else {
			fmt.Fprintf(buffer, "### %s\n", data.Header)
			fmt.Fprintf(buffer, "+ Author %s\n", strings.Join(data.Author, " "))

		}
		for _, context := range data.Context {
			fmt.Fprintf(buffer, "+ %s\n", context)
		}
		fmt.Fprintln(buffer)
	}
	return buffer.String()
}
