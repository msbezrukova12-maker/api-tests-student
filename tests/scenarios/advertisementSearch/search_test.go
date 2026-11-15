package advertisementSearch

import (
	"strings"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/utils"
	"api-tests-template/tests/fixtures"
)

// TestSearchByTitle проверяет полное вхождение, позицию подстроки и регистронезависимость поиска.
func TestSearchByTitle(t *testing.T) {
	user := fixtures.Login(t)

	t.Log("[Precondition] Создаём независимое объявление для проверки поиска")
	request := fixtures.AdvertisementRequest(t)
	request.Fields["title"] = "SearchStart " + utils.RandomString(8) + " SearchMiddle SearchEnd"
	created := advertisements.CreateAdvertisement(t, user.Token, request)

	cases := []struct {
		name   string
		search string
	}{
		{name: "full_title", search: request.Fields["title"]},
		{name: "beginning", search: "SearchStart"},
		{name: "middle", search: "SearchMiddle"},
		{name: "end", search: "SearchEnd"},
		{name: "lowercase", search: strings.ToLower("SearchMiddle")},
		{name: "uppercase", search: strings.ToUpper("SearchMiddle")},
		{name: "mixed_case", search: utils.MixCase("SearchMiddle")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("[Step] Ищем объявление по строке %q", tc.search)
			list := advertisements.GetAdvertisementsList(t, user.Token, tc.search)
			found := assertions.FindAdvertisementByID(t, list.Items, created.ID)
			if found.Title != request.Fields["title"] {
				t.Fatalf("Найдено объявление с неправильным title: got=%q want=%q",
					found.Title, request.Fields["title"])
			}
			t.Log("[Result] Созданное объявление найдено")
		})
	}
}
