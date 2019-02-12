package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"

	"golang.org/x/crypto/sha3"
)

type Flag_value struct {
	encoded_prefix []uint8
	value          string
}

type Node struct {
	node_type    int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value   Flag_value
}

type MerklePatriciaTrie struct {
	db   map[string]Node
	root string
}

// func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {

// 	pathToCheck := StringToHexArray(key)
// 	fmt.Println("get - pathtoCheck : ", pathToCheck)
// 	return mpt.GetHelper(mpt.root, pathToCheck)

// 	//return "", errors.New("path_not_found")
// }

func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {

	// return "", errors.New("path_not_found")

	pathLeft := StringToHexArray(key)
	fmt.Println("input Hex:", pathLeft)
	currentNode := mpt.db[mpt.root]
	fmt.Printf("root node : %+v", currentNode)
	finalValue, err := mpt.getHelper(currentNode, pathLeft)

	//if finalNode.
	fmt.Println("finalValue:", finalValue)

	if err != nil {
		return "", err
	} else {
		return finalValue, nil
	}

}

func (mpt *MerklePatriciaTrie) getHelper(currentNode Node, pathLeft []uint8) (string, error) {

	if len(pathLeft) < 0 {
		//fmt.Println("pathleft is less than 0")
		return "", errors.New("path_not_found")

	} else if len(pathLeft) == 0 {
		//fmt.Println("pathleft is equal 0")
		if currentNode.node_type == 1 {
			//fmt.Println("pathleft is equal 0 - node type 1")
			if currentNode.branch_value[16] != "" {
				//say found and return
				return currentNode.branch_value[16], nil
			} else {
				return "", errors.New("path_not_found")
			}
		} else if currentNode.node_type == 2 {
			//fmt.Println("pathleft is equal 0 - node type 2")
			hex_prefix_array := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
			if hex_prefix_array[0] == 2 || hex_prefix_array[0] == 3 {
				//fmt.Println("path length 0 next node is leaf")
				if len(compact_decode(currentNode.flag_value.encoded_prefix)) == 0 {
					return currentNode.flag_value.value, nil
				} else {
					return "", errors.New("path_not_found")
				}
			}
		} else {
			return "", errors.New("path_not_found")
		}
	} else if len(pathLeft) > 0 {
		//fmt.Println("Current Node Type:", currentNode.node_type)
		//fmt.Println("Node:", currentNode)
		//fmt.Println("pathleft is greater than 0")
		if currentNode.node_type == 1 { //branch node
			//fmt.Println("pathleft is greater than 0 - branch node")
			//fmt.Println("pathLeft", pathLeft)
			//fmt.Println("currentNode.branch_value", currentNode.branch_value[5])
			if currentNode.branch_value[pathLeft[0]] != "" {
				//fmt.Println("Value exists in branch!!")
				//fmt.Println("Hash found:", currentNode.branch_value[pathLeft[0]])
				hash := currentNode.branch_value[pathLeft[0]]
				pathLeft = pathLeft[1:]
				//fmt.Println("Path Left from branch :", pathLeft)
				return mpt.getHelper(mpt.db[hash], pathLeft)
			} else {
				//fmt.Println("value not found in branch:", currentNode.branch_value[pathLeft[0]])
				//fmt.Println("Branch value at the index when not found")
				return "", errors.New("path_not_found")
			}
		} else if currentNode.node_type == 2 {
			//fmt.Println("pathleft is greater than 0 - node_type 2")
			hex_prefix_array := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
			//fmt.Println("!!!!!!!!!!!HEX PREFIX :", hex_prefix_array)
			if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) {
				//fmt.Println("pathleft is greater than 0 - node_type 2 - prefix 0 or 1 Extension")
				counter := 0
				triePath := compact_decode(currentNode.flag_value.encoded_prefix)
				//fmt.Println("currentNode.flag_value.encoded_prefix:", currentNode.flag_value.encoded_prefix)
				//fmt.Println("triePathtriePathtriePath:", triePath)
				//fmt.Println("pathLeftpathLeftpathLeft:", pathLeft)
				for i := 0; i < len(triePath); i++ {
					if triePath[i] == pathLeft[i] {
						counter++
					}
				}
				//fmt.Println("counter34", counter)
				if counter == (len(triePath)) {
					//fmt.Println("counter54", counter)
					pathLeft = pathLeft[counter:]
					//fmt.Println("Path left before call:", pathLeft)
					return mpt.getHelper(mpt.db[currentNode.flag_value.value], pathLeft)
				} else {
					return "", errors.New("path_not_found")
				}
			} else if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) {
				//fmt.Println("pathleft is greater than 0 - node_type 2 - prefix 2 or 3 Leaf")
				counter := 0
				triePath := compact_decode(currentNode.flag_value.encoded_prefix)
				//fmt.Println("leftPath :", pathLeft)
				//fmt.Println("triePath", triePath)
				for i := 0; i < len(triePath); i++ {
					if triePath[i] == pathLeft[i] {
						counter++
					}
				}
				//fmt.Println("Counter", counter)
				if counter == (len(triePath)) {
					pathLeft = pathLeft[counter:]
					if len(pathLeft) == 0 {
						return currentNode.flag_value.value, nil
					} else {
						return mpt.getHelper(mpt.db[currentNode.flag_value.value], pathLeft)
					}
				} else {
					return "", errors.New("path_not_found")
				}
				// return "", errors.New("path_not_found")
			}
		}
	}
	return "", errors.New("path_not_found")
}

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	//hex_key_array := encode(toByteArray(key))
	hex_key_array := StringToHexArray(key)
	if mpt.root == "" {
		fmt.Println("No root")
		hex_key_array = append(hex_key_array, 16)
		nodeL := Node{}
		nodeL.node_type = 2
		nodeL.flag_value.encoded_prefix = compact_encode(hex_key_array)
		fmt.Println("Encoded prefix :", nodeL.flag_value.encoded_prefix)
		nodeL.flag_value.value = new_value
		mpt.root = mpt.InsertRecursive(nodeL, nil, "", nil)
		fmt.Println("Inserted in root")
	} else {
		root_node_hash := mpt.root
		fmt.Println("Root hash before", mpt.root)
		mpt.root = mpt.InsertRecursive(mpt.db[root_node_hash], hex_key_array, new_value, []uint8{})
		fmt.Println("Root hash after", mpt.root)
	}

}

