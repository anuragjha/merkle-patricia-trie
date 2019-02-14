package p1

import (
	"errors"
	"fmt"
	"reflect"
)

func (mpt *MerklePatriciaTrie) delHelper(nodeKey string, currentNode Node, pathLeft []uint8, hashStack []string) (string, error) {
	fmt.Println("--------------------------------------------------------------------------------------")
	fmt.Println("HASH STACK :", hashStack)
	fmt.Println("--------------------------------------------------------------------------------------")
	if len(pathLeft) > 0 && currentNode.node_type != 0 { //path length >0
		fmt.Println("in Del - path length > 0")
		if currentNode.node_type == 1 { // branch and pathleft >0
			fmt.Println("in Del - path length > 0 - node is Branch")
			if currentNode.branch_value[pathLeft[0]] != "" {
				fmt.Println("Value exists in branch!!")
				hash := currentNode.branch_value[pathLeft[0]]
				fmt.Println("Hash found:", currentNode.branch_value[pathLeft[0]])
				fmt.Println("Hash found:", hash)
				oldhash := currentNode.hash_node()
				index := pathLeft[0]
				pathLeft = pathLeft[1:]
				fmt.Println("PathLeft from branch now:", pathLeft)

				next_node := mpt.db[hash]
				fmt.Println("Next Node:", next_node)

				// if next_node.node_type == 1 {
				// 	hashStack = append(hashStack, oldhash)
				// 	return mpt.delHelper(hash, mpt.db[hash], pathLeft, hashStack)

				// } else
				if next_node.node_type == 2 {
					fmt.Println("next_node.flag_value.encoded_prefix : ", next_node.flag_value.encoded_prefix)
					hex_prefix_array := AsciiArrayToHexArray(next_node.flag_value.encoded_prefix)
					fmt.Println("node_type == 2 - hex prefix: ", hex_prefix_array)
					if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) {
						fmt.Println("Next node is a leaf")
						currentNode.branch_value[index] = ""
						mpt.db[oldhash] = currentNode
						fmt.Println("Value cleared from branch:")
						fmt.Println("Branch node:", currentNode)
					}
				}

				hashStack = append(hashStack, oldhash)
				return mpt.delHelper(hash, mpt.db[hash], pathLeft, hashStack)

				// fmt.Println("PathLeft from branch :", pathLeft)
				// // if path ends in branch and so value can be  stored in empty leaf node
				// if len(pathLeft) == 0 { // pathleft gets over at branch node
				// 	fmt.Println("pathleft gets over at branch node")
				// 	if mpt.db[hash].node_type == 1 { // if this branch node contains a value of node
				// 		fmt.Println("next node :-  branch node")
				// 		if mpt.db[hash].branch_value[16] != "" {
				// 			fmt.Println("next node :-  branch node - value in pos 16")
				// 			hashStack = append(hashStack, hash)
				// 			branchnode := mpt.db[hash]
				// 			value := branchnode.branch_value[16]
				// 			branchnode.branch_value[16] = ""
				// 			fmt.Println("Deleting from db")
				// 			delete(mpt.db, hash) // delete :-> empty - value node
				// 			///  l!!!!!!!!! rearrange trie !!!!!
				// 			fmt.Println("rearranging trie : ")
				// 			mpt.rearrangeDeletedTrie(hashStack)
				// 			////// !!!!!!!!!!!!!!!!!!                 !!!!!!!!
				// 			fmt.Println("Returning with value : ", value)
				// 			return value, nil
				// 		} else {
				// 			fmt.Println("Returning with Error")
				// 			return "", errors.New("path_not_found")
				// 		}
				// 	} else if mpt.db[hash].node_type == 2 {
				// 		fmt.Println("next node :-  Ext/Leaf node")
				// 		hex_prefix_array := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
				// 		fmt.Println("node_type == 2 - hex prefix: ", hex_prefix_array)
				// 		if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) { //extension
				// 			fmt.Println("pathleft is greater than 0 - node_type 2 - prefix 0 or 1 Extension")
				// 			return "", errors.New("path_not_found")
				// 		} else if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) { //leaf
				// 			fmt.Println("pathleft is greater than 0 - node_type 2 - prefix 2 or 3 Leaf")
				// 			value := mpt.db[hash].flag_value.value
				// 			// if this ext/leaf node contains a value of node
				// 			// encoded prefix should be empty , if so then ////
				// 			fmt.Println("Deleting from db")
				// 			delete(mpt.db, hash) // delete :-> empty - value node
				// 			///  l!!!!!!!!! rearrange trie !!!!!
				// 			fmt.Println("rearranging trie : ")
				// 			mpt.rearrangeDeletedTrie(hashStack)
				// 			////// !!!!!!!!!!!!!!!!!!                 !!!!!!!!
				// 			fmt.Println("Returning with value : ", value)
				// 			return value, nil
				// 		} else {
				// 			fmt.Println("Returning with Error : ")
				// 			return "", errors.New("path_not_found")
				// 		}
				// 	} else {
				// 		fmt.Println("Returning with Error : ")
				// 		return "", errors.New("path_not_found")
				// 	}
				// } else { //pathLeft > 0
				// 	fmt.Println("pathleft traversing through branch node")
				// 	fmt.Println("currentNode.branch_value before making empty", currentNode.branch_value[pathLeft[0]])
				// 	fmt.Println("\nHash before removing value from the branch[16] field:", currentNode.hash_node())
				// 	// next_node := mpt.db[currentNode.branch_value[pathLeft[0]]]

				// 	fmt.Println("Next Node:", next_node)
				// 	if next_node.node_type == 2 {
				// 		fmt.Println("next_node.flag_value.encoded_prefix : ", next_node.flag_value.encoded_prefix)
				// 		hex_prefix_array := AsciiArrayToHexArray(next_node.flag_value.encoded_prefix)
				// 		fmt.Println("node_type == 2 - hex prefix: ", hex_prefix_array)
				// 		if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) {

				// 			next_node.branch_value[pathLeft[0]] = ""
				// 		}
				// 	}
				// 	fmt.Println("\nHash before removing value from the branch[16] field:", currentNode.hash_node())
				// 	fmt.Println("\nHASHSTACK before : ", hashStack)
				// 	hashStack = append(hashStack, hash)
				// 	fmt.Println("HASHSTACK after: ", hashStack)
				// 	pathLeft = pathLeft[1:]
				// 	fmt.Println("Pathleft : ", pathLeft)
				// 	return mpt.delHelper(hash, mpt.db[hash], pathLeft, hashStack)
				// }
			} else {
				fmt.Println("returning with error")
				return "", errors.New("path_not_found")
			}
		} else if currentNode.node_type == 2 { //ext or leaf and pathleft >0
			hex_prefix_array := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
			fmt.Println("node_type == 2 - hex prefix: ", hex_prefix_array)
			oldhash := currentNode.hash_node()

			if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) { //extension
				fmt.Println("pathleft is greater than 0 - node_type 2 - prefix 1 or 2 Extension")

				nodePath := compact_decode(currentNode.flag_value.encoded_prefix)
				fmt.Println("currentNode.flag_value.encoded_prefix:", currentNode.flag_value.encoded_prefix)
				fmt.Println("triePathtriePathtriePath:", nodePath)
				fmt.Println("pathLeftpathLeftpathLeft:", pathLeft)

				if reflect.DeepEqual(nodePath, pathLeft[:len(nodePath)]) {
					pathLeft = pathLeft[len(nodePath):]
					fmt.Println("Path left before new call:", pathLeft)
					if len(pathLeft) == 0 { // pathleft is zero now
						fmt.Println("Path left is zero now")
						hash := currentNode.flag_value.value
						fmt.Println("hash : ", hash)
						fmt.Println("HASHSTACK before : ", hashStack)
						hashStack = append(hashStack, oldhash)
						fmt.Println("HASHSTACK2 after: ", hashStack)
						//pathleft=0, use value from ext and check in branchvalu[16]
						if mpt.db[hash].branch_value[16] != "" { //value found
							fmt.Println("node Branch = valu at 16")
							//value := mpt.db[hash].branch_value[16]
							//delete node
							fmt.Println("so finding next node")
							nextnode := mpt.db[hash]
							fmt.Println("next node : ", nextnode)
							valuetoreturn := nextnode.branch_value[16]
							nextnode.branch_value[16] = ""
							mpt.db[hash] = nextnode
							fmt.Println("Value after deleted:", mpt.db[hash].branch_value[16])
							fmt.Println("made value at 16 - empty")
							//after cheching how many values in branch_value
							// delete(mpt.db, nodeKey)
							fmt.Println("HASHSTACK before : ", hashStack)
							hashStack = append(hashStack, hash)
							fmt.Println("HASHSTACK after  : ", hashStack)
							//rearrange trie ///// !!!!!!!! call func
							fmt.Println("calling rearrangeDeletedTrie")

							mpt.rearrangeDeletedTrie(hashStack)
							// !!!!!!!!!!!!!!!!!!!!!!!!!!
							fmt.Println("returning : ", valuetoreturn)
							return valuetoreturn, nil
						} //................................value not found
						fmt.Println("returning with error")
						return "", errors.New("path_not_found")
						//"pathleft is greater than 0 - node_type 2 - prefix 1 or 2 Extension")
					} else if len(pathLeft) > 0 { //add to hashstack call on next node
						fmt.Println("leftpath >0, ")
						hash := currentNode.flag_value.value
						fmt.Println("HASHSTACK before : ", hashStack)
						hashStack = append(hashStack, oldhash) //adding current ext
						fmt.Println("HASHSTACK after : ", hashStack)
						fmt.Println("Path left before new call:", pathLeft)
						fmt.Println("calling delHelper")
						return mpt.delHelper(hash, mpt.db[hash], pathLeft, hashStack)
					}
					fmt.Println("returning with error")
					return "", errors.New("path_not_found")
				} else {
					fmt.Println("returning with error")
					return "", errors.New("path_not_found")
				}
			} else if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) { //leaf //pathleft >0
				fmt.Println("pathleft is greater than 0 - node_type 2 - prefix 2 or 3 Leaf")
				// hex without prefix <(ascii with prefix)
				nodePath := compact_decode(currentNode.flag_value.encoded_prefix)
				fmt.Println("leftPath :", pathLeft)
				fmt.Println("nodePath :", nodePath)
				if reflect.DeepEqual(nodePath, pathLeft) {
					//delete node
					fmt.Println("deleting node from map")
					delete(mpt.db, nodeKey)
					fmt.Println("hashStack : ", hashStack)
					//rearrange trie ///// !!!!!!!!
					fmt.Println("calling rearrangeDeletedTrie")
					mpt.rearrangeDeletedTrie(hashStack)
					// !!!!!!!!!!!!!!!!!!!!!!!!!!
					fmt.Println("returning with value : ", currentNode.flag_value.value)
					return currentNode.flag_value.value, nil
				} else {
					fmt.Println("deep equal success")
					fmt.Println("returning with error")
					return "", errors.New("path_not_found")
				}

			}
		}
	} else if len(pathLeft) == 0 && currentNode.node_type != 0 { //pathlength ==0
		fmt.Println("Path left is zero...")
		if currentNode.node_type == 1 { // branch
			if currentNode.branch_value[16] != "" {
				fmt.Println("Lenght is zero")
				previous_hash := currentNode.hash_node()
				previous_value := currentNode.branch_value[16]
				currentNode.branch_value[16] = ""
				mpt.db[previous_hash] = currentNode
				hashStack = append(hashStack, previous_hash) //adding current ext
				mpt.rearrangeDeletedTrie(hashStack)
				return previous_value, nil
			} else {
				fmt.Println("MY ERROR: not found to delete")
			}
		} else if currentNode.node_type == 2 { //ext or leaf
			//extension
			fmt.Println("leaf or extension")
			hex_empty := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
			fmt.Println("Hex_empty:", hex_empty)

			if currentNode.flag_value.encoded_prefix[0] == 0 || currentNode.flag_value.encoded_prefix[0] == 1 {
				fmt.Println("My ERROR : Not found to delete")

			} else if hex_empty[0] == 2 || hex_empty[0] == 3 { //leaf
				hex_empty := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
				fmt.Println("Hex_empty:", hex_empty)
				if reflect.DeepEqual(hex_empty, []uint8{2, 0}) {
					fmt.Println("Leaf is empty...")
					fmt.Println("Found the key to be deleted")
					value := currentNode.flag_value.value

					delete(mpt.db, currentNode.hash_node())

					mpt.rearrangeDeletedTrie(hashStack)
					return value, nil
				} else {
					fmt.Println("My ERROR : Not found to delete")
				}
			}
		}
	} else if currentNode.node_type == 0 || len(pathLeft) < 0 { //path length <0 or nodetype =0
		return "", errors.New("path_not_found")
	}

	fmt.Println("EXit with Err or")
	fmt.Println("current node : ", currentNode)
	fmt.Println("Path left : ", pathLeft)

	return "", errors.New("path_not_found")
}

