package main
import "strings"


type (
    PHolderElem struct {
        holder string
        value string
    }
    Placeholder struct {
        Map []PHolderElem
    }
)

func (phold *Placeholder) add(holder string, value string) {
    phold.Map = append(phold.Map, PHolderElem{holder:holder, value:value})
}

func (phold *Placeholder) parse(str string) string {
    for _,element := range phold.Map {
        str = strings.Replace(str, element.holder, element.value, -1)
    }
    return str
}