package geddit

import "encoding/json"

// Thing types (kind)
const (
	TypeComment   = "t1"
	TypeAccount   = "t2"
	TypeLink      = "t3"
	TypeMessage   = "t4"
	TypeSubreddit = "t5"
	TypeAward     = "t6"
	TypePromo     = "t8" // Promo campain
	TypeListing   = "Listing"
)

// Listing endpoint represents paginated content and it's contents are
// exposed via the Paginator type
// https://github.com/reddit/reddit/wiki/JSON#listing
type listing struct {
	Data struct {
		After    string  `json:"after"`
		Before   string  `json:"before"`
		Children []Thing `json:"children"`
		Modhash  string  `json:"modhash"`
	} `json:"data"`
	Kind string `json:"kind"` // Always "Listing"
}

// Thing endpoint represents the Reddit thing base class.
// https://github.com/reddit/reddit/wiki/JSON#thing-reddit-base-class
type Thing struct {
	Data json.RawMessage `json:"data"` // A data structure formatted based on kind
	ID   string          `json:"id"`   // Item identifier, e.g. "c3v7f8u"
	Kind string          `json:"kind"` // Kind denotes the item's type.
	Name string          `json:"name"` // Fullname of item, e.g. "t1_c3v7f8u"
}

type Created struct {
	Local json.Number `json:"created"`     // Time of creation in local epoch-second format
	UTC   json.Number `json:"created_utc"` // Time of creation in UTC epoch-second format
}

type Votable struct {
	Downs int  `json:"downs"`           // Number of downvotes. (includes own)
	Likes bool `json:"likes,omitempty"` // True if thing is liked by the user
	Ups   int  `json:"ups"`             // Number of upvotes. (includes own)
}

// Account represents a Reddit user account
type Account struct {
	CommentKarma  int    `json:"comment_karma"` // User's comment karma
	GoldCredits   int    `json:"gold_creddits,omitempty"`
	HasMail       bool   `json:"has_mail"`     // User has unread mail?
	HasModMail    bool   `json:"has_mod_mail"` // User has unread mod mail?
	HideRobots    bool   `json:"hide_from_robots,omitempty"`
	ID            string `json:"id"`                 // ID of the account; prepend t2_ to get fullname
	IsFriend      bool   `json:"is_friend"`          // Logged-in user has this user set as a friend
	IsGold        bool   `json:"is_gold"`            // Reddit gold status
	IsMod         bool   `json:"is_mod"`             // This account moderates a subreddits
	LinkKarma     int    `json:"link_karma"`         // User's link karma
	Modhash       string `json:"modhash"`            // Current modhash, not present if not your account
	Name          string `json:"name"`               // The username of the account
	Over18        bool   `json:"over_18"`            // If this account is set to be over 18
	VerifiedEmail bool   `json:"has_verified_email"` // User has a verified email address
	Created
}

// Link represents a subreddit post link
type Link struct {
	AuthorFlairClass string          `json:"author_flair_css_class"` // CSS class of the author's flair
	AuthorFlairText  string          `json:"author_flair_text"`      // Text of the author's flair
	Author           string          `json:"author"`                 // Account name of the poster
	Clicked          bool            `json:"clicked"`                //
	Distinguished    string          `json:"distinguished"`          //
	Domain           string          `json:"domain"`                 // The domain of this link
	Edited           json.RawMessage `json:"edited"`                 //
	Hidden           bool            `json:"hidden"`                 // True if the post is hidden by user
	ID               string          `json:"id"`                     // Item identifier, e.g. "c3v7f8u"
	IsNsfw           bool            `json:"over_18"`                // True if the post is tagged as NSFW
	Likes            bool            `json:"likes"`                  // How the logged-in user has voted on the link
	LinkFlairClass   string          `json:"link_flair_css_class"`   //
	LinkFlairText    string          `json:"link_flair_text"`        //
	MediaEmbed       json.RawMessage `json:"media_embed"`            //
	Media            json.RawMessage `json:"media"`                  //
	Name             string          `json:"name"`                   // Fullname of item, e.g. "t3_c3v7f8u"
	NumComments      int             `json:"num_comments"`           //
	Permalink        string          `json:"permalink"`              // Relative URL of the permanent link for this link
	Saved            bool            `json:"saved"`                  //
	Score            int             `json:"score"`                  //
	Selfpost         bool            `json:"is_self"`                // True if this link is a selfpost
	SelftextHtml     string          `json:"selftext_html"`          //
	Selftext         string          `json:"selftext"`               //
	Stickied         bool            `json:"stickied"`               //
	SubredditID      string          `json:"subreddit_id"`           //
	Subreddit        string          `json:"subreddit"`              //
	Thumbnail        string          `json:"thumbnail"`              //
	Title            string          `json:"title"`                  //
	Url              string          `json:"url"`                    //
	Visited          bool            `json:"visited"`                //
	Created
}

type CommentResult struct {
	ID           string `json:"id"`          // UNKNOWN
	Name         string `json:"link"`        // Full name of item, e.g. "t3_c3v7f8u"
	ComntentHtml string `json:"contentHTML"` // Comment text HTML formatted
	Content      string `json:"contentText"` // Comment text plain
	Replies      string `json:"replies"`     // UNKNOWN
	Parent       string `json:"parent"`      // Parent item
}
