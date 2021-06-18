package birdblog_test

import (
	"birdblog"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGivenJSONDataShouldDecodeIntoStruct(t *testing.T) {
	file, err := os.Open("testdata/conversation.json")
	if err != nil {
		t.Fatal(err)
	}
	var got birdblog.Conversation
	want := birdblog.Conversation{
		"This is super spoiler free, as IMO the way you design games based on other media is to extract patterns and not to translate content.\n\nAlso: I've got thoughts on what I'd use for system, but I'd REALLY rather talk about what a system should be doing to emulate the feel.",
		"...feel free to apply whatever we are talking about to your system of choice, but I'm not shopping or asking for system recommendations. I  probably won't interact much with \"this system does X\" even though I'm sure it does! \n\nThis is just not that conversation.",
		"Last disclaimer: I'm probably not actually building this thing, which is why I am talking about it...it's my way for getting things out of my head so I can focus. If my thoughts make sense and you wanna build something using it, go for it!\n\nOK, ready?",
		"I will take some tips from the video game, but the *show* Castlevania is just so perfect for a TTRPG it is the main emphasis of what I'm talking about here.\n\nSo, my favorite favorite thing about the show is how well it uses dialogue to create story and transition to action...",
		"but more importantly, it uses its dialogue in such a way that it makes the dialogue more like an action scene than a dramatic scene. The level of interplay and back and forth, the way the dialogue advances plot, builds relationships...\nThe dialogue *is* the action.",
		"and the \"action\" scenes are the transitional payoffs between the action of dialogue. \nI really loved the action scenes, but I found myself wanting the dialogue scenes more and more. The action scenes were ways to release the tension of the discoveries and revelations in dialogue.",
		"rather than dialogue feeling like this compulsory space between fight scenes, fight scenes felt like the natural consequence of what the characters were talking and coming into conflict about.\n\nFor me, the role of talking and fight scenes were completely reversed!",
		"This is a show that routinely had multiple episodes where people just...talked....and it felt like almost a let down when a fight broke out!\n\nSo, anyway, from a design point of view, I look for patterns like this, how those patterns make me feel.\n\nThen I think about structure.",
		"It's fun to just jump to mechanics and delve into system, but IME you can do so much more heavy lifting with looking at the structure of how sessions will go. \n\nAlso, structure is portable. A good structure can work with just about any system you please.",
		"The guiding principle I'd use for a Castlevania structure is:\n\nTalk until it explodes.",
		"Which is really a way of saying that you structure the game in a way that most actions are part of the background of a scene and that we put the discussions between characters at the fore.",
		"example:\nCharacters are travelling from place A to B through a forest. We can make the skill checks and rolls of them navigating the forest, but instead of making the narrative about how they navigate the forest, we break those scenes into convos...",
		"\"What are we talking about as we clear through the brush?\"\n\n\"What are arguing about as we find our way back to the trail?\"\n\n\"What do we learn about each other as we enter the town?\"\n\nWe alter the focus of the camera to the words of the characters.",
		"and then...when it doesn't seem like there is any more to say, when there is a some insurmountable topic or conflict...the conversational momentum explodes into big physical action. We flip the camera's focus on the action now.",
		"You can focus on dialogue from a third person point of view and at distance  - \"he says..\" if you are more comfortable that way.\nOr you can delve full into scene chewing character talking if that is your thing. Either works!\nThe GM's role is to keep background moving...",
		"providing the backdrop details as characters talk...if they are helping repair a village, putting their convo at the forefront, but still injecting some details of the scene and having them make the checks and rolls they would normally make as they talk.",
		"The goal is to turn what you are doing as you talk into the subtext of the talk, and to allow that subtext to reify your conversation so that it feels substantial, occasionally resetting and palette-cleansing with physical action.",
		"And that's it, basically!  There are a lot of ways to use system to support this structure. I think you could definitely build a system just *for* the support of this structure, but you can honestly use this for just about any system you enjoy.\n\nHope it helps and makes sense!",
	}

	got, err = birdblog.DecodeJSONIntoStruct(file)
	// fmt.Printf("%#v", got)
	if err != nil {
		t.Fatal(err)
	}

	if !(cmp.Equal(want, got)) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGivenIDGenerateTweetRequest(t *testing.T) {
	var ID string = "1400096718788743180"
	var token string = "dummyToken"
	wantURL, err := url.Parse("https://api.twitter.com/2/tweets/search/recent?query=conversation_id:1400096718788743180&max_results=100&tweet.fields=in_reply_to_user_id,author_id,created_at,conversation_id")
	if err != nil {
		t.Fatal(err)
	}
	got, err := birdblog.NewConversationRequest(token, ID)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(http.MethodGet, got.Method) {
		t.Error(cmp.Diff(http.MethodGet, got.Method))
	}
	if !cmp.Equal(wantURL, got.URL) {
		t.Error(cmp.Diff(wantURL, got.URL))
	}
	wantHeader := "Bearer dummyToken"
	gotHeader := got.Header.Get("Authorization")
	if !cmp.Equal(wantHeader, gotHeader) {
		t.Errorf("bad Authorization header: %s", cmp.Diff(wantHeader, gotHeader))
	}
}

func TestGenerateGhostPostRequest(t *testing.T) {
	var content birdblog.Conversation
	content = append(content, "hello blog!")
	credentials := "someJWTToken"
	site := "myblog.ghost.io"
	gpr, err := birdblog.GhostPostRequest(content, credentials, site)
	wantBody := `
{
	"posts": [{
		"title": "Twitter Draft",
		"html": "hello blog!",
		"status": "draft"
	}]
}`
	gotBody, err := io.ReadAll(gpr.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantBody, string(gotBody)) {
		t.Errorf(cmp.Diff(wantBody, string(gotBody)))
	}

}

func TestCanCreateValidGhostJWT(t *testing.T) {
	now := time.Unix(162401290, 0)
	// dummy values to test that we get consistent results
	wantJWT := "eyJhbGciOiJIUzI1NiIsImN0eSI6IkpXVCIsImtpZCI6ImI0OWMxYzRmNjBlMDAzYjhhYTZjNmYwYyIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIvdjQvYWRtaW4vIiwiZXhwIjoxNjI0MDE1OTAsImlhdCI6MTYyNDAxMjkwfQ._ZJAG7d5_69HnGqvNWTbXPu_XsQtPQ8bwKVjCRua-oA"
	api := "b49c1c4f60e003b8aa6c6f0c:d396b0a338dc588295efe7d71e381cf6a82729d4f2783e47123803cbc0cca621"
	gotJWT, err := birdblog.MakeGhostJWT(api, now)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantJWT, gotJWT) {
		t.Errorf(cmp.Diff(wantJWT, gotJWT))
	}
}
