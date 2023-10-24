package truecaller

import "fmt"

type Response struct {
	Name string `json:"name"`
}

func (e *Response) ParseInformationMessage() string {
	header := "Information:"
	name := fmt.Sprintf("\nName: %s", e.Name)

	return fmt.Sprint(header, name)
}