func (mpt *MerklePatriciaTrie) InsertRecursive(node Node, leftPath []uint8, newValue string, leftNibble []uint8) string {
	currNode := node
	//if leftPath == nil && newValue == "" {
	if len(leftPath) == 0 && newValue == "" {
		hash := currNode.hash_node()
		if mpt.db == nil {
			mpt.db = make(map[string]Node)
		}
		mpt.db[hash] = currNode
		return hash
	} else {
		if currNode.node_type == 2 {
			fmt.Println("Leaf or extension")
			path := currNode.flag_value.encoded_prefix
			//hexPrefixArray := ascii_to_hex(path)
			hexPrefixArray := AsciiArrayToHexArray(path)
			fmt.Println("hex with prefix:", hexPrefixArray)
			decodedHexArray := compact_decode(path)
			fmt.Println("decoded hex:", decodedHexArray)
			fmt.Println("Left path:", leftPath)
			//if leaf
			if hexPrefixArray[0] == 2 || hexPrefixArray[0] == 3 {
				//same
				fmt.Println("It is a leaf")
				//if check_if_equal(decodedHexArray, leftPath) {
				if reflect.DeepEqual(decodedHexArray, leftPath) {
					fmt.Println("Equal!!")
					currNode.flag_value.value = newValue
					leftPath = nil
					return mpt.InsertRecursive(currNode, leftPath, "", leftNibble)
				} else if check_if_common_path_exists(leftPath, decodedHexArray) { //common path exists
					// } else if leftPath[0] == decodedHexArray[0] { //common path exists

					fmt.Println("Common path exists")
					counter := 0
					for i := 0; i < len(decodedHexArray); i++ {
						if decodedHexArray[i] == leftPath[i] {
							counter++
						} else {
							break
						}
					}
					fmt.Println("Counter :", counter)
					common_path := leftPath[:(counter)]
					leftPath = leftPath[counter:]
					leftNibble = decodedHexArray[counter:]
					fmt.Println("Common Path :", common_path)
					fmt.Println("Left path :", leftPath)
					fmt.Println("Left nibble :", leftNibble)
					//create extension node and branch
					nodeE := Node{}
					nodeE.node_type = 2
					nodeE.flag_value.encoded_prefix = compact_encode(common_path)
					//create branch
					nodeB := Node{}
					nodeB.node_type = 1

					// if len(leftPath) > 0 && len(leftNibble) > 0 {
					//create two leaves
					leftPath = append(leftPath, 16)
					leftNibble = append(leftNibble, 16)
					fmt.Println("Left path appended :", leftPath)
					fmt.Println("Left nibble appended:", leftNibble)
					nodeL1 := Node{}
					nodeL1.node_type = 2
					nodeL1.flag_value.encoded_prefix = compact_encode(leftPath[1:])
					nodeL1.flag_value.value = newValue
					if len(leftNibble) > 1 {
						nodeL2 := Node{}
						nodeL2.node_type = 2
						nodeL2.flag_value.encoded_prefix = compact_encode(leftNibble[1:])
						nodeL2.flag_value.value = currNode.flag_value.value
						nodeB.branch_value[leftNibble[0]] = mpt.InsertRecursive(nodeL2, nil, "", nil)
					} else if len(leftNibble) == 1 {
						nodeB.branch_value[16] = currNode.flag_value.value
						fmt.Println("Inserted in the value field of the branch")
					}
					nodeB.branch_value[leftPath[0]] = mpt.InsertRecursive(nodeL1, nil, "", nil)
					fmt.Println("Inserted nodeL1")
					fmt.Println("Inserted nodeL2")

					nodeE.flag_value.value = mpt.InsertRecursive(nodeB, nil, "", nil)
					fmt.Println("Inserted nodeLE")

					hashE := nodeE.hash_node()
					mpt.db[hashE] = nodeE
					return hashE
				} else {
					//create branch
					//create leaves or leaf
					fmt.Println("In leaf and the left path is totally different than a leaf")
					nodeB := Node{}
					nodeB.node_type = 1
					nodeL1 := Node{}
					nodeL1.node_type = 2
					if len(leftPath) > 0 {
						leftPath = append(leftPath, 16)
						nodeL1.flag_value.encoded_prefix = compact_encode(leftPath[1:])
					} else {
						leftPath = append([]uint8{16})
						nodeL1.flag_value.encoded_prefix = compact_encode(leftPath[1:])
						// hash := currNode.hash_node()
						// mpt.db[hash] = currNode
						// return hash
					}

					fmt.Println("leftPath :", leftPath)
					//pepsi
					// newleftpath := leftPath[1:]
					// fmt.Println("newleftpath : ", newleftpath)
					// if len(newleftpath) == 0 {
					// 	// cocacola
					// 	nodeL1.flag_value.encoded_prefix = compact_encode([]uint8{16})
					// } else {
					// 	nodeL1.flag_value.encoded_prefix = compact_encode(newleftpath)
					// }
					//fmt.Println("!!!!!! After newleftpath : ", newleftpath)

					nodeL1.flag_value.value = newValue
					nodeB.branch_value[leftPath[0]] = mpt.InsertRecursive(nodeL1, nil, "", nil)

					if len(decodedHexArray) > 0 {
						// if len(hexPrefixArray) > 1 && len(leftPath) > 1 {
						// leftPath = append(leftPath, 16)
						decodedHexArray = append(decodedHexArray, 16)
						// nodeL1 := Node{}
						// nodeL1.node_type = 2
						// fmt.Println("leftPath :", leftPath)
						// nodeL1.flag_value.encoded_prefix = compact_encode(leftPath[1:])
						// nodeL1.flag_value.value = newValue
						nodeL2 := Node{}
						nodeL2.node_type = 2
						fmt.Println("decodedHexArray :", decodedHexArray)
						nodeL2.flag_value.encoded_prefix = compact_encode(decodedHexArray[1:])
						nodeL2.flag_value.value = currNode.flag_value.value
						nodeB.branch_value[decodedHexArray[0]] = mpt.InsertRecursive(nodeL2, nil, "", nil)
						// nodeB.branch_value[leftPath[0]] = mpt.InsertRecursive(nodeL1, nil, "", nil)
						// hashB := nodeB.hash_node()
						// mpt.db[hashB] = nodeB
						// fmt.Println("Branch and two leaves inserted")
						// return hashB
					} else if len(decodedHexArray) == 0 {
						nodeB.branch_value[16] = currNode.flag_value.value
					}
					hashB := nodeB.hash_node()
					mpt.db[hashB] = nodeB
					fmt.Println("Branch and two leaves inserted")
					return hashB
				}
			} else if hexPrefixArray[0] == 0 || hexPrefixArray[0] == 1 {
				// if equal
				fmt.Println("Extension")
				//if check_if_equal(decodedHexArray, leftPath) {
				if reflect.DeepEqual(decodedHexArray, leftPath) {
					//insert in branch value place
					fmt.Println("Equal in extension..")
					fmt.Println("Left path:", leftPath)
					currNode.flag_value.value = mpt.InsertRecursive(mpt.db[currNode.flag_value.value], nil, newValue, nil)
				} else if check_if_common_path_exists(leftPath, decodedHexArray) {
					counter := 0
					for i := 0; i < len(decodedHexArray); i++ {
						if decodedHexArray[i] == leftPath[i] {
							counter++
						} else {
							break
						}
					}
					fmt.Println("Counter:", counter)
					common_path2 := leftPath[:counter]
					leftPath = leftPath[counter:]
					leftNibble = decodedHexArray[counter:]
					currNode.flag_value.encoded_prefix = compact_encode(common_path2)
					fmt.Println("common_path2:", common_path2)
					fmt.Println("leftPath:", leftPath)
					fmt.Println("leftNibble:", leftNibble)
					nodeBranch := Node{}
					nodeBranch.node_type = 1
					if len(leftNibble) > 0 {
						//create Branch
						//left path create branch and leaf
						leftPath = append(leftPath, 16)

						nodeLeaf := Node{}
						nodeLeaf.node_type = 2
						nodeLeaf.flag_value.encoded_prefix = compact_encode(leftPath[1:])
						nodeLeaf.flag_value.value = newValue
						fmt.Println("Inserting in branch at index:", leftPath[0])
						nodeBranch.branch_value[leftPath[0]] = mpt.InsertRecursive(nodeLeaf, nil, "", nil)
					}
					//check the left nibble size
					if len(leftNibble) == 0 {
						currNode.flag_value.value = mpt.InsertRecursive(mpt.db[currNode.flag_value.value], leftPath, newValue, nil)
						hashExt := currNode.hash_node()
						mpt.db[hashExt] = currNode
						return hashExt
					} else if len(leftNibble) == 1 {
						nodeBranch.branch_value[leftNibble[0]] = currNode.flag_value.value
						currNode.flag_value.value = mpt.InsertRecursive(nodeBranch, nil, "", nil)
						hashExt := currNode.hash_node()
						mpt.db[hashExt] = currNode
						return hashExt
					} else if len(leftNibble) > 1 {
						//create an extension and store the value of that extension in branch => branch hash store in currNode
						nodeExtension := Node{}
						nodeExtension.node_type = 2
						nodeExtension.flag_value.encoded_prefix = compact_encode(leftNibble[1:])
						nodeExtension.flag_value.value = currNode.flag_value.value
						fmt.Println("Inserting in branch at index:", leftNibble[0])
						nodeBranch.branch_value[leftNibble[0]] = mpt.InsertRecursive(nodeExtension, nil, "", nil)
						currNode.flag_value.value = mpt.InsertRecursive(nodeBranch, nil, "", nil)
						hashExt := currNode.hash_node()
						mpt.db[hashExt] = currNode
						return hashExt
					}
					// currNode.flag_value.value = mpt.InsertRecursive(nodeBranch, nil, "", nil)
					// hashExt := currNode.hash_node()
					// mpt.db[hashExt] = currNode
					// return hashExt
				} else {
					//make extension node a branch and one leaf
					nodeBranch := Node{}
					nodeBranch.node_type = 1

					leftPath = append(leftPath, 16)
					//creating a leaf
					nodeLeaf := Node{}
					nodeLeaf.node_type = 2
					nodeLeaf.flag_value.encoded_prefix = compact_encode(leftPath[1:])
					nodeLeaf.flag_value.value = newValue
					nodeBranch.branch_value[leftPath[0]] = mpt.InsertRecursive(nodeLeaf, nil, "", nil)

					if len(decodedHexArray) == 1 {
						nodeBranch.branch_value[decodedHexArray[0]] = currNode.flag_value.value
					} else if len(decodedHexArray) > 1 {
						nodeExt := Node{}
						nodeExt.node_type = 2
						nodeExt.flag_value.encoded_prefix = compact_encode(decodedHexArray[1:])
						nodeExt.flag_value.value = currNode.flag_value.value
						nodeBranch.branch_value[decodedHexArray[0]] = mpt.InsertRecursive(nodeExt, nil, "", nil)
					}
					hashBranch := nodeBranch.hash_node()
					mpt.db[hashBranch] = nodeBranch
					return hashBranch
				}
			}
		} else if currNode.node_type == 1 {
			if currNode.branch_value[leftPath[0]] == "" {
				//store leftPath[0] create a leaf to store the rest
				fmt.Println("Branch and it is new entry")
				leftPath = append(leftPath, 16)
				nodeL3 := Node{}
				nodeL3.node_type = 2
				nodeL3.flag_value.encoded_prefix = compact_encode(leftPath[1:])
				nodeL3.flag_value.value = newValue
				currNode.branch_value[leftPath[0]] = mpt.InsertRecursive(nodeL3, nil, "", nil)
			} else if currNode.branch_value[leftPath[0]] != "" {
				fmt.Println("Branch and path already exists")
				index := leftPath[0]
				nextNode := mpt.db[currNode.branch_value[leftPath[0]]]
				leftPath = leftPath[1:]

				currNode.branch_value[index] = mpt.InsertRecursive(nextNode, leftPath, newValue, nil)
			}
			hashBr := currNode.hash_node()
			mpt.db[hashBr] = currNode
			return hashBr
		}
	}
	return ""
}

