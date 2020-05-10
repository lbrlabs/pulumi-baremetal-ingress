package main

type AddressPool struct {
	Name      string   `json:"name"`
	Protocol  string   `json:"protocol"`
	Addresses []string `json:"addresses"`
}
type MetallbConfig struct {
	AddressPools []AddressPool `json:"address-pools"`
}
