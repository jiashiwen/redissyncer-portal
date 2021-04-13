package commons

import (
	"io/ioutil"
	"time"
)

type Report struct {
	ReportContent map[string]interface{}
}

func (r Report) Json() (jsonResult string, err error) {
	bodyStr, err := json.MarshalIndent(r, "", " ")
	return string(bodyStr), err
}

func (r Report) JsonToFile() error {

	now := time.Now().Format("20060102150405000")
	filename := "report_" + now + ".json"
	bodyStr, err := json.MarshalIndent(r.ReportContent, "", " ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filename, bodyStr, 0666); err != nil {
		return err
	}
	return nil

}
