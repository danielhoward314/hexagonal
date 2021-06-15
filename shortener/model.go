package shortener

type Redirect struct {
	Code      string `json:"code" bson:"code" msgpack:"code" gorm:"primary_key;size:50;not null; unique"`
	URL       string `json:"url" bson:"url" msgpack:"url" gorm:"size:100;not null" validate:"empty=false & format=url"`
	CreatedAt int64  `json:"created_at" bson:"created_at" msgpack:"created_at"`
}
