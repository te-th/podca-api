package podca_api

type Feed struct {
	Id 		int64		`json:"id" datastore:"-"`
	Title 		string 		`json:"title" xml:"title"`
	Link 		string 		`json:"link" xml:"link"`
	Description 	string		`json:"description" xml:"description" datastore:",noindex"`
	Language 	string		`json:"language" xml:"language"`
	Copyright	string		`json:"copyright" xml:"copyright"`
	PubDate		string		`json:"pubDate" xml:"pubDate"`
	Image 		Image		`json:"image" xml:"image"`
	Episodes	[]Episode	`json:"episodes" xml:"item"`
}


type Episode struct {
	//Id 		int64	`datastore:"-"`
	//sFeedId		int64
	Title 		string 	`json:"title" xml:"title"`
	Description 	string	`json:"description" xml:"description" datastore:",noindex"`
	Author		string	`json:"author" xml:"author"`
	Guid		string	`json:"guid" xml:"guid"`
	PubDate		string	`json:"pubDate" xml:"pubDate"`
}

type Image struct {
	//Id 	int64	`datastore:"-"`
	Url 	string	`json:"url" xml:"url"`
	Title 	string	`json:"title" xml:"title"`
	Link 	string	`json:"link" xml:"link"`
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