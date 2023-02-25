package db

import "go.mongodb.org/mongo-driver/mongo/options"

func GetClientOptions() *options.ClientOptions {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://ghousedev:12W2r3w4r568@ghouse.4azrsyy.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)

	return clientOptions
}
