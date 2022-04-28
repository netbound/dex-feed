package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/netbound/dex-feed/db"
	uniswapv3 "github.com/netbound/dex-feed/exchange/uniswap_v3"
	"github.com/netbound/dex-feed/token"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

func main() {
	var (
		factoryAddress = common.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984")
		rpcApi         = "http://localhost:8080/eth"
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcApi)
	if err != nil {
		log.Fatal(err)
	}

	dbOpts := db.Opts{Persistent: true, CacheSize: 2048}
	tdb := token.NewTokenDB(client, dbOpts)

	univ3 := uniswapv3.New(client, tdb, factoryAddress, dbOpts)

	router := gin.Default()

	router.Use(gin.Logger())

	router.GET("/price/:base", func(c *gin.Context) {
		base := common.HexToAddress(c.Param("base"))
		quote := common.HexToAddress(c.Query("quote"))
		fee, err := strconv.ParseInt(c.Query("fee"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(base, quote, fee)

		price, err := univ3.GetPrice(c.Request.Context(), quote, base, fee)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"price": price})
	})

	router.Run(":3000")
}
