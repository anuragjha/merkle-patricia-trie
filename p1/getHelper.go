package p1

import (
	"errors"
	"reflect"
)

func (mpt *MerklePatriciaTrie) getHelper(currentNode Node, pathLeft []uint8) (string, error) {

	if len(pathLeft) == 0 { //should find value of key

		if currentNode.node_type == 1 { // value in branch
			if currentNode.branch_value[16] != "" {
				return currentNode.branch_value[16], nil //value found
			}

		} else if currentNode.node_type == 2 { // value in leaf

			currentNodeHexArray := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
			if currentNodeHexArray[0] == 2 || currentNodeHexArray[0] == 3 { // leaf //pathleft gets to 0 in previous node
				if len(compact_decode(currentNode.flag_value.encoded_prefix)) == 0 { //coz pathLeft = 0
					return currentNode.flag_value.value, nil
				}
			}
		}

	} else if len(pathLeft) > 0 { // still to match pathLeft, pathLeft > 0

		if currentNode.node_type == 1 { //branch node
			if currentNode.branch_value[pathLeft[0]] != "" {
				nextNodeHash := currentNode.branch_value[pathLeft[0]]
				pathLeft = pathLeft[1:]
				return mpt.getHelper(mpt.db[nextNodeHash], pathLeft)
			}

		} else if currentNode.node_type == 2 { //extension or leaf

			currentNodeHexArray := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
			currentNodePath := compact_decode(currentNode.flag_value.encoded_prefix)
			if (currentNodeHexArray[0] == 0) || (currentNodeHexArray[0] == 1) { //extension
				if reflect.DeepEqual(currentNodePath, pathLeft[:len(currentNodePath)]) {
					pathLeft = pathLeft[len(currentNodePath):]
					return mpt.getHelper(mpt.db[currentNode.flag_value.value], pathLeft)
				}

			} else if (currentNodeHexArray[0] == 2) || (currentNodeHexArray[0] == 3) { //leaf
				if len(pathLeft) >= len(currentNodePath) {
					if reflect.DeepEqual(currentNodePath, pathLeft[:len(currentNodePath)]) {

						pathLeft = pathLeft[len(currentNodePath):]
						if len(pathLeft) == 0 {
							return currentNode.flag_value.value, nil
						} else {
							return mpt.getHelper(mpt.db[currentNode.flag_value.value], pathLeft)
						}
					}
				}
			}
		}
	}
	return "", errors.New("path_not_found")
}
