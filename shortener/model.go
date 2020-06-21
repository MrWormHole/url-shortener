package shortener

type Redirect struct {
	Hash      string `json:"hash" bson:"hash" msgpack:"hash"`
	URL       string `json:"url" bson:"url" msgpack:"url" validate:"empty=false & format=url"`
	CreatedAt int64  `json:"created_at" bson:"created_at" msgpack:"created_at"`
}
