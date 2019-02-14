package p1

import (
	"fmt"
	"reflect"
)

func check_if_common_path_exists(leftPath []uint8, decodedHexArray []uint8) bool {
	if len(leftPath) > 0 && len(decodedHexArray) > 0 {
		if leftPath[0] == decodedHexArray[0] {
			return true
		}
	}
	return false
}

func (mpt *MerklePatriciaTrie) InsertRecursive(node Node, leftPath []uint8, newValue string, leftNibble []uint8) string {
	currNode := node
	old_hash := currNode.hash_node()
	//if leftPath == nil && newValue == "" {
	if len(leftPath) == 0 && newValue == "" {
		hash := currNode.hash_node()
		if mpt.db == nil {
			mpt.db = make(map[string]Node)
		}
		mpt.db[hash] = currNode
		return hash
	} else if len(leftPath) == 0 {
		fmt.Println("LEFT PATH IS ZERO TODO...")
		//lib
		fmt.Println("currNode", currNode)
		if currNode.node_type == 1 {
			currNode.branch_value[16] = newValue
			hash := currNode.hash_node()
			delete(mpt.db, old_hash)
			mpt.db[hash] = currNode
			return hash
		} else {
			fmt.Println("This should not be printed this is an error when the left path lenght is zero and the next node is a leaf or extensio...n")
		}
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
					delete(mpt.db, old_hash)
					return mpt.InsertRecursive(currNode, leftPath, "", leftNibble)
				} else if check_if_common_path_exists(leftPath, decodedHexArray) { //common path exists
					// } else if leftPath[0] == decodedHexArray[0] { //common path exists

					fmt.Println("Common path exists")
					counter := 0
					for i := 0; i < len(decodedHexArray); i++ {
						//now me me me
						if i == len(decodedHexArray) || i == len(leftPath) {
							break
						}
						if decodedHexArray[i] == leftPath[i] {
							counter = counter + 1
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
					if len(leftPath) > 0 {
						leftPath = append(leftPath, 16)
					}
					if len(leftNibble) > 0 {
						leftNibble = append(leftNibble, 16)
					}
					fmt.Println("Left path appended :", leftPath)
					fmt.Println("Left nibble appended:", leftNibble)
					if len(leftPath) > 0 {
						nodeL1 := Node{}
						nodeL1.node_type = 2
						index := leftPath[0]    //now
						leftPath = leftPath[1:] //now
						fmt.Println("Left path:", leftPath)
						nodeL1.flag_value.encoded_prefix = compact_encode(leftPath)
						nodeL1.flag_value.value = newValue
						nodeB.branch_value[index] = mpt.InsertRecursive(nodeL1, nil, "", nil)
						fmt.Println("Inserted nodeL1")
					} else if len(leftPath) == 0 {
						nodeB.branch_value[16] = newValue
						fmt.Println("Inserted in the value field of the branch")
					}
					if len(leftNibble) > 1 {
						nodeL2 := Node{}
						nodeL2.node_type = 2
						index := leftNibble[0]      //now
						leftNibble = leftNibble[1:] //now
						fmt.Println("Left nibble:", leftNibble)

						nodeL2.flag_value.encoded_prefix = compact_encode(leftNibble)
						nodeL2.flag_value.value = currNode.flag_value.value
						nodeB.branch_value[index] = mpt.InsertRecursive(nodeL2, nil, "", nil)
						fmt.Println("Inserted nodeL2")
					} else if len(leftNibble) == 0 {
						nodeB.branch_value[16] = currNode.flag_value.value
						fmt.Println("Inserted in the value field of the branch")
					}

					nodeE.flag_value.value = mpt.InsertRecursive(nodeB, nil, "", nil)
					fmt.Println("Inserted nodeB")

					hashE := nodeE.hash_node()
					delete(mpt.db, old_hash)
					mpt.db[hashE] = nodeE
					fmt.Println("Inserted nodeE")
					return hashE
				} else {
					//create branch
					//create leaves or leaf
					fmt.Println("In leaf and the left path is totally different than a leaf")
					//put in the branch[16] field
					//also check if the leaf has just one value
					nodeB := Node{}
					nodeB.node_type = 1

					if len(decodedHexArray) > 0 {
						decodedHexArray = append(decodedHexArray, 16)
						nodeL2 := Node{}
						nodeL2.node_type = 2
						fmt.Println("decodedHexArray :", decodedHexArray)
						nodeL2.flag_value.encoded_prefix = compact_encode(decodedHexArray[1:])
						nodeL2.flag_value.value = currNode.flag_value.value
						nodeB.branch_value[decodedHexArray[0]] = mpt.InsertRecursive(nodeL2, nil, "", nil)
					} else if len(decodedHexArray) == 0 {
						nodeB.branch_value[16] = currNode.flag_value.value
					}

					if len(leftPath) > 0 {
						nodeL1 := Node{}
						nodeL1.node_type = 2
						leftPath = append(leftPath, 16)
						nodeL1.flag_value.encoded_prefix = compact_encode(leftPath[1:])
						nodeL1.flag_value.value = newValue
						nodeB.branch_value[leftPath[0]] = mpt.InsertRecursive(nodeL1, nil, "", nil)

					} else {
						nodeB.branch_value[16] = newValue
					}
					fmt.Println("leftPath :", leftPath)
					hashB := nodeB.hash_node()
					mpt.db[hashB] = nodeB
					fmt.Println("Branch and two leaves inserted")
					delete(mpt.db, old_hash)
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
					fmt.Println("Next node type : ", currNode.node_type)
					//check if next node is a leaf
					//if yes convert it into branch
					//insert this value in branchvalue[16]
					//check the lenght of the leaf it 1 then create leaf store empty value
					//lenght of leaf is 0 .... i think it is the same
					//lenght of leaf is > 1 ....store in leaf
					//lib
					currNode.flag_value.value = mpt.InsertRecursive(mpt.db[currNode.flag_value.value], nil, newValue, nil)
					hash := currNode.hash_node()
					delete(mpt.db, old_hash)
					mpt.db[hash] = currNode
					fmt.Println("Returning hash")
					return hash
					//lib
				} else if check_if_common_path_exists(leftPath, decodedHexArray) {
					counter := 0
					for i := 0; i < len(decodedHexArray); i++ {
						//now me me me
						if i == len(decodedHexArray) || i == len(leftPath) {
							break
						}
						if decodedHexArray[i] == leftPath[i] {
							counter = counter + 1
						} else {
							break
						}
					}
					// for i := 0; i < len(decodedHexArray); i++ {
					// 	if decodedHexArray[i] == leftPath[i] {
					// 		counter++
					// 	} else {
					// 		break
					// 	}
					//}

					fmt.Println("Counter:", counter)
					common_path2 := leftPath[:counter]
					leftPath = leftPath[counter:]
					leftNibble = decodedHexArray[counter:]
					currNode.flag_value.encoded_prefix = compact_encode(common_path2)
					fmt.Println("common_path2:", common_path2)
					fmt.Println("leftPath:", leftPath)
					fmt.Println("leftNibble:", leftNibble)
					if len(leftNibble) > 0 { // lib
						nodeBranch := Node{}
						nodeBranch.node_type = 1
						if len(leftPath) > 0 {
							//create Branch
							//left path create branch and leaf
							leftPath = append(leftPath, 16)

							nodeLeaf := Node{}
							nodeLeaf.node_type = 2
							nodeLeaf.flag_value.encoded_prefix = compact_encode(leftPath[1:])
							nodeLeaf.flag_value.value = newValue
							fmt.Println("Inserting in branch at index:", leftPath[0])
							nodeBranch.branch_value[leftPath[0]] = mpt.InsertRecursive(nodeLeaf, nil, "", nil)
						} else {
							nodeBranch.branch_value[16] = newValue
						}
						//check the left nibble size
						if len(leftNibble) == 0 {
							currNode.flag_value.value = mpt.InsertRecursive(mpt.db[currNode.flag_value.value], leftPath, newValue, nil)
							hashExt := currNode.hash_node()
							delete(mpt.db, old_hash)
							mpt.db[hashExt] = currNode
							return hashExt
						} else if len(leftNibble) == 1 {
							nodeBranch.branch_value[leftNibble[0]] = currNode.flag_value.value
							currNode.flag_value.value = mpt.InsertRecursive(nodeBranch, nil, "", nil)
							hashExt := currNode.hash_node()
							delete(mpt.db, old_hash)
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
							delete(mpt.db, old_hash)
							mpt.db[hashExt] = currNode
							return hashExt
						}
					} else { //lib
						fmt.Println("Proceeding ahead from extension") //lib
						currNode.flag_value.value = mpt.InsertRecursive(mpt.db[currNode.flag_value.value], leftPath, newValue, nil)
						hashExt := currNode.hash_node()
						delete(mpt.db, old_hash)
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
					fmt.Println("Left path:", leftPath)
					nodeLeaf := Node{}
					nodeLeaf.node_type = 2
					nodeLeaf.flag_value.encoded_prefix = compact_encode(leftPath[1:])

					//creating a leaf //pop

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
					fmt.Println("Node:", nodeBranch)
					delete(mpt.db, old_hash)
					fmt.Println("Returning hash")
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
				fmt.Println("Next node type :", nextNode.node_type)
				leftPath = leftPath[1:]
				currNode.branch_value[index] = mpt.InsertRecursive(nextNode, leftPath, newValue, nil)
				fmt.Println("Value updated in branch")
			}
			hashBr := currNode.hash_node()
			delete(mpt.db, old_hash)
			mpt.db[hashBr] = currNode
			fmt.Println("Branch updated")
			return hashBr
		}
	}
	return ""
}
