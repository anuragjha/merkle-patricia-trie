package p1

import (
	"errors"
	"fmt"
)

func (mpt *MerklePatriciaTrie) getHelper(currentNode Node, pathLeft []uint8) (string, error) {

	if len(pathLeft) < 0 {
		fmt.Println("pathleft is less than 0")
		return "", errors.New("path_not_found")

	} else if len(pathLeft) == 0 {
		fmt.Println("pathleft is equal 0")
		if currentNode.node_type == 1 {
			fmt.Println("pathleft is equal 0 - node type 1")
			if currentNode.branch_value[16] != "" {
				//say found and return
				return currentNode.branch_value[16], nil
			} else {
				return "", errors.New("path_not_found")
			}
		} else if currentNode.node_type == 2 {
			fmt.Println("pathleft is equal 0 - node type 2")
			hex_prefix_array := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
			if hex_prefix_array[0] == 2 || hex_prefix_array[0] == 3 {
				fmt.Println("path length 0 next node is leaf")
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
		fmt.Println("Current Node Type:", currentNode.node_type)
		fmt.Println("Node:", currentNode)
		fmt.Println("pathleft is greater than 0")
		if currentNode.node_type == 1 { //branch node
			fmt.Println("pathleft is greater than 0 - branch node")
			fmt.Println("pathLeft", pathLeft)
			fmt.Println("currentNode.branch_value", currentNode.branch_value[5])
			if currentNode.branch_value[pathLeft[0]] != "" {
				fmt.Println("Value exists in branch!!")
				fmt.Println("Hash found:", currentNode.branch_value[pathLeft[0]])
				hash := currentNode.branch_value[pathLeft[0]]
				pathLeft = pathLeft[1:]
				fmt.Println("Path Left from branch :", pathLeft)
				return mpt.getHelper(mpt.db[hash], pathLeft)
			} else {
				fmt.Println("value not found in branch:", currentNode.branch_value[pathLeft[0]])
				fmt.Println("Branch value at the index when not found")
				return "", errors.New("path_not_found")
			}
		} else if currentNode.node_type == 2 {
			fmt.Println("pathleft is greater than 0 - node_type 2")
			hex_prefix_array := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)
			//fmt.Println("!!!!!!!!!!!HEX PREFIX :", hex_prefix_array)
			if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) {
				fmt.Println("pathleft is greater than 0 - node_type 2 - prefix 0 or 1 Extension")
				counter := 0
				triePath := compact_decode(currentNode.flag_value.encoded_prefix)
				fmt.Println("currentNode.flag_value.encoded_prefix:", currentNode.flag_value.encoded_prefix)
				fmt.Println("triePathtriePathtriePath:", triePath)
				fmt.Println("pathLeftpathLeftpathLeft:", pathLeft)
				for i := 0; i < len(triePath); i++ {
					if triePath[i] == pathLeft[i] {
						counter++
					}
				}
				fmt.Println("counter34", counter)
				if counter == (len(triePath)) {
					fmt.Println("counter54", counter)
					pathLeft = pathLeft[counter:]
					fmt.Println("Path left before call:", pathLeft)
					return mpt.getHelper(mpt.db[currentNode.flag_value.value], pathLeft)
				} else {
					return "", errors.New("path_not_found")
				}
			} else if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) {
				fmt.Println("pathleft is greater than 0 - node_type 2 - prefix 2 or 3 Leaf")
				counter := 0
				triePath := compact_decode(currentNode.flag_value.encoded_prefix)
				fmt.Println("leftPath :", pathLeft)
				fmt.Println("triePath", triePath)
				//DEEP EQUAL
				if len(pathLeft) < len(triePath) {
					return "", errors.New("path_not_found")
				}
				for i := 0; i < len(triePath); i++ {
					if triePath[i] == pathLeft[i] {
						counter++
					}
				}
				fmt.Println("Counter", counter)
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
