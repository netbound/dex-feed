.PHONY: abigen

abigen:
	mkdir -p bindings/uniswap_v3
	abigen --abi abi/uniswap_v3/factory.abi.json --pkg univ3factory --out bindings/uniswap_v3/factory.go