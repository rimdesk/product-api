package service

import (
	"context"
	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	"github.com/jinzhu/copier"
	productv1 "github.com/rimdesk/product-api/gen/rimdesk/product/v1"
	"github.com/rimdesk/product-api/pkg/data/entities"
	"github.com/rimdesk/product-api/pkg/types"
)



type productService struct {
	productRepository types.ProductRepository
}


func NewProductService(productRepository types.ProductRepository) types.ProductService {
	return &productService{productRepository: productRepository}
}

func (service *productService) ListProducts(ctx context.Context, request *connect.Request[productv1.ListProductsRequest]) ([]*productv1.Product, error) {
	products, err := service.productRepository.FindAll(request.Msg.String())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	productProtos := make([]*productv1.Product, len(products))
	for _, product := range products {
		productProtos = append(productProtos, product.ToProto())
	}

	return productProtos, nil
}


func (service *productService) GetProduct(ctx context.Context, request *connect.Request[productv1.GetProductRequest]) (*productv1.Product, error) {
	 product,err := service.productRepository.FindById(request.Msg.GetId()); 
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	return product.ToProto(), nil
}

func (service *productService) CreateProduct(ctx context.Context, request *connect.Request[productv1.CreateProductRequest]) (*productv1.Product, error) {
	productRequest := request.Msg.GetProduct()

	if err := protovalidate.Validate(productRequest); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	product, err := entities.NewProductFromRequest(productRequest)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := service.productRepository.Create(product); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return product.ToProto(), nil
}


func (service *productService) UpdateProduct(ctx context.Context, request *connect.Request[productv1.UpdateProductRequest]) (*productv1.Product, error) {

	product, err := service.productRepository.FindById(request.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	productRequest := request.Msg.GetProduct()
	if err := protovalidate.Validate(productRequest); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := copier.Copy(product, productRequest); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := service.productRepository.Update(product); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return product.ToProto(), nil
}

func (service *productService) DeleteProduct(ctx context.Context, request *connect.Request[productv1.DeleteProductRequest]) error {
	product, err := service.productRepository.FindById(request.Msg.GetId())
	if err != nil {
		return connect.NewError(connect.CodeNotFound, err)
	}

	if err := service.productRepository.Delete(product); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	return nil
}

