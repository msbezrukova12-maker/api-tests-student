package advertisementSearch

import (
	"strings"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/utils"
	"api-tests-template/tests/fixtures"
)

func TestSearchByTitle(t *testing.T) {
	user := fixtures.Login(t)

	request := fixtures.AdvertisementRequest(t)
	request.Fields["title"] = "SearchStart " + utils.RandomString(8) + " SearchMiddle SearchEnd"
	advertisements.CreateAdvertisement(t, user.Token, request)

	cases := []struct {
		name   string
		search string
	}{
		{name: "beginning", search: "SearchStart"},
		{name: "middle", search: "SearchMiddle"},
		{name: "end", search: "SearchEnd"},
		{name: "lowercase", search: strings.ToLower("SearchMiddle")},
		{name: "uppercase", search: strings.ToUpper("SearchMiddle")},
		{name: "mixed_case", search: utils.MixCase("SearchMiddle")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			list := advertisements.GetAdvertisementsList(t, user.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(t, list.Items, request.Fields["title"])
		})
	}
}