// rearrangeDeletedTrie rearranges mpt
func (mpt *MerklePatriciaTrie) rearrangeDeletedTrie(hashStack []string) {

	//hashStack = hashStack[:len(hashStack)-1]
	fmt.Println("HASHSTACK")
	for i := range hashStack {
		fmt.Println(hashStack[i])
	}
	counter := len(hashStack) - 1
	rearranged := mpt.rearrangeDeletedTrieHelper(hashStack, counter, "")
	fmt.Println("rearranged : ", rearranged)

}

func (mpt *MerklePatriciaTrie) rearrangeDeletedTrieHelper(hashStack []string, counter int, currenthash string) bool {
	//shoud exit after  : counter >= len(hashStack)
	fmt.Println("rearrangeDeletedTrieHelper:Counter =", counter)

	if counter == -1 {
		mpt.root = currenthash
		return true
	}

	if counter == len(hashStack)-1 {
		fmt.Println("First loop")
		if len(hashStack) == 1 {
			if mpt.db[hashStack[counter]].node_type == 1 { //just one node in stack and should be a branch
				currNode := mpt.db[hashStack[counter]]
				numValues := 0
				n := 0 //mpt.db[currNode.branch_value]
				for i := 0; i < 17; i++ {
					if currNode.branch_value[i] != "" {
						numValues++
						if numValues == 1 {
							n = i
							// 	nextnode := mpt.db[currNode.branch_value[i]]
						}
					}
				}
				if numValues == 1 {
					nextnode := mpt.db[currNode.branch_value[n]]
					nodetype := nextnode.node_type
					if nodetype == 1 {
						nodeE := Node{}
						nodeE.node_type = 2
						nodeE.flag_value.encoded_prefix = compact_encode([]uint8{uint8(n)})
						delete(mpt.db, hashStack[counter])
						hash := nodeE.hash_node()
						mpt.db[hash] = nodeE
						return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
					} else if nodetype == 2 {
						hex_prefix_array := AsciiArrayToHexArray(nextnode.flag_value.encoded_prefix)
						if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) { //extension
							nodeE := Node{}
							nodeE.node_type = 2
							combinedHexArray := append([]uint8{uint8(n)}, compact_decode(nextnode.flag_value.encoded_prefix)...)
							nodeE.flag_value.encoded_prefix = compact_encode(combinedHexArray)
							nodeE.flag_value.value = nextnode.flag_value.value
							delete(mpt.db, hashStack[counter])
							delete(mpt.db, hashStack[counter-1])
							hash := nodeE.hash_node()
							mpt.db[hash] = nodeE
							return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
						} else if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) { //leaf
							nodeL := Node{}
							nodeL.node_type = 2
							hex_prefix_array = append([]uint8{uint8(n)}, hex_prefix_array...)
							nodeL.flag_value.encoded_prefix = compact_encode(hex_prefix_array)
							nodeL.flag_value.value = nextnode.flag_value.value
							delete(mpt.db, hashStack[counter])
							hash := nodeL.hash_node()
							mpt.db[hash] = nodeL
							return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, nodeL.hash_node())
						}
					} else {
						return false
					}
				} else if numValues > 1 {
					node := mpt.db[hashStack[counter]]
					hash := node.hash_node()
					mpt.db[hash] = node //pepsi1
					return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
				}
			} else {
				fmt.Println("Not possible")
				return false
			}

		} else if len(hashStack) > 1 {
			if mpt.db[hashStack[counter]].node_type == 1 { //............................. branch
				currNode := mpt.db[hashStack[counter]]
				fmt.Println("Current Branch Node:", currNode)
				numValues := 0
				n := 0 //mpt.db[currNode.branch_value]
				for i := 0; i < 17; i++ {
					if currNode.branch_value[i] != "" {
						numValues++
						if numValues == 1 {
							n = i
							// 	nextnode := mpt.db[currNode.branch_value[i]]
						}
					}
				}
				fmt.Println("Number of values in the branch:", numValues)
				if numValues == 1 {
					fmt.Println("One value in branch")
					if n == 16 {
						fmt.Println("It is the branchvalue[16] field n=16")
						//convert to leaf and store value with key = ""
						//put the value in branch[16] of the following branch
						prevnode := mpt.db[hashStack[counter-1]]
						prevnodetype := prevnode.node_type
						if prevnodetype == 1 { // branch
							// TODO :- create a empty(key) leaf node with value (branch_value[16])
							fmt.Println("Previous is a branch")
							nodeL := Node{}
							nodeL.node_type = 2
							nodeL.flag_value.encoded_prefix = compact_encode([]uint8{16})
							nodeL.flag_value.value = currNode.branch_value[16]
							delete(mpt.db, hashStack[counter])
							hash := nodeL.hash_node()                                         //pepsi
							mpt.db[hash] = nodeL                                              //pepsi
							return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash) //?
						} else if prevnodetype == 2 {
							fmt.Println("Previous is a ext/leaf")
							hex_prefix_array := AsciiArrayToHexArray(prevnode.flag_value.encoded_prefix)
							fmt.Println("node_type == 2 - hex prefix: ", hex_prefix_array)
							if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) { //extension
								fmt.Println("Previous is a ext")
								// TODO :- add the value to prev Ext node - and make a leaf
								hex_array := compact_decode(prevnode.flag_value.encoded_prefix)
								hex_array = append(hex_array, 16)
								nodeL := Node{}
								nodeL.node_type = 2
								nodeL.flag_value.encoded_prefix = compact_encode(hex_array)
								nodeL.flag_value.value = currNode.branch_value[16]
								delete(mpt.db, hashStack[counter])
								delete(mpt.db, hashStack[counter-1])
								hash := nodeL.hash_node() //pepsi
								mpt.db[hash] = nodeL
								return mpt.rearrangeDeletedTrieHelper(hashStack, counter-2, hash) //?
							} else if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) { //leaf
								fmt.Println("Previous is a leaf")
								fmt.Println("should never reach here - 123")
								return false
							}
						}

					} else if n < 16 && n >= 0 {
						//hash value
						fmt.Println("the branch has just one hash ")
						nextnode := mpt.db[currNode.branch_value[n]]
						nodetype := nextnode.node_type
						u := uint8(n)

						if nodetype == 1 {
							fmt.Println("Next is a branch")
							//Creating Extension node
							nodeE := Node{}
							nodeE.node_type = 2
							nodeE.flag_value.encoded_prefix = compact_encode([]uint8{u})
							nodeE.flag_value.value = currNode.branch_value[n]
							delete(mpt.db, hashStack[counter])
							hash := nodeE.hash_node() //pepsi
							mpt.db[hash] = nodeE      //pepsi
							return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
						} else if nodetype == 2 {
							fmt.Println("Next is a ext/leaf")
							hex_prefix_array := AsciiArrayToHexArray(nextnode.flag_value.encoded_prefix)
							fmt.Println("node_type == 2 - hex prefix: ", hex_prefix_array)
							if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) { //extension
								fmt.Println("Next is a ext")
								var new_array []uint8
								previous_node := mpt.db[hashStack[counter-1]]
								if previous_node.node_type == 2 {
									hex_with_prefix_array := AsciiArrayToHexArray(previous_node.flag_value.encoded_prefix)
									if hex_with_prefix_array[0] == 0 || hex_with_prefix_array[0] == 1 {
										fmt.Println("Extension node..")
										new_array = compact_decode(previous_node.flag_value.encoded_prefix)
										fmt.Println("New array:", new_array)
									}
								}
								new_array = append(new_array, []uint8{u}...)
								fmt.Println("New array:", new_array)

								hex_array := compact_decode(nextnode.flag_value.encoded_prefix)
								new_array = append(new_array, hex_array...)
								fmt.Println("New array:", new_array)

								nextnode.flag_value.encoded_prefix = compact_encode(new_array)
								delete(mpt.db, hashStack[counter])   //pepsi maybe
								delete(mpt.db, hashStack[counter-1]) //pepsi maybe

								hash := nextnode.hash_node() //pepsi
								mpt.db[hash] = nextnode      //pepsi
								return mpt.rearrangeDeletedTrieHelper(hashStack, counter-2, hash)
							} else if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) { //leaf
								fmt.Println("Next is a leaf")
								//check if the previous is extension
								// if so club it with the leaf
								//now me me me me
								var hex_array []uint8
								nodeUnknown := mpt.db[hashStack[counter-1]]
								if nodeUnknown.node_type == 2 {
									hex_with_prefix_array := AsciiArrayToHexArray(nodeUnknown.flag_value.encoded_prefix)
									if hex_with_prefix_array[0] == 0 || hex_with_prefix_array[0] == 1 {
										fmt.Println("Extension node..")
										hex_array = compact_decode(nodeUnknown.flag_value.encoded_prefix)
									}
								}
								hex_array = append(hex_array, []uint8{u}...)

								to_be_added_array := compact_decode(nextnode.flag_value.encoded_prefix)
								hex_array = append(hex_array, to_be_added_array...)
								// hex_array = append([]uint8{u}, hex_array...)
								hex_array = append(hex_array, 16)
								fmt.Println("The final hex array1234567 : ", hex_array)
								nextnode.flag_value.encoded_prefix = compact_encode(hex_array)
								delete(mpt.db, hashStack[counter])   //pepsi maybe
								delete(mpt.db, hashStack[counter-1]) //pepsi
								hash := nextnode.hash_node()         //pepsi
								mpt.db[hash] = nextnode              //pepsi
								return mpt.rearrangeDeletedTrieHelper(hashStack, counter-2, hash)
							}
						}
					}
				} else if numValues > 1 {
					fmt.Println("Branch has more than one values in it")
					node := mpt.db[hashStack[counter]]
					////
					if node.node_type == 1 {
						for i := range node.branch_value {
							fmt.Println("values of branch : ", i, node.branch_value[i])
						}
					}
					///
					hash := node.hash_node()
					fmt.Println("mpt db len :", len(mpt.db))
					//pepsi added now
					delete(mpt.db, hashStack[counter])
					fmt.Println("mpt db len :", len(mpt.db))
					counter = counter - 1
					mpt.db[hash] = node
					return mpt.rearrangeDeletedTrieHelper(hashStack, counter, hash)
				}
			}
		}
	} else {
		fmt.Println("\nLOOP:", counter)
		node := mpt.db[hashStack[counter]]
		if node.node_type == 1 {
			fmt.Println("It is a branch")
			for i := 0; i < 17; i++ {
				if node.branch_value[i] == hashStack[counter+1] {
					node.branch_value[i] = currenthash
				}
			}
			delete(mpt.db, hashStack[counter]) //pepsi
			hash := node.hash_node()           //pepsi
			mpt.db[hash] = node                //pepsi
			return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
		} else if node.node_type == 2 {
			fmt.Println("It is a extension/leaf")
			node.flag_value.value = currenthash
			delete(mpt.db, hashStack[counter]) //pepsi
			hash := node.hash_node()           //pepsi
			mpt.db[hash] = node                //pepsi
			return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
		} else {
			fmt.Println("Node Type: ", node.node_type)
			fmt.Println("Node : ", node)
		}
	}
	fmt.Println("Final False")
	return false
}
