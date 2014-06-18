// testXML project main.go
package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

type (
	XMLload struct {
		XMLName       xml.Name          `xml:"load"`
		Duration      int               `xml:"duration,attr"`
		Unit          string            `xml:"unit,attr"`
		ArrivalPhases []XMLarrivalphase `xml:"arrivalphase"`
	}
	XMLarrivalphase struct {
		XMLName  xml.Name `xml:"arrivalphase"`
		Phase    int      `xml:"phase,attr"`
		Duration int      `xml:"duration,attr"`
		Unit     string   `xml:"unit,attr"`
		Users    XMLusers `xml:"users"`
	}
	XMLusers struct {
		XMLName     xml.Name `xml:"users"`
		Maxnumber   int      `xml:"maxnumber,attr"`
		Arrivalrate int      `xml:"arrivalrate,attr"`
		Unit        string   `xml:"unit,attr"`
	}
)

func setDuration(content string, durationUnit string, durationValue int, usersMaxNumber int, usersArrivalRate int, usersUnit string) string {
	var buffer bytes.Buffer
	inputReader := strings.NewReader(content)
	decoder := xml.NewDecoder(inputReader)
	encoder := xml.NewEncoder(&buffer)
	encoder.Indent("", " ")
	buffer.WriteString(xml.Header)
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		//Inspect element
		switch token := t.(type) {
		case xml.StartElement:
			if token.Name.Local == "load" {
				fmt.Printf("XXXX %#v\n", token)
				var load XMLload
				decoder.DecodeElement(&load, &token)
				load.Duration = durationValue
				load.Unit = durationUnit
				for i, _ := range load.ArrivalPhases {
					load.ArrivalPhases[i].Users.Arrivalrate = usersArrivalRate
					load.ArrivalPhases[i].Users.Maxnumber = usersMaxNumber
					load.ArrivalPhases[i].Users.Unit = usersUnit
				}
				encoder.Encode(load)

			} else {
				err := encoder.EncodeToken(t)
				if err != nil {
					fmt.Printf("error=%s\b", err.Error())
				}
			}
		case xml.EndElement:
			err := encoder.EncodeToken(t)
			if err != nil {
				fmt.Printf("error=%s\b", err.Error())
			}
		}
	}
	encoder.Flush()
	return buffer.String()
}
func main() {
	str := `<?xml version="1.0"?>
<!DOCTYPE tsung SYSTEM "/usr/share/tsung/tsung-1.0.dtd">
<tsung loglevel="notice" version="1.0">
    <clients>
        <client host="tsung1" weight="1" maxusers="64000" cpu="12">
        </client>
    </clients>

    <servers>
        <server host="50.97.233.134" port="8080" type="tcp"/>
    </servers>

    <load duration="5" unit="minute">
        <arrivalphase phase="1" duration="10" unit="second">
            <users maxnumber="100" arrivalrate="500" unit="second"/>
        </arrivalphase>
    </load>

    <sessions>
        <session probability="100" name="get" type="ts_http">
            <for from="1" to="2000" var="i">
                <request>
                    <http url="/" method="GET" version="1.1" content_type="application/json"/>
                </request>

                <thinktime value="5"></thinktime>
            </for>
        </session>
    </sessions>
</tsung>`

	fmt.Println(setDuration(str, "hours", 9999, 100, 10, "second"))

}
