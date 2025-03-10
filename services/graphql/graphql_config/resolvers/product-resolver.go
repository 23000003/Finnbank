package resolvers

import (
	"finnbank/services/common/grpc/products"

	"github.com/graphql-go/graphql"
)

func (s *StructGraphQLResolvers) GetProductQueryType(prodServer products.ProductServiceClient) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				// http://localhost:8083/graphql/product?query={product(id:1){name,info,price}}
				"product": &graphql.Field{
					Type:        productType,
					Description: "Get product by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Args["id"].(int)
						if ok {
							// Call gRPC product service
							res, err := prodServer.GetByIdProduct(p.Context, &products.ProductByIdRequest{
								ID: int64(id),
							})
							return res.Product, err
						}
						return nil, nil
					},
				},
				// http://localhost:8083/graphql/product?query={list{id,name,info,price}}
				"list": &graphql.Field{
					Type:        graphql.NewList(productType),
					Description: "Get product list",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {

						// Call gRPC product service
						res, err := prodServer.GetAllProducts(p.Context, &products.GetAllProductsRequest{
							Message: "HEY",
						})

						s.log.Info("Response: %v", res)

						return res.Product, err
					},
				},
			},
		})
}

func (s *StructGraphQLResolvers) GetProductMutationType(prodServer products.ProductServiceClient) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			/* Create new product item
			http://localhost:8083/graphql/product?query=mutation+_{create(name:"a Kaoala",info:"Inca on verbena (wiki)",price:1.99){id,name,info,price}}
			// */
			"create": &graphql.Field{
				Type:        productType,
				Description: "Create new product",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"info": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"price": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					name := params.Args["name"].(string)
					info := params.Args["info"].(string)
					price := float32(params.Args["price"].(float64))

					// Call gRPC product service
					res, err := prodServer.CreateProduct(params.Context, &products.CreateProductRequest{
						Name:  name,
						Info:  info,
						Price: price,
					})

					if err != nil {
						return nil, err
					}
					return res.Product, nil
				},
			},

			/* Update product by id
			   http://localhost:8083/graphql/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
			*/
			"update": &graphql.Field{
				Type:        productType,
				Description: "Update product by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"info": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"price": &graphql.ArgumentConfig{
						Type: graphql.Float,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// id, _ := params.Args["id"].(int)
					// name, nameOk := params.Args["name"].(string)
					// info, infoOk := params.Args["info"].(string)
					// price, priceOk := params.Args["price"].(float64)
					// product := Product{}

					// Call gRPC product service

					// Too lazy to implement since this is just tests
					//  but just the same concept if u understand it :)

					// for i, p := range products {
					// 	if int64(id) == p.ID {
					// 		if nameOk {
					// 			products[i].Name = name
					// 		}
					// 		if infoOk {
					// 			products[i].Info = info
					// 		}
					// 		if priceOk {
					// 			products[i].Price = price
					// 		}
					// 		product = products[i]
					// 		break
					// 	}
					// }
					return nil, nil
				},
			},

			/* Delete product by id
			   http://localhost:8083/graphql/product?query=mutation+_{delete(id:1){id,name,info,price}}
			*/
			"delete": &graphql.Field{
				Type:        productType,
				Description: "Delete product by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// id, _ := params.Args["id"].(int)
					// product := Product{}
					// for i, p := range products {
					// 	if int64(id) == p.ID {
					// 		product = products[i]
					// 		// Remove from product list
					// 		products = append(products[:i], products[i+1:]...)
					// 	}
					// }

					// Call gRPC product service
					// Too lazy to implement since this is just tests
					//  but just the same concept if u understand it :)

					return nil, nil
				},
			},
		},
	})
}
