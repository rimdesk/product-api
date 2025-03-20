package server

import (
	"context"

	"connectrpc.com/connect"
	productv1 "github.com/rimdesk/product-api/gen/rimdesk/product/v1"
	"github.com/rimdesk/product-api/gen/rimdesk/product/v1/productv1connect"
	"github.com/rimdesk/product-api/pkg/types"
)

type productServer struct {
	productService types.ProductService
}


func (server *productServer) CreateProduct(ctx context.Context, request *connect.Request[productv1.CreateProductRequest]) (*connect.Response[productv1.CreateProductResponse], error) {
	product, err := server.productService.CreateProduct(ctx, request)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&productv1.CreateProductResponse{Product: product}), nil
}

func (server *productServer) DeleteProduct(ctx context.Context, request *connect.Request[productv1.DeleteProductRequest]) (*connect.Response[productv1.DeleteProductResponse], error) {
	if err := server.productService.DeleteProduct(ctx, request); err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.DeleteProductResponse{}), nil
}

func (server *productServer) GetProduct(ctx context.Context, request *connect.Request[productv1.GetProductRequest]) (*connect.Response[productv1.GetProductResponse], error) {
	product, err := server.productService.GetProduct(ctx, request)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&productv1.GetProductResponse{Product: product}), nil
}

func (server *productServer) ListProducts(ctx context.Context, request *connect.Request[productv1.ListProductsRequest]) (*connect.Response[productv1.ListProductsResponse], error) {
	products, err := server.productService.ListProducts(ctx, request)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&productv1.ListProductsResponse{Products: products}), nil
}

func (server *productServer) UpdateProduct(ctx context.Context, request *connect.Request[productv1.UpdateProductRequest]) (*connect.Response[productv1.UpdateProductResponse], error) {
	product, err  := server.productService.UpdateProduct(ctx, request)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&productv1.UpdateProductResponse{
		Product: product,
	}), nil
}

func NewProductServer(productService types.ProductService) productv1connect.ProductServiceHandler {
	return &productServer{productService: productService}
}
