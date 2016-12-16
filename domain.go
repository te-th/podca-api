package podca_api

type Feed struct {
	Id 		int64		`datastore:"-"`
	Title 		string 		`xml:"title"`
	Link 		string 		`xml:"link"`
	Description 	string		`xml:"description" datastore:",noindex"`
	Language 	string		`xml:"language"`
	Copyright	string		`xml:"copyright"`
	PubDate		string		`xml:"pubDate"`
	Image 		Image		`xml:"image"`
	Episodes	[]Episode	`xml:"item"`
}


type Episode struct {
	//Id 		int64	`datastore:"-"`
	//sFeedId		int64
	Title 		string 	`xml:"title"`
	Description 	string	`xml:"description" datastore:",noindex"`
	Author		string	`xml:"author"`
	Guid		string	`xml:"guid"`
	PubDate		string	`xml:"pubDate"`
}

type Image struct {
	//Id 	int64	`datastore:"-"`
	Url 	string	`xml:"url"`
	Title 	string	`xml:"title"`
	Link 	string	`xml:"link"`
}


// Podcast struct is strong coupled to the Apple iTunes format
type Podcast struct{
	Id   		int64  		`json:"id"`
	ArtistName 	string 		`json:"artistName"`
	CollectionName 	string		`json:"collectionName"`
	FeedUrl 	string		`json:"feedUrl"`
	CollectionId 	int64		`json:"collectionId"`
	TrackId 	int64		`json:"trackId"`
	Genres 		[]string	`json:"genres"`
}