func check_if_common_path_exists(leftPath []uint8, decodedHexArray []uint8) bool {
	if len(leftPath) > 0 && len(decodedHexArray) > 0 {
		if leftPath[0] == decodedHexArray[0] {
			return true
		}
	}
	return false
}

//createLeafNode
func createLeafNode(path []uint8, new_value string) (string, Node) {
	flagStruct := Flag_value{path, new_value}
	n1 := &Node{
		node_type:  2,
		flag_value: flagStruct,
	}
	nHash := n1.hash_node()
	return nHash, *n1
}

//createExtensionNode
func createExtensionNode(path []uint8, next_hash_value string) (string, Node) {
	flagStruct := Flag_value{path, next_hash_value}
	n1 := &Node{
		node_type:  2,
		flag_value: flagStruct,
	}
	nHash := n1.hash_node()
	return nHash, *n1
}

//createBranchNode
func createBranchNode(path []uint8, new_value string) (string, Node) {
	n1 := &Node{
		node_type:    1,
		branch_value: [17]string{},
	}
	nHash := n1.hash_node()
	return nHash, *n1
}

//Delete func
func (mpt *MerklePatriciaTrie) Delete(key string) string {

	pathLeft := StringToHexArray(key)
	fmt.Println("in Del - input Hex:", pathLeft)
	currentNode := mpt.db[mpt.root]
	fmt.Printf("in Del - root node : %+v", currentNode)
	//sendinf root hash in hashstack
	value, err := mpt.delHelper(mpt.root, currentNode, pathLeft, []string{mpt.root})
	fmt.Println("in Del -deleted value: ", value)
	if err != nil {
		return value
	}
	return ""

}

