package loxone_stats

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//
type Statistic struct {
	Uuid string
	Month      int
	Year       int
	Uri        string
	Statistics Statistics
}

//
type Statistics struct {
	Name       string           `xml:"Name,attr"`
	NumOutputs int              `xml:"NumOutputs,attr"`
	Outputs    string           `xml:"Outputs,attr"`
	Values     []StatisticValue `xml:"S"`
}

//
type StatisticValue struct {
	DateTime string  `xml:"T,attr"`
	Value    float64 `xml:"V,attr"`
}

// fetch the statistic values of a single statistic
// from the loxone miniserver
func (s *Statistic) Fetch(m *Miniserver) error {
	//
	var url = fmt.Sprintf("%s://%s/stats/%s", m.Protocol, m.Host, s.Uri)
	client := &http.Client{}

	//
	req, err := m.authenticatedRequest("GET", url)
	if err != nil {
		return err
	}

	//
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	//
	stats, err := ParseStatisticValues(res.Body)
	if err != nil {
		return err
	}

	// set statistic values
	s.Statistics = stats

	return nil
}

//
func ParseStatisticValues(r io.Reader) (s Statistics, err error) {
	// transform data
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return Statistics{}, err
	}

	// parse xml data
	var statistics Statistics
	err = xml.Unmarshal(data, &statistics)
	if err != nil {
		return Statistics{}, err
	}

	return statistics, nil
}
