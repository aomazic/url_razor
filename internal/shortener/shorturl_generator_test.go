package shortener

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const UserId = "e0dba740-fc4b-4977-872c-d360239e6b1a"

func TestShortLinkGenerator(t *testing.T) {
	initialLink1 := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortLink1 := GenerateShortLink(initialLink1, UserId)

	initialLink2 := "https://www.eddywm.com/lets-build-a-url-shortener-in-go-with-redis-part-2-storage-layer/"
	shortLink2 := GenerateShortLink(initialLink2, UserId)

	initialLink3 := "https://spectrum.ieee.org/automaton/robotics/home-robots/hello-robots-stretch-mobile-manipulator"
	shortLink3 := GenerateShortLink(initialLink3, UserId)

	assert.Equal(t, shortLink1, "jTa4L57P")
	assert.Equal(t, shortLink2, "d66yfx7N")
	assert.Equal(t, shortLink3, "dhZTayYQ")
}

func FuzzGenerateShortLink(f *testing.F) {
	testcases := []struct {
		initialLink string
		userId      string
	}{
		{"https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html", UserId},
		{"https://www.eddywm.com/lets-build-a-url-shortener-in-go-with-redis-part-2-storage-layer/", UserId},
		{"https://spectrum.ieee.org/automaton/robotics/home-robots/hello-robots-stretch-mobile-manipulator", UserId},
		{"", ""}, // Empty strings
		{"http://a.b", "testuser"},
	}
	for _, tc := range testcases {
		f.Add(tc.initialLink, tc.userId)
	}
	f.Fuzz(func(t *testing.T, initialLink string, userId string) {
		shortLink := GenerateShortLink(initialLink, userId)

		if len(shortLink) != 8 {
			t.Errorf("GenerateShortLink(%q, %q) = %q, want length 8, got %d", initialLink, userId, shortLink, len(shortLink))
		}
	})
}
