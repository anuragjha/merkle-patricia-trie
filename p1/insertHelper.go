package p1

import (
	"reflect"
)

func (mpt *MerklePatriciaTrie) insertHelper(node Node, leftPath []uint8, newValue string, leftNibble []uint8) string {
	currentNode := node
	currentNodeOldHash := currentNode.hash_node()
	if len(leftPath) == 0 && newValue == "" { //initialize
		currentNodeNewHash := currentNode.hash_node()
		if len(mpt.db) == 0 {
			mpt.db = make(map[string]Node)
			mpt.root = ""
		}
		mpt.db[currentNodeNewHash] = currentNode
		return currentNodeNewHash
	} else if len(leftPath) == 0 { ////// pathLeft = 0

		if currentNode.node_type == 1 { // branch - value to be added in branch
			currentNode.branch_value[16] = newValue
			currentNodeNewHash := currentNode.hash_node()
			delete(mpt.db, currentNodeOldHash)
			mpt.db[currentNodeNewHash] = currentNode
			return currentNodeNewHash
		} else if currentNode.node_type == 2 { // only leaf condition here
			currentNodeHexArray := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
			if currentNodeHexArray[0] == 2 || currentNodeHexArray[0] == 3 { //leaf

				currentNodeHexArrayDecoded := compact_decode(currentNode.flag_value.encoded_prefix)
				if reflect.DeepEqual(currentNodeHexArray, []uint8{2, 0}) { //empty key
					currentNode.flag_value.value = newValue
					currentNodeNewHash := currentNode.hash_node()
					delete(mpt.db, currentNodeOldHash)
					mpt.db[currentNodeNewHash] = currentNode
					return currentNodeNewHash
				} else if len(currentNodeHexArrayDecoded) > 0 { //pathleft in leaf is > 0
					currentNodeHexArrayDecoded = append(currentNodeHexArrayDecoded, 16)
					//leaf - previous leaf value to be moved to branch-leaf
					nodeLeaf := Node{}
					nodeLeaf.node_type = 2
					nodeLeaf.flag_value.value = currentNode.flag_value.value
					nodeLeaf.flag_value.encoded_prefix = compact_encode(currentNodeHexArrayDecoded[1:])
					//branch - newValue to be added in 16 of branch
					nodeBranch := Node{}
					nodeBranch.node_type = 1
					nodeBranch.branch_value[16] = newValue
					nodeBranch.branch_value[currentNodeHexArrayDecoded[0]] = mpt.insertHelper(nodeLeaf, nil, "", nil)
					branchNodeHash := nodeBranch.hash_node()
					delete(mpt.db, currentNodeOldHash)
					mpt.db[branchNodeHash] = nodeBranch
					return branchNodeHash
				}
			} else {
				return "unable to insert"
			}
		}
	} else { /////// pathLeft > 0
		if currentNode.node_type == 2 { //ext or leaf
			path := currentNode.flag_value.encoded_prefix
			currentNodeHexArray := AsciiArrayToHexArray(path)
			currentNodeHexArrayDecoded := compact_decode(path)

			if currentNodeHexArray[0] == 2 || currentNodeHexArray[0] == 3 { // leaf

				if reflect.DeepEqual(currentNodeHexArrayDecoded, leftPath) {
					currentNode.flag_value.value = newValue
					leftPath = nil
					delete(mpt.db, currentNodeOldHash)
					currentNodeNewHash := currentNode.hash_node()
					mpt.db[currentNodeNewHash] = currentNode
					return currentNodeNewHash

				} else if checkForCommonPath(leftPath, currentNodeHexArrayDecoded) { //common path exists

					counter := 0
					for i := 0; i < len(currentNodeHexArrayDecoded); i++ {
						//now me me me
						if i == len(currentNodeHexArrayDecoded) || i == len(leftPath) {
							break
						}
						if currentNodeHexArrayDecoded[i] == leftPath[i] {
							counter = counter + 1
						} else {
							break
						}
					}

					commonPath := leftPath[:(counter)]
					leftPath = leftPath[counter:]
					leftNibble = currentNodeHexArrayDecoded[counter:]

					//create extension node and branch
					nodeExt := Node{}
					nodeExt.node_type = 2
					nodeExt.flag_value.encoded_prefix = compact_encode(commonPath)
					//create branch
					nodeBranch := Node{}
					nodeBranch.node_type = 1

					// if len(leftPath) > 0 && len(leftNibble) > 0 {
					//create two leaves
					if len(leftPath) > 0 {
						leftPath = append(leftPath, 16)
					}
					if len(leftNibble) > 0 {
						leftNibble = append(leftNibble, 16)
					}

					if len(leftPath) > 0 {
						nodeLeaf1 := Node{}
						nodeLeaf1.node_type = 2
						index := leftPath[0]
						leftPath = leftPath[1:]

						nodeLeaf1.flag_value.encoded_prefix = compact_encode(leftPath)
						nodeLeaf1.flag_value.value = newValue
						nodeBranch.branch_value[index] = mpt.insertHelper(nodeLeaf1, nil, "", nil)

					} else if len(leftPath) == 0 {
						nodeBranch.branch_value[16] = newValue

					}
					if len(leftNibble) > 1 {
						nodeLeaf2 := Node{}
						nodeLeaf2.node_type = 2
						index := leftNibble[0]
						leftNibble = leftNibble[1:]

						nodeLeaf2.flag_value.encoded_prefix = compact_encode(leftNibble)
						nodeLeaf2.flag_value.value = currentNode.flag_value.value
						nodeBranch.branch_value[index] = mpt.insertHelper(nodeLeaf2, nil, "", nil)

					} else if len(leftNibble) == 0 {
						nodeBranch.branch_value[16] = currentNode.flag_value.value
					}

					nodeExt.flag_value.value = mpt.insertHelper(nodeBranch, nil, "", nil)

					hashE := nodeExt.hash_node()
					delete(mpt.db, currentNodeOldHash)
					mpt.db[hashE] = nodeExt

					return hashE
				} else {

					nodeBranch := Node{}
					nodeBranch.node_type = 1

					if len(currentNodeHexArrayDecoded) > 0 {
						currentNodeHexArrayDecoded = append(currentNodeHexArrayDecoded, 16)
						nodeLeaf2 := Node{}
						nodeLeaf2.node_type = 2
						nodeLeaf2.flag_value.encoded_prefix = compact_encode(currentNodeHexArrayDecoded[1:])
						nodeLeaf2.flag_value.value = currentNode.flag_value.value
						nodeBranch.branch_value[currentNodeHexArrayDecoded[0]] = mpt.insertHelper(nodeLeaf2, nil, "", nil)

					} else if len(currentNodeHexArrayDecoded) == 0 {
						nodeBranch.branch_value[16] = currentNode.flag_value.value
					}

					if len(leftPath) > 0 {
						nodeL1 := Node{}
						nodeL1.node_type = 2
						leftPath = append(leftPath, 16)
						nodeL1.flag_value.encoded_prefix = compact_encode(leftPath[1:])
						nodeL1.flag_value.value = newValue
						nodeBranch.branch_value[leftPath[0]] = mpt.insertHelper(nodeL1, nil, "", nil)

					} else {
						nodeBranch.branch_value[16] = newValue
					}

					hashB := nodeBranch.hash_node()
					mpt.db[hashB] = nodeBranch

					delete(mpt.db, currentNodeOldHash)
					return hashB
				}
			} else if currentNodeHexArray[0] == 0 || currentNodeHexArray[0] == 1 { //extension

				if reflect.DeepEqual(currentNodeHexArrayDecoded, leftPath) {
					currentNode.flag_value.value = mpt.insertHelper(mpt.db[currentNode.flag_value.value], nil, newValue, nil)
					currentNodeNewHash := currentNode.hash_node()
					delete(mpt.db, currentNodeOldHash)
					mpt.db[currentNodeNewHash] = currentNode
					return currentNodeNewHash
				} else if checkForCommonPath(leftPath, currentNodeHexArrayDecoded) {
					counter := 0
					for i := 0; i < len(currentNodeHexArrayDecoded); i++ {
						//now me me me
						if i == len(currentNodeHexArrayDecoded) || i == len(leftPath) {
							break
						}
						if currentNodeHexArrayDecoded[i] == leftPath[i] {
							counter = counter + 1
						} else {
							break
						}
					}

					commonPath := leftPath[:counter]
					leftPath = leftPath[counter:]
					leftNibble = currentNodeHexArrayDecoded[counter:]
					currentNode.flag_value.encoded_prefix = compact_encode(commonPath)

					if len(leftNibble) > 0 {
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

							nodeBranch.branch_value[leftPath[0]] = mpt.insertHelper(nodeLeaf, nil, "", nil)
						} else {
							nodeBranch.branch_value[16] = newValue
						}
						//check the left nibble size
						if len(leftNibble) == 0 {
							currentNode.flag_value.value = mpt.insertHelper(mpt.db[currentNode.flag_value.value], leftPath, newValue, nil)
							hashExt := currentNode.hash_node()
							delete(mpt.db, currentNodeOldHash)
							mpt.db[hashExt] = currentNode
							return hashExt
						} else if len(leftNibble) == 1 {
							nodeBranch.branch_value[leftNibble[0]] = currentNode.flag_value.value
							currentNode.flag_value.value = mpt.insertHelper(nodeBranch, nil, "", nil)
							hashExt := currentNode.hash_node()
							delete(mpt.db, currentNodeOldHash)
							mpt.db[hashExt] = currentNode
							return hashExt
						} else if len(leftNibble) > 1 {
							//create an extension and store the value of that extension in branch => branch hash store in currNode
							nodeExt := Node{}
							nodeExt.node_type = 2
							nodeExt.flag_value.encoded_prefix = compact_encode(leftNibble[1:])
							nodeExt.flag_value.value = currentNode.flag_value.value
							nodeBranch.branch_value[leftNibble[0]] = mpt.insertHelper(nodeExt, nil, "", nil)
							currentNode.flag_value.value = mpt.insertHelper(nodeBranch, nil, "", nil)
							hashExt := currentNode.hash_node()
							delete(mpt.db, currentNodeOldHash)
							mpt.db[hashExt] = currentNode
							return hashExt
						}
					} else {

						currentNode.flag_value.value = mpt.insertHelper(mpt.db[currentNode.flag_value.value], leftPath, newValue, nil)
						hashExt := currentNode.hash_node()
						delete(mpt.db, currentNodeOldHash)
						mpt.db[hashExt] = currentNode
						return hashExt
					}

				} else { //leaf hash
					//make extension node a branch and one leaf
					nodeBranch := Node{}
					nodeBranch.node_type = 1
					leftPath = append(leftPath, 16)

					nodeLeaf := Node{}
					nodeLeaf.node_type = 2
					nodeLeaf.flag_value.encoded_prefix = compact_encode(leftPath[1:])

					//creating a leaf

					nodeLeaf.flag_value.value = newValue
					nodeBranch.branch_value[leftPath[0]] = mpt.insertHelper(nodeLeaf, nil, "", nil)

					if len(currentNodeHexArrayDecoded) == 1 {

						nodeBranch.branch_value[currentNodeHexArrayDecoded[0]] = currentNode.flag_value.value
					} else if len(currentNodeHexArrayDecoded) > 1 {
						nodeExt := Node{}
						nodeExt.node_type = 2
						nodeExt.flag_value.encoded_prefix = compact_encode(currentNodeHexArrayDecoded[1:])
						nodeExt.flag_value.value = currentNode.flag_value.value
						nodeBranch.branch_value[currentNodeHexArrayDecoded[0]] = mpt.insertHelper(nodeExt, nil, "", nil)
					}

					newBranchNodeHash := nodeBranch.hash_node()
					mpt.db[newBranchNodeHash] = nodeBranch
					delete(mpt.db, currentNodeOldHash)

					return newBranchNodeHash
				}
			}
		} else if currentNode.node_type == 1 {
			if currentNode.branch_value[leftPath[0]] == "" {
				//store leftPath[0] create a leaf to store the rest

				leftPath = append(leftPath, 16)
				nodeL3 := Node{}
				nodeL3.node_type = 2
				nodeL3.flag_value.encoded_prefix = compact_encode(leftPath[1:])
				nodeL3.flag_value.value = newValue
				currentNode.branch_value[leftPath[0]] = mpt.insertHelper(nodeL3, nil, "", nil)
			} else if currentNode.branch_value[leftPath[0]] != "" {

				index := leftPath[0]
				nextNode := mpt.db[currentNode.branch_value[leftPath[0]]]

				leftPath = leftPath[1:]
				currentNode.branch_value[index] = mpt.insertHelper(nextNode, leftPath, newValue, nil)

			}
			currentNodeNewHash := currentNode.hash_node()
			delete(mpt.db, currentNodeOldHash)
			mpt.db[currentNodeNewHash] = currentNode

			return currentNodeNewHash
		}
	}
	return "unable to insert"
}

func checkForCommonPath(leftPath []uint8, decodedHexArray []uint8) bool {
	if len(leftPath) > 0 && len(decodedHexArray) > 0 {
		if leftPath[0] == decodedHexArray[0] {
			return true
		}
	}
	return false
}
