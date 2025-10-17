package internal

import (
	"encoding/xml"
	"strings"
	"testing"

	"marktstammdatenregister.dev/internal/spec"
)

func TestXMLReaderAcceptsAliasRootAndElement(t *testing.T) {
	table := spec.Table{
		Root:        "MarktakteureUndRollen",
		RootAliases: []string{"Marktrollen"},
		Element:     "Marktrolle",
		ElementAliases: []string{
			"MarktakteurUndRolle",
		},
		Fields: []spec.Field{
			{Name: "MarktakteurMastrNummer"},
			{Name: "Marktrolle"},
		},
	}

	const doc = `<?xml version="1.0" encoding="UTF-8"?>
<MarktakteureUndRollen>
  <MarktakteurUndRolle>
    <MarktakteurMastrNummer>MASTR123</MarktakteurMastrNummer>
    <Marktrolle>VNB</Marktrolle>
  </MarktakteurUndRolle>
</MarktakteureUndRollen>`

	dec := xml.NewDecoder(strings.NewReader(doc))
	reader := NewXMLReader(&table, dec)

	item, err := reader.Read()
	if err != nil {
		t.Fatalf("Read() returned unexpected error: %v", err)
	}

	if got, want := item["MarktakteurMastrNummer"], "MASTR123"; got != want {
		t.Fatalf("MarktakteurMastrNummer = %q, want %q", got, want)
	}
	if got, want := item["Marktrolle"], "VNB"; got != want {
		t.Fatalf("Marktrolle = %q, want %q", got, want)
	}
}
