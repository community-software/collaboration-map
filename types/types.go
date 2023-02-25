package types

type Test struct {
	Msg string `bson:"msg, omitempty"`
}

type Spot struct {
	Title       string  `bson:"title, omitempty"`
	Description string  `bson:"description, omitempty"`
	Category    string  `bson:"category, omitempty"`
	Lat         float64 `bson:"lat, omitempty"`
	Lon         float64 `bson:"lon, omitempty"`
	Pic         string  `bson:"pic, omitempty"`
	Author      string  `bson:"author, omitempty"`
}
type User struct {
	Username string  `bson:"username, omitempty"`
	ID       int64   `bson:"id, omitempty"`
	Chats    []int64 `bson:"chats, omitempty"`
}
type Chat struct {
	ID    int64  `bson:"id, omitempty"`
	Title string `bson:"title, omitempty"`
}
type ActiveSpotCreation struct {
	ChatID int64 `bson:"chat_id, omitempty"`
	Step   int   `bson:"step, omitempty"`
	Data   Spot  `bson:"data, omitempty"`
}
type ActiveSpotCreations map[int64]ActiveSpotCreation
