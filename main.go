package main

import (
	"fmt"

	"./p1"
)

func main() {
	fmt.Println("Project Merkle Patricia Trie")

	mpt := p1.MerklePatriciaTrie{}

	//mpt.Initial()
	mpt.Insert("a", "10")
	mpt.Insert("b", "20")
	mpt.Insert("p", "30")
	mpt.Insert("c", "40")
	mpt.Insert("aa", "50")

	mpt.Delete("aa")
	fmt.Println("===============Trie===============")
	inserted_trie := mpt.Order_nodes()
	fmt.Println(inserted_trie)

}
