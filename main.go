package main

import (
	"fmt"

	"./p1"
)

func main() {
	fmt.Println("Project Merkle Patricia Trie")

	mpt := p1.MerklePatriciaTrie{}

	//mpt.Initial()
	// mpt.Insert("a", "apple")
	// mpt.Insert("b", "banana")
	// mpt.Insert("p", "papaya") //TODO -- uncomment and see if insert working without error
	// mpt.Insert("c", "cream")
	// mpt.Insert("aa", "aananas")
	//mpt.Delete("aa")

	mpt.Insert("a", "apple")
	mpt.Insert("p", "banana")
	//mpt.Insert("aaap", "orange")
	//inserted_trie := mpt.Order_nodes()
	mpt.Insert("a", "new")

	fmt.Println("===============Trie===============")
	inserted_trie := mpt.Order_nodes()
	fmt.Println(inserted_trie)
	//fmt.Println("===============Trie===============")

}
