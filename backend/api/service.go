package api

import (
	"context"
	"fmt"
	"log"

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
	CatalogList(context.Context, int) ([]repository.DocCatalog, error)
	CatalogListByCategory(context.Context, int, string, string) ([]repository.DocCatalog, error)
}

type ServiceCirclePage struct {
	repo repository.Repository
}

// pagination with index in cache
// return helper.MaxPerPage array leng or less
// id reference per page from cache.Pagination
func (s *ServiceCirclePage) CatalogList(ctx context.Context, page int) ([]repository.DocCatalog, error) {
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
	log.Print("next page: ", TotalProductNextPage)
	if TotalProductNextPage < 1 {
		return nil, mongo.ErrNoDocuments
	}

	return s.fetchCatalog(ctx, id, page, TotalProductNextPage)
}

// TODO
// rotation in specific category
func (s *ServiceCirclePage) CatalogListByCategory(ctx context.Context, page int, category string, lastId string) ([]repository.DocCatalog, error) {
	if page == 1 {
		ctlg, err := s.repo.CatalogFirstPageWithCategory(ctx, int64(helper.MaxProductPerPage), category)
		return ctlg, err
	}

	id, idErr := primitive.ObjectIDFromHex(lastId)
	if idErr != nil {
		return nil, fmt.Errorf("invalid id")
	}

	return s.repo.CatalogGtId(ctx, id, int64(helper.MaxProductPerPage))
}

// fetch catalog like fetching circle catalog
// it used for page two or more in pagination
// @id item starts from, @contentLimit limit catalog you want to fetch
// if the pagination is at the bottom of the item, while there are still items that have not been returned
// then a batch 2 query will be executed to the top item in the database database
func (s *ServiceCirclePage) fetchCatalog(ctx context.Context, id primitive.ObjectID, page int, contentLimit int) (ctlgLs []repository.DocCatalog, err error) {
	ctlgLs, err = s.repo.CatalogGteId(ctx, id, int64(contentLimit))

	log.Print("batch 1 get:", len(ctlgLs))
	log.Println(err)
	if err != nil {
		if err != mongo.ErrEmptySlice {
			return nil, err
		}
	}

	// insufficient number of items that should be returned
	contentLimit = contentLimit - len(ctlgLs)
	log.Print("deviation: ", contentLimit)
	if contentLimit <= 0 {
		return ctlgLs, nil
	}

	ctlgLsFirst, err := s.repo.CatalogFirstLine(ctx, int64(contentLimit))
	log.Print("batch 2 get:", len(ctlgLsFirst))
	if err != nil {
		return nil, err
	}

	ctlgLs = append(ctlgLs, ctlgLsFirst...)
	return
}
