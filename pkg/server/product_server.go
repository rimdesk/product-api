package server

import (
	"context"

	"connectrpc.com/connect"
	productv1 "github.com/rimdesk/product-api/gen/protos/rimdesk/product/v1"
	"github.com/rimdesk/product-api/gen/protos/rimdesk/product/v1/productv1connect"
	"github.com/rimdesk/product-api/pkg/service"
)

type productServer struct {
	productService service.ProductService
}

// SearchWarehouse implements productv1connect.ProductServiceHandler.
func (server *productServer) SearchWarehouse(ctx context.Context, request *connect.Request[productv1.SearchWarehouseRequest]) (*connect.Response[productv1.SearchWarehouseResponse], error) {
	product, _ := server.productService.SearchWarehouse(ctx, request)
	response := &productv1.SearchWarehouseResponse{
		Product: product.Product,
	}
	return connect.NewResponse(response), nil
}

func (server *productServer) CreateProduct(ctx context.Context, request *connect.Request[productv1.CreateProductRequest]) (*connect.Response[productv1.CreateProductResponse], error) {
	product, _ := server.productService.CreateProduct(ctx, request)
	return connect.NewResponse(&productv1.CreateProductResponse{Product: product.Product}), nil
}

func (server *productServer) DeleteProduct(ctx context.Context, request *connect.Request[productv1.DeleteProductRequest]) (*connect.Response[productv1.DeleteProductResponse], error) {
	_, _ = server.productService.DeleteProduct(ctx, request)
	return connect.NewResponse(&productv1.DeleteProductResponse{}), nil
}

func (server *productServer) GetProduct(ctx context.Context, request *connect.Request[productv1.GetProductRequest]) (*connect.Response[productv1.GetProductResponse], error) {
	product, _ := server.productService.GetProduct(ctx, request)
	response := &productv1.GetProductResponse{
		Product: product.Product,
	}
	return connect.NewResponse(response), nil
}

func (server *productServer) ListProducts(ctx context.Context, request *connect.Request[productv1.ListProductsRequest]) (*connect.Response[productv1.ListProductsResponse], error) {
	response, _ := server.productService.ListProducts(ctx, request)
	return connect.NewResponse(response), nil
}

func (server *productServer) UpdateProduct(ctx context.Context, request *connect.Request[productv1.UpdateProductRequest]) (*connect.Response[productv1.UpdateProductResponse], error) {
	product, _ := server.productService.UpdateProduct(ctx, request)
	return connect.NewResponse(&productv1.UpdateProductResponse{
		Product: product.Product,
	}), nil
}

func NewProductServer(productService service.ProductService) productv1connect.ProductServiceHandler {
	return &productServer{productService: productService}
}
