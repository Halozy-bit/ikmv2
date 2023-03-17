package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/ikmv2/backend/pkg/cache"
	"github.com/ikmv2/backend/pkg/helper"
	"github.com/ikmv2/backend/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidID = fmt.Errorf("invalid id")
)

type Service interface {
	CatalogList(context.Context, int) ([]repository.CatalogDisplay, error)
	CatalogListByCategory(context.Context, int, string) ([]repository.CatalogDisplay, error)
	Product(context.Context, primitive.ObjectID) (Product, error)
	Umkm(context.Context, primitive.ObjectID) (repository.Umkm, error)
	TotalPage(...string) int
}

type ServiceCirclePage struct {
	repo repository.Repository
}

// pagination with index in cache
// return helper.MaxPerPage array leng or less
// id reference per page from cache.Pagination
func (s *ServiceCirclePage) CatalogList(ctx context.Context, page int) ([]repository.CatalogDisplay, error) {
	id := cache.Pagination.Page(page)
	if id == primitive.NilObjectID {
		return nil, fmt.Errorf("page not found")
	}

	// NOTE
	// potential bug amount if maualy add data to database
	// add to cache pagination
	totalProduct, err := s.repo.CountCatalog(ctx)
	if err != nil {
		return nil, err
	}

	TotalProductNextPage := helper.CountTtlProductNxtPage(page, int(totalProduct))
	if TotalProductNextPage < 1 {
		return nil, mongo.ErrNoDocuments
	}

	return s.fetchCatalog(ctx, id, page, TotalProductNextPage)
}

// TODO
// make it works
func (s *ServiceCirclePage) CatalogListByCategory(ctx context.Context, page int, category string) ([]repository.CatalogDisplay, error) {
	id := cache.Pagination.CategoryPage(category, page)
	if id == primitive.NilObjectID {
		return s.repo.CatalogFirstPageWithCategory(ctx, int64(helper.MaxProductPerPage), category)
	}

	// NOTE
	// potential bug amount if maualy add data to database
	// add to cache pagination
	totalProduct, err := s.repo.CountCatalogCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	TotalProductNextPage := helper.CountTtlProductNxtPage(page, int(totalProduct))
	if TotalProductNextPage < 1 {
		return nil, mongo.ErrNoDocuments
	}

	return s.fetchCatalog(ctx, id, page, TotalProductNextPage, category)
}

type Product struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"nama" json:"nama"`
	Category    string             `bson:"kategori" json:"kategori"`
	Weight      string             `bson:"ukuran" json:"ukuran"`
	Pirt        string             `bson:"pirt" json:"pirt"`
	Variant     string             `bson:"varian" json:"varian"`
	Composition string             `bson:"komposisi" json:"komposisi"`
	Description string             `bson:"deskripsi" json:"deskripsi"`
	Owner       string             `bson:"owner" json:"owner"`
	Foto        repository.Foto    `bson:"foto" json:"foto"`
}

func (s *ServiceCirclePage) Product(ctx context.Context, id primitive.ObjectID) (Product, error) {
	p, err := s.repo.FindProduct(ctx, id)
	if err != nil {
		return Product{}, err
	}

	prod := Product{
		Id:          p.Id,
		Name:        p.Name,
		Category:    p.Category,
		Weight:      strings.Join(p.Weight, ", "),
		Variant:     strings.Join(p.Variant, ", "),
		Pirt:        p.Pirt,
		Composition: strings.Join(p.Composition, ", "),
		Description: p.Description,
		Owner:       p.Owner,
		Foto:        p.Foto,
	}
	return prod, nil
}

func (s *ServiceCirclePage) Umkm(ctx context.Context, id primitive.ObjectID) (repository.Umkm, error) {
	return s.repo.FindUmkm(ctx, id)
}

// fetch catalog like fetching circle catalog
// it used for page two or more in pagination
// @id item starts from, @contentLimit limit catalog you want to fetch
// if the pagination is at the bottom of the item, while there are still items that have not been returned
// then a batch 2 query will be executed to the top item in the database database
func (s *ServiceCirclePage) fetchCatalog(ctx context.Context, id primitive.ObjectID, page int, contentLimit int, category ...string) (ctlgLs []repository.CatalogDisplay, err error) {
	if len(category) > 0 {
		ctlgLs, err = s.repo.CatalogGteIdByCategory(ctx, id, int64(contentLimit), category[0])
	} else {
		ctlgLs, err = s.repo.CatalogGteId(ctx, id, int64(contentLimit))
	}

	if err != nil {
		if err != mongo.ErrEmptySlice {
			return nil, err
		}
	}

	// insufficient number of items that should be returned
	contentLimit = contentLimit - len(ctlgLs)
	if contentLimit <= 0 {
		return ctlgLs, nil
	}

	ctlgLsFirst, err := s.repo.CatalogFirstLine(ctx, int64(contentLimit))
	if err != nil {
		return nil, err
	}

	ctlgLs = append(ctlgLs, ctlgLsFirst...)
	return
}

func (s *ServiceCirclePage) TotalPage(category ...string) int {
	if len(category) > 0 {
		return cache.Pagination.LengCategory(category[0])
	}
	return cache.Pagination.Leng()
}