func (mpt *MerklePatriciaTrie) delHelper(nodeKey string, currentNode Node, pathLeft []uint8, hashStack []string) (string, error) {
	fmt.Println("--------------------------------------------------------------------------------------")
	if len(pathLeft) > 0 && currentNode.node_type != 0 { //path length >0
		fmt.Println("in Del - path length > 0")
		if currentNode.node_type == 1 { // branch and pathleft >0
			fmt.Println("in Del - path length > 0 - node is Branch")
			if currentNode.branch_value[pathLeft[0]] != "" {
				fmt.Println("Value exists in branch!!")
				fmt.Println("Hash found:", currentNode.branch_value[pathLeft[0]])
				hash := currentNode.branch_value[pathLeft[0]]
				// pathLeft = pathLeft[1:]
				fmt.Println("PathLeft from branch :", pathLeft)
				// if path ends in branch and so value can be  stored in empty leaf node
				if len(pathLeft) == 0 { // pathleft gets over at branch node
					fmt.Println("pathleft gets over at branch node")
					if mpt.db[hash].node_type == 1 { // if this branch node contains a value of node
						fmt.Println("next node :-  branch node")
						if mpt.db[hash].branch_value[16] != "" {
							fmt.Println("next node :-  branch node - value in pos 16")
							hashStack = append(hashStack, hash)
							branchnode := mpt.db[hash]
							value := branchnode.branch_value[16]
							branchnode.branch_value[16] = ""
							fmt.Println("Deleting from db")
							delete(mpt.db, hash) // delete :-> empty - value node
							///  l!!!!!!!!! rearrange trie !!!!!
							fmt.Println("rearranging trie : ")
							mpt.rearrangeDeletedTrie(hashStack)
							////// !!!!!!!!!!!!!!!!!!                 !!!!!!!!
							fmt.Println("Returning with value : ", value)
							return value, nil
						} else {
							fmt.Println("Returning with Error")
							return "", errors.New("path_not_found")
						}
					} else if mpt.db[hash].node_type == 2 {
						fmt.Println("next node :-  Ext/Leaf node")
						hex_prefix_array := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
						fmt.Println("node_type == 2 - hex prefix: ", hex_prefix_array)
						if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) { //extension
							fmt.Println("pathleft is greater than 0 - node_type 2 - prefix 0 or 1 Extension")
							return "", errors.New("path_not_found")
						} else if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) { //leaf
							fmt.Println("pathleft is greater than 0 - node_type 2 - prefix 2 or 3 Leaf")
							value := mpt.db[hash].flag_value.value
							// if this ext/leaf node contains a value of node
							// encoded prefix should be empty , if so then ////
							fmt.Println("Deleting from db")
							delete(mpt.db, hash) // delete :-> empty - value node
							///  l!!!!!!!!! rearrange trie !!!!!
							fmt.Println("rearranging trie : ")
							mpt.rearrangeDeletedTrie(hashStack)
							////// !!!!!!!!!!!!!!!!!!                 !!!!!!!!
							fmt.Println("Returning with value : ", value)
							return value, nil
						} else {
							fmt.Println("Returning with Error : ")
							return "", errors.New("path_not_found")
						}
					} else {
						fmt.Println("Returning with Error : ")
						return "", errors.New("path_not_found")
					}
				} else { //pathLeft > 0
					fmt.Println("pathleft traversing through branch node")
					fmt.Println("currentNode.branch_value before making empty", currentNode.branch_value[pathLeft[0]])
					fmt.Println("\nHash before removing value from the branch[16] field:", currentNode.hash_node())
					currentNode.branch_value[pathLeft[0]] = ""
					fmt.Println("\nHash before removing value from the branch[16] field:", currentNode.hash_node())
					fmt.Println("\nHASHSTACK before : ", hashStack)
					hashStack = append(hashStack, hash)
					fmt.Println("HASHSTACK after: ", hashStack)
					pathLeft = pathLeft[1:]
					fmt.Println("Pathleft : ", pathLeft)
					return mpt.delHelper(hash, mpt.db[hash], pathLeft, hashStack)
				}
			} else {
				fmt.Println("returning with error")
				return "", errors.New("path_not_found")
			}
		} else if currentNode.node_type == 2 { //ext or leaf and pathleft >0
			hex_prefix_array := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
			fmt.Println("node_type == 2 - hex prefix: ", hex_prefix_array)
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
						hashStack = append(hashStack, hash)
						fmt.Println("HASHSTACK2 after: ", hashStack)
						//pathleft=0, use value from ext and check in branchvalu[16]
						if mpt.db[hash].branch_value[16] != "" { //value found
							fmt.Println("node Branch = valu at 16")
							//value := mpt.db[hash].branch_value[16]
							//delete node
							fmt.Println("so finding next node")
							nextnode := mpt.db[hash]
							fmt.Println("next node : ", nextnode)
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
							fmt.Println("returning : ", currentNode.flag_value.value)
							return currentNode.flag_value.value, nil
						} //................................value not found
						fmt.Println("returning with error")
						return "", errors.New("path_not_found")
						//"pathleft is greater than 0 - node_type 2 - prefix 1 or 2 Extension")
					} else if len(pathLeft) > 0 { //add to hashstack call on next node
						fmt.Println("leftpath >0, ")
						hash := currentNode.flag_value.value
						fmt.Println("HASHSTACK before : ", hashStack)
						hashStack = append(hashStack, hash) //adding current ext
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
					fmt.Println("deep equal success")
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
	} // else if len(pathLeft) == 0 && currentNode.node_type != 0 { //pathlength ==0
	// 	if currentNode.node_type == 1 { // branch
	// 		//go to the
	// 	} else if currentNode.node_type == 2 { //ext or leaf
	// 		//extension
	// 		if currentNode.flag_value.encoded_prefix[0] == 0 || currentNode.flag_value.encoded_prefix[0] == 1 {

	// 			//leaf
	// 		} else if currentNode.flag_value.encoded_prefix[0] == 2 || currentNode.flag_value.encoded_prefix[0] == 3 {

	// 		}
	// 	}

	// } else if currentNode.node_type == 0 || len(pathLeft) < 0 { //path length <0 or nodetype =0
	// 	return "", errors.New("path_not_found")
	// }
	fmt.Println("EXit with Err or")
	return "", errors.New("path_not_found")
}

