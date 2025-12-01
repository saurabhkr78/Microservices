package catalog

import (
	"context"

	"encoding/json"
	"log"

	"errors"
	elastic "gopkg.in/olivere/elastic.v5"
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}
type elasticRepository struct {
	Client *elastic.Client
}

var (
	ErrNotFound = errors.New("Entity not found")
)

type productDocument struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client}, nil

}
func (r *elasticRepository) Close() {
}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	_, err := r.Client.Index().
		Index("catalog").
		Id(p.ID).
		BodyJson(productDocument{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}).Do(ctx)
	return err
}
func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	res, err := r.Client.Get().
		Index("catalog").
		Type("Product").
		Id(id).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, ErrNotFound
	}
	P := &productDocument{}
	if err := json.Unmarshal(*res.Source, &P); err != nil {
		return nil, err
	}
	return &Product{
		ID:          res.Id,
		Name:        P.Name,
		Description: P.Description,
		Price:       P.Price,
	}, err
}
func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	res, err := r.Client.Search().
		Index("catalog").
		Type("Product").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println("Error executing search:", err)
		return nil, err
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := &productDocument{}
		if err := json.Unmarshal(*hit.Source, &p); err != nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}

	}
	return products, err
}
func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(
			items, elastic.NewMultiGetItem().
				Index("catalog").
				Type("Product").
				Id(id),
		)
	}
	res, err := r.Client.MultiGet().
		Add(items...).
		Do(ctx)
	if err != nil {
		log.Println("Error executing multi get:", err)
		return nil, err
	}
	products := []Product{}
	for _, doc := range res.Docs {
		p := &productDocument{}
		if err := json.Unmarshal(*doc.Source, &p); err != nil {
			products = append(products, Product{
				ID:          doc.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, err
}
func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	res, err := r.Client.Search().
		Index("catalog").
		Type("Product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println("Error executing search:", err)
		return nil, err
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := &productDocument{}
		if err := json.Unmarshal(*hit.Source, &p); err != nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, err
}
