package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CategoryField = "kategori"
)

// TODO
// Foto change to array focuse on per product
// Thumbnail for displaying catalog

type CatalogDisplay struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"nama" json:"nama"`
	Category  string             `bson:"kategori" json:"kategori"`
	Owner     string             `bson:"owner" json:"owner"`
	Thumbnail string             `bson:"thumbnail" json:"thumbnail"`
}

type Product struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"nama" json:"nama"`
	Category    string             `bson:"kategori" json:"kategori"`
	Weight      []string           `bson:"ukuran" json:"ukuran"`
	Pirt        string             `bson:"pirt" json:"pirt"`
	Variant     []string           `bson:"varian" json:"varian"`
	Composition []string           `bson:"komposisi" json:"komposisi"`
	Description string             `bson:"deskripsi" json:"deskripsi"`
	Owner       string             `bson:"owner" json:"owner"`
	Foto        Foto               `bson:"foto" json:"foto"`
}

type Foto struct {
	Cover   string `bson:"cover" json:"cover"`
	Detail1 string `bson:"detail1" json:"detail1"`
	Detail2 string `bson:"detail2" json:"detail2"`
}

type Umkm struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"nama" json:"nama"`
	Description string             `bson:"deskripsi" json:"deskripsi"`
	BadanHukum  string             `bson:"badan_hukum" json:"badan_hukum"`
	Branding    string             `bson:"branding" json:"branding"`
	Marketplace Marketplace        `bson:"marketplace" json:"marketplace"`
	SocialMedia SocialMedia        `bson:"social_media" json:"social_media"`
	Foto        FotoOwner          `bson:"foto_owner" json:"foto_owner"`
	Story       []string           `bson:"cerita" json:"cerita"`
}

type SocialMedia struct {
	Instagram string `bson:"instagram" json:"instagram"`
	Tiktok    string `bson:"tiktok" json:"tiktok"`
	Facebook  string `bson:"facebook" json:"facebook"`
	Youtube   string `bson:"youtube" json:"youtube"`
}

type FotoOwner struct {
	FotoRompi     string `bson:"f_rompi" json:"f_rompi"`
	FotoWawancara string `bson:"f_wawancara" json:"f_wawancara"`
	FotoProduksi  string `bson:"f_produksi" json:"f_produksi"`
}

type Marketplace struct {
	Shopee     string `bson:"shopee" json:"shopee"`
	Tokopedia  string `bson:"tokopedia" json:"tokopedia"`
	TiktokShop string `bson:"tt_shop" json:"tiktok_shop"`
	Web        string `bson:"web" json:"web"`
}

func ProductToDocument(p Product) bson.D {
	return bson.D{
		{Key: "nama", Value: p.Name},
		{Key: "kategori", Value: p.Category},
		{Key: "ukuran", Value: p.Weight},
		{Key: "varian", Value: p.Variant},
		{Key: "pirt", Value: p.Pirt},
		{Key: "komposisi", Value: p.Composition},
		{Key: "deskripsi", Value: p.Description},
		{Key: "owner", Value: p.Owner},
		{Key: "foto", Value: bson.D{
			{Key: "cover", Value: p.Foto.Cover},
			{Key: "detail1", Value: p.Foto.Detail1},
			{Key: "detail2", Value: p.Foto.Detail2},
		}},
	}
}

func UmkmToDocument(u Umkm) bson.D {
	return bson.D{
		{Key: "nama", Value: u.Name},
		{Key: "deskripsi", Value: u.Description},
		{Key: "badan_hukum", Value: u.BadanHukum},
		{Key: "branding", Value: u.Branding},
		{Key: "marketplace", Value: bson.D{
			{Key: "shopee", Value: u.Marketplace.Shopee},
			{Key: "tokopedia", Value: u.Marketplace.Tokopedia},
			{Key: "tt_shop", Value: u.Marketplace.TiktokShop},
			{Key: "web", Value: u.Marketplace.Web},
		}},
		{Key: "social_media", Value: bson.D{
			{Key: "instagram", Value: u.SocialMedia.Instagram},
			{Key: "tiktok", Value: u.SocialMedia.Tiktok},
			{Key: "facebook", Value: u.SocialMedia.Facebook},
			{Key: "youtube", Value: u.SocialMedia.Youtube},
		}},
		{Key: "foto_owner", Value: bson.D{
			{Key: "f_rompi", Value: u.Foto.FotoRompi},
			{Key: "f_wawancara", Value: u.Foto.FotoWawancara},
			{Key: "f_produksi", Value: u.Foto.FotoProduksi},
		}},
		{Key: "cerita", Value: u.Story},
	}
}
