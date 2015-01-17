package rego

import "encoding/json"
import "time"

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

// Thing endpoint represents the Reddit thing base class.
// https://github.com/reddit/reddit/wiki/JSON#thing-reddit-base-class
type Thing struct {
	Data json.RawMessage `json:"data"` // A data structure formatted based on kind
	ID   string          `json:"id"`   // Item identifier, e.g. "c3v7f8u"
	Kind string          `json:"kind"` // Kind denotes the item's type.
	Name string          `json:"name"` // Fullname of item, e.g. "t1_c3v7f8u"
}

// Created implements the Created class.
type Created struct {
	Local json.Number `json:"created"`     // Time of creation in local epoch-second format
	UTC   json.Number `json:"created_utc"` // Time of creation in UTC epoch-second format
}

// Time returns created time in local format
func (c *Created) Time() time.Time {
	return timeFromNumber(c.UTC)
}

// Votable implements the Votable class
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
	Saved            bool            `json:"saved"`                  // True if this post is saved by the logged in user
	Score            int             `json:"score"`                  // The net-score of the link
	Selfpost         bool            `json:"is_self"`                // True if this link is a selfpost
	SelftextHTML     string          `json:"selftext_html"`          //
	Selftext         string          `json:"selftext"`               //
	Stickied         bool            `json:"stickied"`               //
	SubredditID      string          `json:"subreddit_id"`           //
	Subreddit        string          `json:"subreddit"`              //
	Thumbnail        string          `json:"thumbnail"`              //
	Title            string          `json:"title"`                  //
	URL              string          `json:"url"`                    //
	Visited          bool            `json:"visited"`                //
	Created
	Votable
}

// Comment represents a subreddit post comment
type Comment struct {
	ApprovedBy       string          `json:"approved_by"`            // Who approved this comment, null if not a mod
	AuthorFlairClass string          `json:"author_flair_css_class"` // CSS class of the author's flair
	AuthorFlairText  string          `json:"author_flair_text"`      // Text of the author's flair
	Author           string          `json:"author"`                 // Account name of the poster
	BannedBy         string          `json:"banned_by"`              // Who removed this comment, null if not a mod
	BodyHTML         string          `json:"body_html"`              // Formatted HTML text as displayed on Reddit
	Body             string          `json:"body"`                   // Raw unformatted text of the comment
	Distinguished    string          `json:"distinguished"`          //
	Edited           json.RawMessage `json:"edited"`                 //
	Likes            bool            `json:"likes"`                  // How the logged-in user has voted on the link
	LinkAuthor       string          `json:"link_author"`            // Author of the parent link
	LinkID           string          `json:"id"`                     // ID of the link this comment is in
	LinkTitle        string          `json:"link_title"`             // Title of the parent link
	LinkURL          string          `json:"title_url"`              // Link URL of the parent link
	Name             string          `json:"name"`                   // Fullname of item, e.g. "t3_c3v7f8u"
	NumReports       int             `json:"num_reports"`            // Number of times comment has been reported, null if not a mod
	ParentID         string          `json:"parent_id"`              // ID of the thing this comment is a reply to
	Saved            bool            `json:"saved"`                  // True if this post is saved by the logged in user
	ScoreHidden      bool            `json:"score_hidden"`           // Whether the comment's score is currently hidden.
	Score            int             `json:"score"`                  // The net-score of the link
	SubredditID      string          `json:"subreddit_id"`           // ID of the subreddit
	Subreddit        string          `json:"subreddit"`              // Subreddit name
	Created
	Votable
}

// CommentResult is returned when submitting a new comment
type CommentResult struct {
	ID          string `json:"id"`          // UNKNOWN
	Name        string `json:"link"`        // Full name of item, e.g. "t3_c3v7f8u"
	ContentHTML string `json:"contentHTML"` // Comment text HTML formatted
	Content     string `json:"contentText"` // Comment text plain
	Replies     string `json:"replies"`     // UNKNOWN
	Parent      string `json:"parent"`      // Parent item
}
