func modifyXmlValues(content string, load XMLload, cloud Clients) io.Reader {
	var buffer bytes.Buffer
	inputReader := strings.NewReader(content)
	decoder := xml.NewDecoder(inputReader)
	encoder := xml.NewEncoder(&buffer)
	encoder.Indent("", " ")
	buffer.WriteString(xml.Header)
	buffer.WriteString("<!DOCTYPE tsung SYSTEM '/usr/share/tsung/tsung-1.0.dtd'>\n")
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch token := t.(type) {
		case xml.StartElement:
			switch token.Name.Local {
			case "load":
				encoder.Encode(load)
			case "clients":
				encoder.Encode(cloud)
			default:
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
	return bytes.NewReader(buffer.Bytes())
}