// rearrangeDeletedTrie rearranges mpt
func (mpt *MerklePatriciaTrie) rearrangeDeletedTrie(hashStack []string) {

	hashStack = hashStack[:len(hashStack)-1]
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
							fmt.Println("Previous is a branch")
							//Creating Extension node
							nodeE := Node{}
							nodeE.node_type = 2
							nodeE.flag_value.encoded_prefix = compact_encode([]uint8{u})
							nodeE.flag_value.value = currNode.branch_value[n]
							delete(mpt.db, hashStack[counter])
							hash := nextnode.hash_node() //pepsi
							mpt.db[hash] = nextnode      //pepsi
							return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
						} else if nodetype == 2 {
							fmt.Println("Previous is a ext/leaf")
							hex_prefix_array := AsciiArrayToHexArray(nextnode.flag_value.encoded_prefix)
							fmt.Println("node_type == 2 - hex prefix: ", hex_prefix_array)
							if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) { //extension
								fmt.Println("Previous is a ext")
								hex_array := compact_decode(nextnode.flag_value.encoded_prefix)
								hex_array = append([]uint8{u}, hex_array...)
								nextnode.flag_value.encoded_prefix = compact_encode(hex_array)
								delete(mpt.db, hashStack[counter]) //pepsi maybe
								hash := nextnode.hash_node()       //pepsi
								mpt.db[hash] = nextnode            //pepsi
								return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
							} else if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) { //leaf
								fmt.Println("Previous is a leaf")
								hex_array := compact_decode(nextnode.flag_value.encoded_prefix)
								hex_array = append([]uint8{u}, hex_array...)
								hex_array = append(hex_array, 16)
								nextnode.flag_value.encoded_prefix = compact_encode(hex_array)
								delete(mpt.db, hashStack[counter]) //pepsi maybe
								hash := nextnode.hash_node()       //pepsi
								mpt.db[hash] = nextnode            //pepsi
								return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
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

func compact_encode(hex_array []uint8) []uint8 {
	var term int
	if hex_array[len(hex_array)-1] == 16 { //checking if last element in array is 16
		term = 1 // if last element is 16, term = 1 and remove the last element from hex_array
		hex_array = hex_array[:len(hex_array)-1]
	}
	oddlen := len(hex_array) % 2               //checking if length is odd (oddlen = 1) otherwise (oddlen = 0)
	flags := []uint8{(uint8(2*term + oddlen))} //calculating flags value
	//changing hex_array based on value of oddlen
	if oddlen == 1 { //odd       // prefix -> flags so... ( either, 1 - Odd length extension
		hex_array = append(flags, hex_array...) //      3 - Odd length leaf        )
	} else { //when oddlen = 0 // even
		flags = append(flags, uint8(0))
		hex_array = append(flags, hex_array...) //prefix -> flags + 0 so...  (either, 00 -  Even Length Extention)
	} // or 20 - Even length Leaf
	encoded_arr := []uint8{} //array to return
	for i := 0; i < len(hex_array); i += 2 {
		encoded_arr = append(encoded_arr, 16*hex_array[i]+hex_array[i+1])
	}
	//fmt.Println("ascii array :", ascii_array)
	return encoded_arr
} // closing compact_encode

func compact_decode(encoded_arr []uint8) []uint8 {
	decoded_arr := AsciiArrayToHexArray(encoded_arr) //converting ascii array to hex array

	prefix := decoded_arr[0]
	switch prefix {
	case 0:
		decoded_arr = decoded_arr[2:]
	case 1:
		decoded_arr = decoded_arr[1:]
	case 2:
		decoded_arr = decoded_arr[2:]
	case 3:
		decoded_arr = decoded_arr[1:]
	}
	return decoded_arr
} //closing compact_decode

func test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value + string(node.flag_value.encoded_prefix)
	}
	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

