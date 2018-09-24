package battleground

import (
	"fmt"

	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/go-loom/util"
	"github.com/loomnetwork/zombie_battleground/types/zb"
)

type Versioned interface {
	GetVersion() string
	MakeKey([]byte) []byte
	ListHeroes(contract.StaticContext, *zb.ListHeroesRequest) (*zb.ListHeroesResponse, error)
	GetHero(ctx contract.StaticContext, req *zb.GetHeroRequest) (*zb.GetHeroResponse, error)
	GetCollection(ctx contract.StaticContext, req *zb.GetCollectionRequest) (*zb.GetCollectionResponse, error)
}

type V1 struct {
	KeyPrefix string
}

func getVersionedObject(version string) (Versioned, error) {
	if version == "v1" {
		return &V1{}, nil
	}
	return nil, fmt.Errorf("version not found")
}

func (v *V1) GetVersion() string {
	return "v1"
}

func (v *V1) MakeKey(key []byte) []byte {
	return util.PrefixKey([]byte(v.GetVersion()), key)
}

func (v *V1) ListHeroes(ctx contract.StaticContext, req *zb.ListHeroesRequest) (*zb.ListHeroesResponse, error) {
	heroList, err := loadHeroes(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &zb.ListHeroesResponse{Heroes: heroList.Heroes}, nil
}

func (v *V1) GetHero(ctx contract.StaticContext, req *zb.GetHeroRequest) (*zb.GetHeroResponse, error) {
	heroList, err := loadHeroes(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	hero := getHeroById(heroList.Heroes, req.HeroId)
	if hero == nil {
		return nil, contract.ErrNotFound
	}
	return &zb.GetHeroResponse{Hero: hero}, nil
}

func (v *V1) GetCollection(ctx contract.StaticContext, req *zb.GetCollectionRequest) (*zb.GetCollectionResponse, error) {
	collectionList, err := loadCardCollection(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &zb.GetCollectionResponse{Cards: collectionList.Cards}, nil
}
