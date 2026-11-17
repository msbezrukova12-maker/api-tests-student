package advertisementSearch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/utils"
	"api-tests-template/tests/fixtures"
)

func TestSearchByTitle(t *testing.T) {
	user := fixtures.Login(t)

	request := fixtures.AdvertisementRequest(t)
	request.Fields["title"] = "SearchStart " + utils.RandomString(8) + " SearchMiddle SearchEnd"
	created := advertisements.CreateAdvertisement(t, user.Token, request)

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
			found := assertions.FindAdvertisementByID(t, list.Items, created.ID)
			require.Equal(t, request.Fields["title"], found.Title)
		})
	}
}