// StringToHexArray takes in string as input and returns a splitted-hex(0-15) array : 1st
func StringToHexArray(s string) []uint8 {
	uint8Array := []uint8(s)
	hexArray := AsciiArrayToHexArray(uint8Array)
	return hexArray
}

// AsciiArrayToHexArray takes in string as input and returns a splitted-hex(0-15) array : 1st
func AsciiArrayToHexArray(encoded_arr []uint8) []uint8 {
	hexArray := []uint8{}
	for _, element := range encoded_arr {
		div, mod := element/16, element%16
		hexArray = append(hexArray, div)
		hexArray = append(hexArray, mod)
	}
	return hexArray
}

func main() { //before pepsi
	mpt := &MerklePatriciaTrie{make(map[string]Node), ""}
	mpt.CreateTestMpt()
	//fmt.Printf("%+v\n", mpt)
	// fmt.Printf("mpt-blank: %+v\n", mpt)
	// ////insert do -> verb
	// mpt.Insert("do", "verb")
	// fmt.Printf("mpt-do: %+v\n", mpt)
	////////////
	fmt.Println("####################################")
	fmt.Println("====================================")
	v, e := mpt.Get("dog")
	if e != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("value1 : ", v)
	}

	fmt.Println("====================================")
	fmt.Println("====================================")

	value := mpt.Delete("dog")
	fmt.Println("value2 : ", value)

	fmt.Println("====================================")
	fmt.Println("====================================")

	v1, e1 := mpt.Get("dog")
	if e1 != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("value3 : ", v1)
	}

	fmt.Println("====================================")
	v2, e2 := mpt.Get("do")
	if e2 != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("do : ", v2)
	}
	v1, e1 = mpt.Get("dog")
	if e1 != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("dog : ", v1)
	}
	v1, e1 = mpt.Get("doge")
	if e1 != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("doge : ", v1)
	}
	v1, e1 = mpt.Get("horse")
	if e1 != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("horse : ", v1)
	}
	fmt.Println("====================================")

	fmt.Println("inserting horse")
	mpt.Insert("horse", "stallion")
	v1, e1 = mpt.Get("horse")
	if e1 != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("horse : ", v1)
	}

	fmt.Println("inserting dog")
	mpt.Insert("dog", "puppy")

	vd, ed := mpt.Get("dog")
	if ed != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("dog : ", vd)
	}

	fmt.Println("EOP")
	fmt.Println("====================================")
	fmt.Println("####################################")

}

