package routes

import "github.com/rimdesk/product-api/gen/rimdesk/product/v1/productv1connect"

var PublicRoutes = []string{
	productv1connect.ProductServiceGetProductProcedure,
}