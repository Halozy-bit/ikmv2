package api

import (
	"context"
	"fmt"

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
	CatalogListByCategory(context.Context, int, string) ([]repository.DocCatalog, error)
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
	if TotalProductNextPage < 1 {
		return nil, mongo.ErrNoDocuments
	}

	return s.fetchCatalog(ctx, id, page, TotalProductNextPage)
}

// TODO
// make it works
func (s *ServiceCirclePage) CatalogListByCategory(ctx context.Context, page int, category string) ([]repository.DocCatalog, error) {
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

// fetch catalog like fetching circle catalog
// it used for page two or more in pagination
// @id item starts from, @contentLimit limit catalog you want to fetch
// if the pagination is at the bottom of the item, while there are still items that have not been returned
// then a batch 2 query will be executed to the top item in the database database
func (s *ServiceCirclePage) fetchCatalog(ctx context.Context, id primitive.ObjectID, page int, contentLimit int, category ...string) (ctlgLs []repository.DocCatalog, err error) {
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