//main to insert node
// func main() {
// 	////initialize mpt to nil
// 	mpt := &MerklePatriciaTrie{make(map[string]Node), ""}
// 	fmt.Printf("mpt-blank: %+v\n", mpt)
// 	////insert do -> verb
// 	mpt.Insert("do", "verb")
// 	fmt.Printf("mpt-do: %+v\n", mpt)
// 	////insert dog -> puppy
// 	mpt.Insert("dog", "puppy")
// 	fmt.Printf("mpt-dog: %+v\n", mpt)

// }

//main to test encode decode
// func main() {
// 	bla := StringToHexArray("dog")
// 	//bla = append(bla, uint8(16))
// 	fmt.Println("bla :", bla)

// 	bla1 := compact_encode(bla)
// 	fmt.Println("bla1 :", bla1)

// 	bla2 := compact_decode(bla1)
// 	fmt.Println("bla2 :", bla2)

// abla := []uint8{1, 6, 1}
//abla := []uint8{32, 100, 111}
//abla := []uint8{6, 4, 6, 15, 16}
//abla1 := compact_encode(abla)
//fmt.Println("abla1=ascii :", abla1)
// abla2 := compact_decode(abla1)
// fmt.Println("abla2=hex : ", abla2)

// 	//root node
// 	// flagStruct := Flag_value{bla1, "verb"}
// 	// n1 := Node{
// 	// 	node_type:  2,
// 	// 	flag_value: flagStruct,
// 	// }
// 	// hashVal := n1.hash_node() //hash to put in db !!!!
// 	// fmt.Println("hashVal", hashVal)

// }

