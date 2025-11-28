package main

type Server struct {
	// accountClient *account.Client
	// catalogClient *catalog.Client
	// orderClient   *order.Client
}

func NewGraphQLServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {
	// 	accountClient, err := account.NewClient(accountUrl)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	catalogClient, err := catalog.NewClient(catalogUrl)
	// 	if err != nil {
	// 		accountClient.Close() // close parent client or on which  dependent client
	// 		return nil, err
	// 	}

	// 	orderClient, err := order.NewClient(orderUrl)
	// 	if err != nil {
	// 		accountClient.Close()
	// 		catalogClient.Close()
	// 		return nil, err
	// 	}

	return &Server{
		//		accountClient,
		//		catalogClient,
		//		orderClient,
	}, nil
}

//	func (s *server) Mutation() MutationResolver {
//		return &mutationResolver{server: s}
//	}
//
//	func (s *server) Query() QueryResolver {
//		return &queryResolver{server: s}
//	}
//
//	func (s *server) Account() AccountResolver {
//		return &accoountResolver{server: s}
//	}
func (s *server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
