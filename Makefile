.PHONY: abigen

abigen:
	mkdir -p bindings/uniswap_v3/factory
	mkdir -p bindings/uniswap_v3/pool
	abigen --abi abi/uniswap_v3/factory.abi.json --pkg univ3factory --out bindings/uniswap_v3/factory/factory.go
	abigen --abi abi/uniswap_v3/pool.abi.json --pkg univ3pool --out bindings/uniswap_v3/pool/pool.go