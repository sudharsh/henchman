package henchman

import (
	"github.com/flosch/pongo2"
)

func init() {
	pongo2.RegisterFilter("success", filterSuccess)
	pongo2.RegisterFilter("ignored", filterIgnored)
	pongo2.RegisterFilter("reset", filterReset)
	pongo2.RegisterFilter("failure", filterFailure)
}

func filterSuccess(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(in.String() == "success"), nil
}

func filterIgnored(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(in.String() == "ignored"), nil
}

func filterReset(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(in.String() == "reset"), nil
}

func filterFailure(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(in.String() == "failure"), nil
}
