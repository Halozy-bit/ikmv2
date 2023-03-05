package api

import (
	"context"
	"fmt"

	"github.com/ikmv2/backend/pkg/cache"
	"github.com/ikmv2/backend/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const MaxProductPerPage int = 16

type Service struct {
	repo repository.Repository
}

// pagination with index in cache
func (s Service) CatalogList(ctx context.Context, page int, lastId string) ([]repository.DocCatalog, error) {
	if page <= 1 || lastId == "" {
		page = 1
		lastId = cache.Get(cache.TopCatalog).(string)
	}

	id, idErr := primitive.ObjectIDFromHex(lastId)
	if idErr != nil {
		return nil, fmt.Errorf("invalid id")
	}

	totalProduct, err := s.repo.CountCatalog(ctx)
	if err != nil {
		return nil, err
	}

	TotalProductNextPage := CountTtlProductNxtPage(page, int(totalProduct))
	ctlgLs, err := s.fetchCatalog(ctx, id, TotalProductNextPage)

	return ctlgLs, err
}

func (s Service) CatalogListUsingCategory(ctx context.Context, page int, category string, lastId string) ([]repository.DocCatalog, error) {
	if page == 1 {
		ctlg, err := s.repo.CatalogFirstPage(ctx, int64(MaxProductPerPage))
		return ctlg, err
	}

	id, idErr := primitive.ObjectIDFromHex(lastId)
	if idErr != nil {
		return nil, fmt.Errorf("invalid id")
	}

	return s.repo.CatalogPageFromId(ctx, id, int64(MaxProductPerPage))
}

// fetch catalog is fetching circle catalog
// @id item starts from, @contentLimit limit catalog you want to fetch
// if the pagination is at the bottom of the item, while there are still items that have not been returned
// then a batch 2 query will be executed to the top item in the database database
func (s Service) fetchCatalog(ctx context.Context, id primitive.ObjectID, contentLimit int) ([]repository.DocCatalog, error) {
	ctlgLs, err := s.repo.CatalogPageFromId(ctx, id, int64(contentLimit))
	if err != nil {
		return nil, err
	}

	deviation := contentLimit - len(ctlgLs)

	if deviation == 0 {
		return ctlgLs, nil
	}

	ctlgLs_batch2, err := s.repo.CatalogFirstPage(ctx, int64(deviation))
	if err != nil {
		return nil, err
	}

	ctlgLs = append(ctlgLs, ctlgLs_batch2...)
	return ctlgLs, nil
}
