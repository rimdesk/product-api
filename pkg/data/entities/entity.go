package entities

import (
	"log"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	productv1 "github.com/rimdesk/product-api/gen/protos/rimdesk/product/v1"
	"github.com/rimdesk/product-api/pkg/data/domains"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	CompanyID   string
	Type        string
	Name        string `gorm:"uniqueIndex:idx_name_barcode"`
	CategoryID  string
	Barcode     string `gorm:"uniqueIndex:idx_name_barcode"`
	Description string
	Amount      float64
	SupplyPrice float64
	RetailPrice float64
}

func (p *Product) ToDomain() *domains.ProductDomain {
	domain := new(domains.ProductDomain)
	if err := copier.Copy(domain, p); err != nil {
		log.Println("failed to copy entity to domain:", err)
	}
	return domain
}

func (p *Product) BeforeCreate(*gorm.DB) error {
	p.ID = uuid.NewString()
	return nil
}

func (i *Product) ToProto() *productv1.Product {
	productDomain := new(productv1.Product)
	if err := copier.Copy(productDomain, i); err != nil {
		log.Println("failed to copy model to domain:", err)
	}
	productDomain.CreatedAt = timestamppb.New(i.CreatedAt)
	return &productv1.Product{}
}


func NewProductFromRequest(request *productv1.ProductRequest) (*Product, error) {
	newProduct := new(Product)
	if err := copier.Copy(newProduct, request); err != nil {
		return nil, err
	}
	return newProduct, nil
}