func (mpt *MerklePatriciaTrie) CreateTestMpt() error {
	mpt.db = make(map[string]Node)

	// 0: Null, 1: Branch, 2: Ext or Leaf
	flagValueNodeG := Flag_value{
		encoded_prefix: compact_encode([]uint8{5, 16}),
		value:          "coin",
	}

	nodeG := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeG,
	}
	hashNodeG := nodeG.hash_node()
	mpt.db[hashNodeG] = nodeG

	fmt.Println("Compact_encode([]uint8{5, 16}:", compact_encode([]uint8{5, 16}))

	flagValueNodeF := Flag_value{
		encoded_prefix: nil,
		value:          "puppy",
	}

	nodeF := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeF,
		branch_value: [17]string{"", "", "", "", "", "", hashNodeG, "", "", "", "", "", "", "", "", "", "puppy"},
	}
	hashNodeF := nodeF.hash_node()
	mpt.db[hashNodeF] = nodeF

	flagValueNodeE := Flag_value{
		encoded_prefix: compact_encode([]uint8{7}),
		value:          hashNodeF,
	}

	nodeE := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeE,
	}
	hashNodeE := nodeE.hash_node()
	mpt.db[hashNodeE] = nodeE

	flagValueNodeD := Flag_value{
		encoded_prefix: nil,
		value:          "verb",
	}

	nodeD := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeD,
		branch_value: [17]string{"", "", "", "", "", "", hashNodeE, "", "", "", "", "", "", "", "", "", "verb"},
	}
	hashNodeD := nodeD.hash_node()
	mpt.db[hashNodeD] = nodeD

	flagValueNodeB := Flag_value{
		encoded_prefix: compact_encode([]uint8{6, 15}),
		value:          hashNodeD,
	}

	nodeB := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeB,
	}
	hashNodeB := nodeB.hash_node()
	mpt.db[hashNodeB] = nodeB

	flagValueNodeC := Flag_value{
		encoded_prefix: compact_encode([]uint8{6, 15, 7, 2, 7, 3, 6, 5, 16}),
		value:          "", //"stallion",
	}

	nodeC := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeC,
	}
	hashNodeC := nodeC.hash_node()
	mpt.db[hashNodeC] = nodeC

	flagValueNodeA := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}

	nodeA := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeA,
		branch_value: [17]string{"", "", "", "", hashNodeB, "", "", "", hashNodeC, "", "", "", "", "", "", "", ""},
	}
	hashNodeA := nodeA.hash_node()
	mpt.db[hashNodeA] = nodeA

	flagValueRoot := Flag_value{
		encoded_prefix: compact_encode([]uint8{6}),
		value:          hashNodeA,
	}

	nodeRoot := Node{
		node_type:  2, //Extension Root
		flag_value: flagValueRoot,
	}
	hashNodeRoot := nodeRoot.hash_node()
	mpt.db[hashNodeRoot] = nodeRoot
	mpt.root = hashNodeRoot

	////////////////////////////////////////

	//Add another node to MPT

	//mpt.root = Compact_encode(key) //a
	return errors.New("Problem occured while creating Root Node")
}

func (mpt *MerklePatriciaTrie) CreateTestMpt2() error {
	mpt.db = make(map[string]Node)

	// 0: Null, 1: Branch, 2: Ext or Leaf
	flagValueNodeG := Flag_value{
		encoded_prefix: compact_encode([]uint8{5, 16}),
		value:          "coin",
	}

	nodeG := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeG,
	}
	hashNodeG := nodeG.hash_node()
	mpt.db[hashNodeG] = nodeG

	fmt.Println("Compact_encode([]uint8{5, 16}:", compact_encode([]uint8{5, 16}))

	flagValueNodeF := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}

	nodeF := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeF,
		branch_value: [17]string{"", "", "", "", "", "", hashNodeG, "", "", "", "", "", "", "", "", "", "puppy"},
	}
	hashNodeF := nodeF.hash_node()
	mpt.db[hashNodeF] = nodeF

	flagValueNodeE := Flag_value{
		encoded_prefix: compact_encode([]uint8{7}),
		value:          hashNodeF,
	}

	nodeE := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeE,
	}
	hashNodeE := nodeE.hash_node()
	mpt.db[hashNodeE] = nodeE

	flagValueNodeJ := Flag_value{
		encoded_prefix: compact_encode([]uint8{8}),
		value:          "book",
	}

	nodeJ := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeJ,
	}
	hashNodeJ := nodeJ.hash_node()
	mpt.db[hashNodeJ] = nodeJ

	flagValueNodeH := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}

	nodeH := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeH,
		branch_value: [17]string{"", "", hashNodeJ, "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
	}
	hashNodeH := nodeH.hash_node()
	mpt.db[hashNodeH] = nodeH

	flagValueNodeD := Flag_value{
		encoded_prefix: nil,
		value:          "verb",
	}

	nodeD := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeD,
		branch_value: [17]string{"", "", hashNodeH, "", "", "", hashNodeE, "", "", "", "", "", "", "", "", "", "verb"},
	}
	hashNodeD := nodeD.hash_node()
	mpt.db[hashNodeD] = nodeD

	flagValueNodeB := Flag_value{
		encoded_prefix: compact_encode([]uint8{6, 15}),
		value:          hashNodeD,
	}

	nodeB := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeB,
	}
	hashNodeB := nodeB.hash_node()
	mpt.db[hashNodeB] = nodeB

	flagValueNodeC := Flag_value{
		encoded_prefix: compact_encode([]uint8{6, 15, 7, 2, 7, 3, 6, 5, 16}),
		value:          "stallion",
	}

	nodeC := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeC,
	}
	hashNodeC := nodeC.hash_node()
	mpt.db[hashNodeC] = nodeC

	flagValueNodeA := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}

	nodeA := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeA,
		branch_value: [17]string{"", "", "", "", hashNodeB, "", "", "", hashNodeC, "", "", "", "", "", "", "", ""},
	}
	hashNodeA := nodeA.hash_node()
	mpt.db[hashNodeA] = nodeA

	flagValueRoot := Flag_value{
		encoded_prefix: compact_encode([]uint8{6}),
		value:          hashNodeA,
	}

	nodeRoot := Node{
		node_type:  2, //Extension Root
		flag_value: flagValueRoot,
	}
	hashNodeRoot := nodeRoot.hash_node()
	mpt.db[hashNodeRoot] = nodeRoot
	mpt.root = hashNodeRoot

	////////////////////////////////////////

	//Add another node to MPT

	//mpt.root = Compact_encode(key) //a
	return errors.New("Problem occured while creating Root Node")
}
